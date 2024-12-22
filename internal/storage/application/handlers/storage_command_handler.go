package handlers

import (
	"bytes"
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type StorageCommandHandler struct {
	service *service.StorageService
}

func NewStorageCommandHandler(service *service.StorageService) *StorageCommandHandler {
	return &StorageCommandHandler{
		service: service,
	}
}

// HandleCreateFolder 处理创建文件夹
func (h *StorageCommandHandler) HandleCreateFolder(ctx context.Context, cmd *commands.CreateFolderCommand, tenantID, createdBy string) herrors.Herr {
	// 1. 参数验证
	if err := h.validateCreateFolder(cmd); err != nil {
		return herrors.NewBadReqError(err.Error())
	}

	// 2. 构建文件夹对象
	folder := &model.Folder{
		Name:      cmd.Name,
		ParentID:  cmd.ParentID,
		CreatedBy: createdBy,
		TenantID:  tenantID,
	}

	// 3. 调用服务创建文件夹
	if err := h.service.CreateFolder(ctx, folder); err != nil {
		// 根据错误类型返回不同的错误
		switch err.Error() {
		case "folder not found":
			return herrors.NewBadReqError(err.Error())
		case "permission denied":
			return herrors.NewBadReqError(err.Error())
		case "folder already exists":
			return herrors.NewBadReqError(err.Error())
		default:
			return herrors.NewServerHError(err)
		}
	}

	return nil
}

// HandleUploadFile 处理文件上传
func (h *StorageCommandHandler) HandleUploadFile(ctx context.Context, cmd *commands.UploadFileCommand, tenantID, createdBy string) (*model.File, herrors.Herr) {
	// 1. 参数验证
	if err := h.validateUploadFile(ctx, cmd); err != nil {
		return nil, err
	}

	// 2. 如果没有指定文件夹ID，创建默认文件夹
	if cmd.FolderID == "" || cmd.FolderID == "0" {
		folder, err := h.service.CreateDefaultFolder(ctx, tenantID, createdBy)
		if err != nil {
			return nil, herrors.NewServerHError(err)
		}
		cmd.FolderID = folder.ID
	}

	// 4. 创建文件读取器
	reader := bytes.NewReader(cmd.File)

	// 5. 上传文件
	file, err := h.service.UploadFile(ctx, reader, cmd.FileName, int64(len(cmd.File)), cmd.FolderID, tenantID, createdBy)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	return file, nil
}

// HandleDeleteFile 处理删除文件
func (h *StorageCommandHandler) HandleDeleteFile(ctx context.Context, cmd *commands.DeleteFileCommand) herrors.Herr {
	if err := h.service.DeleteFile(ctx, cmd.ID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleDeleteFolder 处理删除文件夹
func (h *StorageCommandHandler) HandleDeleteFolder(ctx context.Context, cmd *commands.DeleteFolderCommand) herrors.Herr {
	if err := h.service.DeleteFolder(ctx, cmd.ID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleRenameFolder 处理重命名文件夹
func (h *StorageCommandHandler) HandleRenameFolder(ctx context.Context, cmd *commands.RenameFolderCommand) herrors.Herr {
	if err := h.service.RenameFolder(ctx, cmd.ID, cmd.Name); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleMoveFile 处理移动文件
func (h *StorageCommandHandler) HandleMoveFile(ctx context.Context, cmd *commands.MoveFileCommand) herrors.Herr {
	if err := h.service.MoveFile(ctx, cmd.ID, cmd.FolderID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleMoveFolder 处理移动文件夹
func (h *StorageCommandHandler) HandleMoveFolder(ctx context.Context, cmd *commands.MoveFolderCommand) herrors.Herr {
	if err := h.service.MoveFolder(ctx, cmd.ID, cmd.TargetParentID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// validateCreateFolder 验证创建文件夹参数
func (h *StorageCommandHandler) validateCreateFolder(cmd *commands.CreateFolderCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if cmd.Name == "" {
		return herrors.NewBadReqError("folder name is empty")
	}
	if cmd.ParentID == "" {
		cmd.ParentID = "0"
	}
	return nil
}

// validateMoveFolder 验证移动文件夹参数
func (h *StorageCommandHandler) validateMoveFolder(cmd *commands.MoveFolderCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if cmd.ID == "" {
		return herrors.NewBadReqError("invalid folder id")
	}
	if cmd.TargetParentID == "" {
		return herrors.NewBadReqError("invalid parent folder id")
	}
	if cmd.ID == cmd.TargetParentID {
		return herrors.NewBadReqError("cannot move folder to itself")
	}
	return nil
}

// validateUploadFile 验证上传文件参数
func (h *StorageCommandHandler) validateUploadFile(ctx context.Context, cmd *commands.UploadFileCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if len(cmd.File) == 0 {
		return herrors.NewBadReqError("file is empty")
	}
	if cmd.FileName == "" {
		return herrors.NewBadReqError("file name is empty")
	}
	if cmd.FolderID != "" && cmd.FolderID != "0" {
		// 验证文件夹是否存在
		_, err := h.service.GetFolder(ctx, cmd.FolderID)
		if err != nil {
			return herrors.NewBadReqError("invalid folder id")
		}
	}
	return nil
}

// HandleBatchDeleteFiles 处理批量删除文件
func (h *StorageCommandHandler) HandleBatchDeleteFiles(ctx context.Context, cmd *commands.BatchDeleteFilesCommand) herrors.Herr {
	if err := h.validateBatchDeleteFiles(cmd); err != nil {
		return herrors.NewBadReqError(err.Error())
	}
	if err := h.service.BatchDeleteFiles(ctx, cmd.FileIDs); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleBatchMoveFiles 处理批量移动文件
func (h *StorageCommandHandler) HandleBatchMoveFiles(ctx context.Context, cmd *commands.BatchMoveFilesCommand) herrors.Herr {
	if err := h.validateBatchMoveFiles(cmd); err != nil {
		return herrors.NewBadReqError(err.Error())
	}
	if err := h.service.BatchMoveFiles(ctx, cmd.FileIDs, cmd.TargetFolderID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleShareFile 处理分享文件
func (h *StorageCommandHandler) HandleShareFile(ctx context.Context, cmd *commands.ShareFileCommand) (*dto.FileShareDto, herrors.Herr) {
	if err := h.validateShareFile(cmd); err != nil {
		return nil, herrors.NewBadReqError(err.Error())
	}
	share, err := h.service.ShareFile(ctx, cmd.FileID, cmd.ExpireTime, cmd.Password, cmd.CreatedBy)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	shareDto := &dto.FileShareDto{
		ID:         share.ID,
		FileID:     share.FileID,
		ShareCode:  share.ShareCode,
		ExpireTime: share.ExpireTime,
		CreatedBy:  share.CreatedBy,
		CreatedAt:  share.CreatedAt,
	}

	return shareDto, nil
}

// HandleRecycleFile 处理回收文件
func (h *StorageCommandHandler) HandleRecycleFile(ctx context.Context, cmd *commands.RecycleFileCommand) herrors.Herr {
	if err := h.service.RecycleFile(ctx, cmd.ID, cmd.DeletedBy); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// HandleRestoreFile 处理恢复文件
func (h *StorageCommandHandler) HandleRestoreFile(ctx context.Context, cmd *commands.RestoreFileCommand) herrors.Herr {
	if err := h.service.RestoreFile(ctx, cmd.ID); err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// validateBatchDeleteFiles 验证批量删除文件参数
func (h *StorageCommandHandler) validateBatchDeleteFiles(cmd *commands.BatchDeleteFilesCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if len(cmd.FileIDs) == 0 {
		return herrors.NewBadReqError("file ids is empty")
	}
	return nil
}

// validateBatchMoveFiles 验证批量移动文件参数
func (h *StorageCommandHandler) validateBatchMoveFiles(cmd *commands.BatchMoveFilesCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if len(cmd.FileIDs) == 0 {
		return herrors.NewBadReqError("file ids is empty")
	}
	if cmd.TargetFolderID == "" {
		return herrors.NewBadReqError("invalid target folder id")
	}
	return nil
}

// validateShareFile 验证分享文件参数
func (h *StorageCommandHandler) validateShareFile(cmd *commands.ShareFileCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if cmd.FileID == "" {
		return herrors.NewBadReqError("invalid file id")
	}
	if cmd.ExpireTime <= 0 {
		return herrors.NewBadReqError("invalid expire time")
	}
	return nil
}

// validateRecycleFile 验证回收文件参数
func (h *StorageCommandHandler) validateRecycleFile(cmd *commands.RecycleFileCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if cmd.ID == "" {
		return herrors.NewBadReqError("invalid file id")
	}
	return nil
}

// validateRestoreFile 验证恢复文件参数
func (h *StorageCommandHandler) validateRestoreFile(cmd *commands.RestoreFileCommand) herrors.Herr {
	if cmd == nil {
		return herrors.NewBadReqError("command is nil")
	}
	if cmd.ID == "" {
		return herrors.NewBadReqError("invalid file id")
	}
	return nil
}
