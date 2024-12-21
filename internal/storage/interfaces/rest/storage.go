package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
)

type StorageController struct {
	queryHandler   *handlers.StorageQueryHandler
	commandHandler *handlers.StorageCommandHandler
	service        *service.StorageService
}

func NewStorageController(
	queryHandler *handlers.StorageQueryHandler,
	commandHandler *handlers.StorageCommandHandler,
	service *service.StorageService,
) *StorageController {
	return &StorageController{
		queryHandler:   queryHandler,
		commandHandler: commandHandler,
		service:        service,
	}
}

func (s *StorageController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	storage := g.Group("/v1/storage", jwt.Handler(t))
	{
		// 文件夹管理
		storage.GET("/folders", hserver.NewHandlerFu[queries.ListFoldersQuery](s.ListFolders))
		storage.GET("/folders/tree", hserver.NewHandlerFu[queries.GetFolderTreeQuery](s.GetFolderTree))
		storage.POST("/folders", hserver.NewHandlerFu[commands.CreateFolderCommand](s.CreateFolder))
		storage.PUT("/folders/:id/move", hserver.NewHandlerFu[commands.MoveFolderCommand](s.MoveFolder))
		storage.PUT("/folders/:id/rename", hserver.NewHandlerFu[commands.RenameFolderCommand](s.RenameFolder))
		storage.DELETE("/folders/:id", hserver.NewHandlerFu[commands.DeleteFolderCommand](s.DeleteFolder))

		// 文件管理
		storage.GET("/files", hserver.NewHandlerFu[queries.ListFilesQuery](s.ListFiles))
		storage.POST("/files", s.UploadFile)
		storage.DELETE("/files/:id", hserver.NewHandlerFu[commands.DeleteFileCommand](s.DeleteFile))
		storage.PUT("/files/:id/move", hserver.NewHandlerFu[commands.MoveFileCommand](s.MoveFile))

		// 批量操作
		storage.POST("/files/batch/delete", hserver.NewHandlerFu[commands.BatchDeleteFilesCommand](s.BatchDeleteFiles))
		storage.POST("/files/batch/move", hserver.NewHandlerFu[commands.BatchMoveFilesCommand](s.BatchMoveFiles))

		// 文件分享
		storage.POST("/files/:id/share", hserver.NewHandlerFu[commands.ShareFileCommand](s.ShareFile))
		storage.GET("/share/:code", s.GetFileShare) // 不需要认证

		// 回收站
		storage.GET("/recycle", hserver.NewHandlerFu[queries.ListRecycleFilesQuery](s.ListRecycleFiles))
		storage.POST("/files/:id/recycle", hserver.NewHandlerFu[commands.RecycleFileCommand](s.RecycleFile))
		storage.POST("/files/:id/restore", hserver.NewHandlerFu[commands.RestoreFileCommand](s.RestoreFile))

		// 预览文件
		storage.GET("/files/:id/preview", hserver.NewHandlerFu[queries.GetFilePreviewQuery](s.PreviewFile))

		// 下载文件
		storage.GET("/files/:id/download", s.DownloadFile)
	}
}

// ListFolders 查询文件夹列表
// @Summary 查询文件夹列表
// @Description 查询文件夹列表
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param parent_id query int64 false "父文件夹ID"
// @Param name query string false "文件夹名称"
// @Param current query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} base_info.Success{data=[]dto.FolderDto}
// @Router /v1/storage/folders [get]
func (s *StorageController) ListFolders(ctx context.Context, q *queries.ListFoldersQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	q.TenantID = actx.GetTenantId(ctx)
	data, err := s.queryHandler.HandleListFolders(ctx, q)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// ListFiles 查询文件列表
// @Summary 查询文件列表
// @Description 查询文件列表
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param folder_id query int64 false "文件夹ID"
// @Param name query string false "文件名"
// @Param type query string false "文件类型"
// @Param storage_type query string false "存储类型"
// @Param current query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} base_info.Success{data=[]dto.FileDto}
// @Router /v1/storage/files [get]
func (s *StorageController) ListFiles(ctx context.Context, q *queries.ListFilesQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	q.TenantID = actx.GetTenantId(ctx)
	data, err := s.queryHandler.HandleListFiles(ctx, q)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// CreateFolder 创建文件夹
// @Summary 创建文件夹
// @Description 创建文件夹
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param data body commands.CreateFolderCommand true "创建参数"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/folders [post]
func (s *StorageController) CreateFolder(ctx context.Context, cmd *commands.CreateFolderCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()

	// 创建文件夹
	err := s.commandHandler.HandleCreateFolder(ctx, cmd, actx.GetTenantId(ctx), actx.GetUserId(ctx))
	if err != nil {
		return result.WithError(err)
	}

	return result
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件
// @Tags 存储管理
// @Accept multipart/form-data
// @Produce json
// @Param folder_id formData int64 false "文件夹ID,默认为根目录"
// @Param storage_type formData string false "存储类型(minio/aliyun/tencent),默认为minio"
// @Param file formData file true "文件"
// @Success 200 {object} base_info.Success{data=dto.FileDto}
// @Router /v1/storage/files [post]
func (s *StorageController) UploadFile(ctx context.Context, c *app.RequestContext) {
	result := hserver.DefaultResponseResult()

	// 获取文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, result.WithError(herrors.NewBadReqError("获取文件失败")))
		return
	}

	// 读取文件内容
	fileContent, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusOK, result.WithError(herrors.NewBadReqError("读取文件失败")))
		return
	}
	defer fileContent.Close()

	// 读取文件内容到内存
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusOK, result.WithError(herrors.NewBadReqError("读取文件失败")))
		return
	}
	// 获取文件夹ID,默认为"0"(根目录)
	folderID := "0"
	if id := c.FormValue("folder_id"); len(id) > 0 {
		folderID = string(id)
	}
	// 构建上传命令
	cmd := &commands.UploadFileCommand{
		FolderID: folderID,
		File:     fileBytes,
		FileName: fileHeader.Filename,
	}

	// 上传文件
	fileModel, herr := s.commandHandler.HandleUploadFile(ctx, cmd, actx.GetTenantId(ctx), actx.GetUserId(ctx))
	if herr != nil {
		hserver.ResponseFailureErr(ctx, c, herr)
		return
	}

	// 转换为DTO
	fileDto := &dto.FileDto{
		ID:          fileModel.ID,
		Name:        fileModel.Name,
		Path:        fileModel.Path,
		FolderID:    fileModel.FolderID,
		Size:        fileModel.Size,
		Type:        fileModel.Type,
		StorageType: string(fileModel.StorageType),
		URL:         fileModel.URL,
		CreatedBy:   fileModel.CreatedBy,
		CreatedAt:   fileModel.CreatedAt,
	}

	c.JSON(http.StatusOK, result.WithData(fileDto))
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param id path int64 true "文件ID"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/files/{id} [delete]
func (s *StorageController) DeleteFile(ctx context.Context, cmd *commands.DeleteFileCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleDeleteFile(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteFolder 删除文件夹
func (s *StorageController) DeleteFolder(ctx context.Context, cmd *commands.DeleteFolderCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleDeleteFolder(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// RenameFolder 重命名文件夹
func (s *StorageController) RenameFolder(ctx context.Context, cmd *commands.RenameFolderCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleRenameFolder(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// MoveFile 移动文件
func (s *StorageController) MoveFile(ctx context.Context, cmd *commands.MoveFileCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleMoveFile(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// GetFolderTree 获取文件夹树形结构
// @Summary 获取文件夹树形结构
// @Description 获取文件夹树形结构
// @Tags 存储管理
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=[]dto.FolderTreeDto}
// @Router /v1/storage/folders/tree [get]
func (s *StorageController) GetFolderTree(ctx context.Context, q *queries.GetFolderTreeQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	q.TenantID = actx.GetTenantId(ctx)
	data, err := s.queryHandler.HandleGetFolderTree(ctx, q)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// PreviewFile 预览文件
// @Summary 预览文件
// @Description 预览文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} base_info.Success{data=string}
// @Router /v1/storage/files/{id}/preview [get]
func (s *StorageController) PreviewFile(ctx context.Context, q *queries.GetFilePreviewQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()

	// 获取预览URL
	previewURL, err := s.service.PreviewFile(ctx, q.ID)
	if err != nil {
		return result.WithError(herrors.NewErr(err))
	}

	return result.WithData(previewURL)
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 下载文件
// @Tags 存储管理
// @Accept json
// @Produce octet-stream
// @Param id path int true "文件ID"
// @Router /v1/storage/files/{id}/download [get]
func (s *StorageController) DownloadFile(ctx context.Context, c *app.RequestContext) {
	// 获取文件ID
	fileID := c.Param("id")
	if fileID == "" {
		c.String(http.StatusBadRequest, "file id is empty")
		return
	}

	// 下载文件
	reader, filename, err := s.service.DownloadFile(ctx, fileID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(filename)))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件内容
	if _, err := io.Copy(c.Response.BodyWriter(), reader); err != nil {
		hlog.Errorf("send file error: %v", err)
	}
}

// GetFileShare 获取分享文件
// @Summary 获取分享文件
// @Description 获取分享文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param code path string true "分享码"
// @Param password query string false "访问密码"
// @Success 200 {object} base_info.Success{data=dto.FileDto}
// @Router /v1/storage/share/{code} [get]
func (s *StorageController) GetFileShare(ctx context.Context, c *app.RequestContext) {
	result := hserver.DefaultResponseResult()

	// 获取分享码和密码
	shareCode := c.Param("code")
	password := c.Query("password")

	// 获取分享文件
	file, err := s.service.GetShareFile(ctx, shareCode, password)
	if err != nil {
		c.JSON(http.StatusOK, result.WithError(herrors.NewErr(err)))
		return
	}

	// 转换为DTO
	fileDto := &dto.FileDto{
		ID:          file.ID,
		Name:        file.Name,
		Path:        file.Path,
		Size:        file.Size,
		Type:        file.Type,
		StorageType: string(file.StorageType),
		URL:         file.URL,
		CreatedAt:   file.CreatedAt,
	}

	c.JSON(http.StatusOK, result.WithData(fileDto))
}

// BatchDeleteFiles 批量删除文件
// @Summary 批量删除文件
// @Description 批量删除文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param data body commands.BatchDeleteFilesCommand true "文件ID列表"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/files/batch/delete [post]
func (s *StorageController) BatchDeleteFiles(ctx context.Context, cmd *commands.BatchDeleteFilesCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleBatchDeleteFiles(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// BatchMoveFiles 批量移动文件
// @Summary 批量移动文件
// @Description 批量移动文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param data body commands.BatchMoveFilesCommand true "移动参数"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/files/batch/move [post]
func (s *StorageController) BatchMoveFiles(ctx context.Context, cmd *commands.BatchMoveFilesCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleBatchMoveFiles(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// ShareFile 分享文件
// @Summary 分享文件
// @Description 分享文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param id path int64 true "文件ID"
// @Param data body commands.ShareFileCommand true "分享参数"
// @Success 200 {object} base_info.Success{data=dto.FileShareDto}
// @Router /v1/storage/files/{id}/share [post]
func (s *StorageController) ShareFile(ctx context.Context, cmd *commands.ShareFileCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	cmd.CreatedBy = actx.GetUserId(ctx)
	share, err := s.commandHandler.HandleShareFile(ctx, cmd)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(share)
}

// MoveFolder 移动文件夹
func (s *StorageController) MoveFolder(ctx context.Context, cmd *commands.MoveFolderCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleMoveFolder(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// ListRecycleFiles 查询回收站文件列表
// @Summary 查询回收��文件列表
// @Description 查询回收站文件列表
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param name query string false "文件名"
// @Param type query string false "文件类型"
// @Param storage_type query string false "存储类型"
// @Param current query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} base_info.Success{data=dto.FileDto}
// @Router /v1/storage/recycle [get]
func (s *StorageController) ListRecycleFiles(ctx context.Context, q *queries.ListRecycleFilesQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	q.TenantID = actx.GetTenantId(ctx)
	data, err := s.queryHandler.HandleListRecycleFiles(ctx, q)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// RecycleFile 回收文件
// @Summary 回收文件
// @Description 将文件移动到回收站
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param id path int64 true "文件ID"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/files/{id}/recycle [post]
func (s *StorageController) RecycleFile(ctx context.Context, cmd *commands.RecycleFileCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	cmd.DeletedBy = actx.GetUserId(ctx)
	if err := s.commandHandler.HandleRecycleFile(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// RestoreFile 恢复文件
// @Summary 恢复文件
// @Description 从回收站恢复文件
// @Tags 存储管理
// @Accept json
// @Produce json
// @Param id path int64 true "文件ID"
// @Success 200 {object} base_info.Success
// @Router /v1/storage/files/{id}/restore [post]
func (s *StorageController) RestoreFile(ctx context.Context, cmd *commands.RestoreFileCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := s.commandHandler.HandleRestoreFile(ctx, cmd); err != nil {
		return result.WithError(err)
	}
	return result
}

// ... 其他方法实现
