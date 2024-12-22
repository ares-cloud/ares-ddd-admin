package service

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/google/uuid"
)

type StorageService struct {
	repo           repository.IStorageRepository
	storageFactory storage.StorageFactory
}

func NewStorageService(repo repository.IStorageRepository, factory storage.StorageFactory) *StorageService {
	return &StorageService{
		repo:           repo,
		storageFactory: factory,
	}
}

// GetDefaultStorageType 获取默认存储类型
func (s *StorageService) GetDefaultStorageType(ctx context.Context) (model.StorageType, error) {
	// 从配置中获取默认存储类型
	return model.StorageTypeMinio, nil // 这里可以根据实际配置返回
}

// CreateDefaultFolder 创建默认文件夹(按月份)
func (s *StorageService) CreateDefaultFolder(ctx context.Context, tenantID, createdBy string) (*model.Folder, error) {
	now := time.Now()
	yearStr := now.Format("2006")
	monthStr := now.Format("01")

	// 1. 获取或创建年份文件夹
	yearFolder, err := s.getOrCreateFolder(ctx, "0", yearStr, tenantID, createdBy)
	if err != nil {
		return nil, fmt.Errorf("create year folder error: %v", err)
	}

	// 2. 获取或创建月份文件夹
	monthFolder, err := s.getOrCreateFolder(ctx, yearFolder.ID, monthStr, tenantID, createdBy)
	if err != nil {
		return nil, fmt.Errorf("create month folder error: %v", err)
	}

	return monthFolder, nil
}

// getOrCreateFolder 获取或创建文件夹
func (s *StorageService) getOrCreateFolder(ctx context.Context, parentID string, name, tenantID, createdBy string) (*model.Folder, error) {
	// 1. 查询是否存在
	folders, total, err := s.repo.ListFolders(ctx, parentID, query.NewQueryBuilder().
		Where("name", query.Eq, name).
		Where("tenant_id", query.Eq, tenantID))
	if err != nil {
		return nil, err
	}

	// 2. 已存在则直接返回
	if total > 0 {
		return folders[0], nil
	}

	// 3. 不存在则创建
	folder := &model.Folder{
		Name:      name,
		ParentID:  parentID,
		TenantID:  tenantID,
		CreatedBy: createdBy,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// 4. 设置完整路径
	if parentID != "0" {
		parent, err := s.repo.GetFolder(ctx, parentID)
		if err != nil {
			return nil, err
		}
		folder.Path = path.Join(parent.Path, name)
	} else {
		folder.Path = name
	}

	// 5. 保存到数据库
	if err := s.repo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// generateUniqueFileName 生成唯一文件名
func (s *StorageService) generateUniqueFileName(originalName string) string {
	ext := path.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// UploadFile 上传文件
func (s *StorageService) UploadFile(ctx context.Context, file io.Reader, filename string, size int64, folderID string, tenantID, createdBy string) (*model.File, error) {
	// 如果没有指定文件夹，创建默认文件夹
	if folderID == "0" {
		folder, err := s.CreateDefaultFolder(ctx, tenantID, createdBy)
		if err != nil {
			return nil, err
		}
		folderID = folder.ID
	}

	// 获取存储实现
	store, err := s.storageFactory.GetCurrentStorage()
	if err != nil {
		return nil, err
	}

	// 获取文件夹信息(用于构建存储路径)
	folder, err := s.repo.GetFolder(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// 生成唯一文件名
	generatedName := s.generateUniqueFileName(filename)

	// 上传文件(使用生成的文件名)
	fileModel, err := store.Upload(ctx, file, generatedName, size, folder.Path)
	if err != nil {
		return nil, err
	}

	// 设置文件信息
	fileModel.Name = filename               // 保存原始文件名
	fileModel.GeneratedName = generatedName // 保存生成的文件名
	fileModel.FolderID = folderID
	fileModel.Type = filepath.Ext(filename)
	fileModel.CreatedBy = createdBy
	fileModel.TenantID = tenantID

	// 保存到数据库
	if err := s.repo.CreateFile(ctx, fileModel); err != nil {
		// 删除已上传的文件
		_ = store.Delete(ctx, fileModel)
		return nil, err
	}

	return fileModel, nil
}

// DeleteFile 删除文件
func (s *StorageService) DeleteFile(ctx context.Context, id string) error {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return err
	}

	// 获取存储实现
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return err
	}

	// 删除存储中的文件
	if err := store.Delete(ctx, file); err != nil {
		return err
	}

	// 从数据库删除
	return s.repo.DeleteFile(ctx, id)
}

// DeleteFolder 删除文件夹
func (s *StorageService) DeleteFolder(ctx context.Context, id string) error {
	// 获取文件夹下的所有文件
	files, _, err := s.repo.ListFiles(ctx, id, query.NewQueryBuilder())
	if err != nil {
		return err
	}

	// 删除所有文件
	for _, file := range files {
		if err := s.DeleteFile(ctx, file.ID); err != nil {
			return err
		}
	}

	// 获取文件夹
	folders, _, err := s.repo.ListFolders(ctx, id, query.NewQueryBuilder())
	if err != nil {
		return err
	}

	// 递归删除子文件夹
	for _, folder := range folders {
		if err := s.DeleteFolder(ctx, folder.ID); err != nil {
			return err
		}
	}

	// 删除文件夹
	return s.repo.DeleteFolder(ctx, id)
}

// RenameFolder 重命名文件夹
func (s *StorageService) RenameFolder(ctx context.Context, id string, newName string) error {
	// 获取文件夹信息
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return err
	}

	// 更新名称
	folder.Name = newName
	folder.UpdatedAt = time.Now().Unix()

	// 更新数据库
	return s.repo.UpdateFolder(ctx, folder)
}

// MoveFile 移动文件
func (s *StorageService) MoveFile(ctx context.Context, id string, targetFolderID string) error {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return err
	}

	// 获取目标文件夹信息
	targetFolder, err := s.repo.GetFolder(ctx, targetFolderID)
	if err != nil {
		return err
	}

	// 获取存储实现
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return err
	}

	// 更新文件路径
	oldPath := file.Path
	file.FolderID = targetFolderID
	file.Path = path.Join(targetFolder.Path, file.Name)
	file.UpdatedAt = time.Now().Unix()

	// 移动存储中的文件
	if err := store.Move(ctx, file, oldPath); err != nil {
		return err
	}

	// 更新数据库
	return s.repo.UpdateFile(ctx, file)
}

// BatchDeleteFiles 批量删除文件
func (s *StorageService) BatchDeleteFiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := s.DeleteFile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// BatchMoveFiles 批量移动文件
func (s *StorageService) BatchMoveFiles(ctx context.Context, ids []string, targetFolderID string) error {
	for _, id := range ids {
		if err := s.MoveFile(ctx, id, targetFolderID); err != nil {
			return err
		}
	}
	return nil
}

// ShareFile 分享文件
func (s *StorageService) ShareFile(ctx context.Context, id string, expireTime int64, password string, createdBy string) (*model.FileShare, error) {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}

	// 生成分享码
	shareCode := generateShareCode()

	// 创建分享记录
	share := &model.FileShare{
		FileID:     file.ID,
		ShareCode:  shareCode,
		Password:   password,
		ExpireTime: time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
		CreatedAt:  time.Now().Unix(),
		CreatedBy:  createdBy,
	}

	// 保存分享记录
	if err := s.repo.CreateFileShare(ctx, share); err != nil {
		return nil, err
	}

	return share, nil
}

// RecycleFile 回收文件
func (s *StorageService) RecycleFile(ctx context.Context, id string, deletedBy string) error {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return err
	}

	// 更新文件状态
	file.IsDeleted = true
	file.DeletedAt = time.Now().Unix()
	file.DeletedBy = deletedBy
	file.OriginalPath = file.Path

	// 移动到回收站目录
	recyclePath := path.Join("recycle", file.Path)
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return err
	}

	// 移动文件
	if err := store.Move(ctx, file, recyclePath); err != nil {
		return err
	}

	// 更新数据库
	return s.repo.UpdateFile(ctx, file)
}

// RestoreFile 恢复文件
func (s *StorageService) RestoreFile(ctx context.Context, id string) error {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	if !file.IsDeleted {
		return herrors.NewBadReqError("file is not in recycle bin")
	}

	// 恢复文件状态
	file.IsDeleted = false
	file.DeletedAt = time.Now().Unix()
	file.DeletedBy = ""
	file.Path = file.OriginalPath
	file.OriginalPath = ""

	// 从回收站恢复
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return err
	}

	if err := store.Move(ctx, file, file.Path); err != nil {
		return err
	}

	// 更新数据库
	return s.repo.UpdateFile(ctx, file)
}

// generateShareCode 生成分享码
func generateShareCode() string {
	return uuid.New().String()[:8]
}

// GetShareFile 获取分享文件
func (s *StorageService) GetShareFile(ctx context.Context, shareCode, password string) (*model.File, error) {
	// 获取分享信息
	share, err := s.repo.GetFileShare(ctx, shareCode)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	// 检查是否过期
	if time.Now().After(time.Unix(share.ExpireTime, 0)) {
		return nil, herrors.NewBadReqError("share link expired")
	}

	// 检查密码
	if share.Password != "" && share.Password != password {
		return nil, herrors.NewBadReqError("invalid password")
	}

	// 获取文件信息
	file, err := s.repo.GetFile(ctx, share.FileID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	return file, nil
}

// PreviewFile 预览文件
func (s *StorageService) PreviewFile(ctx context.Context, id string) (string, error) {
	// 获取文件信���
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return "", err
	}

	// 获取存储实现
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return "", err
	}

	// 获取预览URL
	return store.GetPreviewURL(ctx, file)
}

// DownloadFile 下载文件
func (s *StorageService) DownloadFile(ctx context.Context, id string) (io.ReadCloser, string, error) {
	// 获取文件信息
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return nil, "", err
	}

	// 获取存储实现
	store, err := s.storageFactory.GetStorage(file.StorageType)
	if err != nil {
		return nil, "", err
	}

	// 下载文件
	reader, err := store.Download(ctx, file)
	if err != nil {
		return nil, "", err
	}

	return reader, file.Name, nil
}

// CreateFolder 创建文件夹
func (s *StorageService) CreateFolder(ctx context.Context, folder *model.Folder) error {
	// 1. 检查父文件夹是否存在
	if folder.ParentID != "0" {
		parent, err := s.repo.GetFolder(ctx, folder.ParentID)
		if err != nil {
			return herrors.NewServerHError(err)
		}
		// 设置完整路径
		folder.Path = path.Join(parent.Path, folder.Name)
	} else {
		folder.Path = folder.Name
	}

	// 2. 检查同名文件夹是否存在
	qb := query.NewQueryBuilder().
		Where("parent_id", query.Eq, folder.ParentID).
		Where("name", query.Eq, folder.Name).
		Where("tenant_id", query.Eq, folder.TenantID)

	_, total, err := s.repo.ListFolders(ctx, folder.ParentID, qb)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if total > 0 {
		return herrors.NewBadReqError("folder already exists")
	}

	// 3. 设置创建时间
	now := time.Now().Unix()
	folder.CreatedAt = now
	folder.UpdatedAt = now

	// 4. 创建文件夹
	if err := s.repo.CreateFolder(ctx, folder); err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}

// MoveFolder 移动文件夹
func (s *StorageService) MoveFolder(ctx context.Context, id string, newParentID string) error {
	// 1. 获取文件夹信息
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 2. 检查新父文件夹是否存在
	newParent, err := s.repo.GetFolder(ctx, newParentID)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 3. 检查是否移动到自己的子文件夹
	if s.isSubFolder(ctx, folder.ID, newParentID) {
		return herrors.NewBadReqError("cannot move folder to its subfolder")
	}

	// 4. 更新文件夹路径
	oldPath := folder.Path
	folder.ParentID = newParentID
	folder.Path = path.Join(newParent.Path, folder.Name)
	folder.UpdatedAt = time.Now().Unix()

	// 5. 更新所有子文件夹和文件的路径
	if err := s.updateSubPaths(ctx, folder.ID, oldPath, folder.Path); err != nil {
		return herrors.NewServerHError(err)
	}

	// 6. 更新数据库
	return s.repo.UpdateFolder(ctx, folder)
}

// isSubFolder 检查是否是子文件夹
func (s *StorageService) isSubFolder(ctx context.Context, parentID, folderID string) bool {
	if parentID == folderID {
		return true
	}

	folder, err := s.repo.GetFolder(ctx, folderID)
	if err != nil {
		return false
	}

	if folder.ParentID == "0" {
		return false
	}

	return s.isSubFolder(ctx, parentID, folder.ParentID)
}

// updateSubPaths 更新子路径
func (s *StorageService) updateSubPaths(ctx context.Context, folderID string, oldPath, newPath string) error {
	// 1. 更新子文件夹路径
	folders, _, err := s.repo.ListFolders(ctx, folderID, query.NewQueryBuilder())
	if err != nil {
		return err
	}

	for _, folder := range folders {
		folder.Path = strings.Replace(folder.Path, oldPath, newPath, 1)
		if err := s.repo.UpdateFolder(ctx, folder); err != nil {
			return err
		}
		// 递归更新子文件夹
		if err := s.updateSubPaths(ctx, folder.ID, oldPath, newPath); err != nil {
			return err
		}
	}

	// 2. 更新文件路径
	files, _, err := s.repo.ListFiles(ctx, folderID, query.NewQueryBuilder())
	if err != nil {
		return err
	}

	for _, file := range files {
		file.Path = strings.Replace(file.Path, oldPath, newPath, 1)
		if err := s.repo.UpdateFile(ctx, file); err != nil {
			return err
		}
	}

	return nil
}

// GetFolderTree 获取文件夹树形结构
func (s *StorageService) GetFolderTree(ctx context.Context, tenantID string) ([]*model.FolderTree, error) {
	return s.repo.GetFolderTree(ctx, tenantID)
}

// validateCreateFolder 验证创建文件夹业务逻辑
func (s *StorageService) validateCreateFolder(ctx context.Context, folder *model.Folder) error {
	// 检查父文件夹是否存在
	if folder.ParentID != "0" {
		parent, err := s.repo.GetFolder(ctx, folder.ParentID)
		if err != nil {
			return herrors.NewBadReqError("parent folder not found")
		}
		if parent.TenantID != folder.TenantID {
			return herrors.NewBadReqError("parent folder belongs to different tenant")
		}
	}

	// 检查同名文件夹是否存在
	qb := query.NewQueryBuilder().
		Where("parent_id", query.Eq, folder.ParentID).
		Where("name", query.Eq, folder.Name).
		Where("tenant_id", query.Eq, folder.TenantID)

	_, total, err := s.repo.ListFolders(ctx, folder.ParentID, qb)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if total > 0 {
		return herrors.NewBadReqError("folder already exists")
	}

	return nil
}

// validateMoveFolder 验证移动文件夹业务逻辑
func (s *StorageService) validateMoveFolder(ctx context.Context, folderID, newParentID string) error {
	// 获取文件夹信息
	folder, err := s.repo.GetFolder(ctx, folderID)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 检查新父文件夹是否存在
	newParent, err := s.repo.GetFolder(ctx, newParentID)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 检查租户权限
	if folder.TenantID != newParent.TenantID {
		return fmt.Errorf("folders belong to different tenants")
	}

	// 检查是否移动到子文件夹
	if s.isSubFolder(ctx, folder.ID, newParentID) {
		return fmt.Errorf("cannot move folder to its child")
	}

	return nil
}

// GetFolder 获取文件夹信息
func (s *StorageService) GetFolder(ctx context.Context, id string) (*model.Folder, error) {
	// 1. 参数验证
	if id == "" {
		return nil, herrors.NewBadReqError("folder id is empty")
	}

	// 2. 从仓储获取文件夹
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, herrors.DataIsExist
		}
		return nil, herrors.NewServerHError(err)
	}

	return folder, nil
}
