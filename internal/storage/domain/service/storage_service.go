package service

import (
	"context"
	"io"
	"math/rand"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type StorageService struct {
	repo    repository.IStorageRepository
	storage storage.StorageFactory
}

func NewStorageService(repo repository.IStorageRepository, storage storage.StorageFactory) *StorageService {
	return &StorageService{repo: repo, storage: storage}
}

// UploadFile 上传文件
func (s *StorageService) UploadFile(ctx context.Context, reader io.Reader, fileName string, size int64, folderID, tenantID, createdBy string) (*model.File, herrors.Herr) {
	// 1. 检查文件夹是否存在
	if folderID != "" && folderID != "0" {
		folder, err := s.repo.GetFolder(ctx, folderID)
		if err != nil {
			return nil, errors.StorageError(err)
		}
		if folder == nil {
			return nil, errors.FolderNotFound(folderID)
		}
	}

	// 2. 获取当前存储实例
	storage, err := s.storage.GetCurrentStorage()
	if err != nil {
		return nil, errors.StorageError(err)
	}

	// 3. 构建文件对象
	file := &model.File{
		Name:        fileName,
		Size:        size,
		FolderID:    folderID,
		TenantID:    tenantID,
		CreatedBy:   createdBy,
		StorageType: model.StorageTypeLocal, // 使用当前配置的存储类型
		CreatedAt:   time.Now().Unix(),
	}

	// 4. 验证文件属性
	if err := file.Validate(); err != nil {
		return nil, errors.InvalidFileName(err.Error())
	}

	// 5. 上传文件到存储
	uploadedFile, err := storage.Upload(ctx, reader, fileName, size, folderID)
	if err != nil {
		return nil, errors.StorageError(err)
	}

	// 6. 更新文件信息
	file.Path = uploadedFile.Path
	file.URL = uploadedFile.URL
	file.StorageType = uploadedFile.StorageType

	// 7. 保存到数据库
	savedFile, err := s.repo.SaveFile(ctx, file, nil)
	if err != nil {
		// 删除已上传的文件
		_ = storage.Delete(ctx, uploadedFile)
		return nil, errors.StorageError(err)
	}

	return savedFile, nil
}

// CreateFolder 创建文件夹
func (s *StorageService) CreateFolder(ctx context.Context, folder *model.Folder) herrors.Herr {
	// 1. 验证文件夹属性
	if err := folder.Validate(); err != nil {
		return herrors.NewBadReqError(err.Error())
	}

	// 2. 检查父文件夹是否存在
	if folder.ParentID != "" && folder.ParentID != "0" {
		parent, err := s.repo.GetFolder(ctx, folder.ParentID)
		if err != nil {
			return errors.StorageError(err)
		}
		if parent == nil {
			return errors.FolderNotFound(folder.ParentID)
		}
	}

	// 3. 创建文件夹
	if err := s.repo.SaveFolder(ctx, folder); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// DeleteFile 删除文件
func (s *StorageService) DeleteFile(ctx context.Context, id string) herrors.Herr {
	// 1. 检查文件是否存在
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if file == nil {
		return errors.FileNotFound(id)
	}

	// 2. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return errors.StorageError(err)
	}

	// 3. 从存储中删除文件
	if err := storage.Delete(ctx, file); err != nil {
		return errors.StorageError(err)
	}

	// 4. 从数据库中删除记录
	if err := s.repo.DeleteFile(ctx, id); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// GetFolder 获取文件夹
func (s *StorageService) GetFolder(ctx context.Context, id string) (*model.Folder, herrors.Herr) {
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return nil, errors.StorageError(err)
	}
	if folder == nil {
		return nil, errors.FolderNotFound(id)
	}
	return folder, nil
}

// MoveFile 移动文件
func (s *StorageService) MoveFile(ctx context.Context, fileID string, targetFolderID string) herrors.Herr {
	// 1. 检查文件是否存在
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return errors.StorageError(err)
	}
	if file == nil {
		return errors.FileNotFound(fileID)
	}

	// 2. 检查目标文件夹是否存在
	var targetFolder *model.Folder
	if targetFolderID != "0" {
		targetFolder, err = s.repo.GetFolder(ctx, targetFolderID)
		if err != nil {
			return errors.StorageError(err)
		}
		if targetFolder == nil {
			return errors.FolderNotFound(targetFolderID)
		}
	}

	// 3. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return errors.StorageError(err)
	}

	// 4. 构建新路径
	newPath := ""
	if targetFolder != nil {
		newPath = targetFolder.Path + "/" + file.Name
	} else {
		newPath = "/" + file.Name
	}

	// 5. 在存储中移动文件
	if err := storage.Move(ctx, file, newPath); err != nil {
		return errors.StorageError(err)
	}

	// 6. 更新数据库记录
	file.FolderID = targetFolderID
	file.Path = newPath
	if _, err := s.repo.SaveFile(ctx, file, nil); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// CreateDefaultFolder 创建默认文件夹
func (s *StorageService) CreateDefaultFolder(ctx context.Context, tenantID, createdBy string) (*model.Folder, herrors.Herr) {
	// 1. 构建默认文件夹对象
	folder := &model.Folder{
		Name:      "默认文件夹",
		ParentID:  "0",
		TenantID:  tenantID,
		CreatedBy: createdBy,
		CreatedAt: time.Now().Unix(),
	}

	// 2. 验证文件夹属性
	if err := folder.Validate(); err != nil {
		return nil, herrors.NewBadReqError(err.Error())
	}

	// 3. 创建文件夹
	if err := s.repo.SaveFolder(ctx, folder); err != nil {
		return nil, errors.StorageError(err)
	}

	return folder, nil
}

// DeleteFolder 删除文件夹
func (s *StorageService) DeleteFolder(ctx context.Context, id string) herrors.Herr {
	// 1. 检查文件夹是否存在
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if folder == nil {
		return errors.FolderNotFound(id)
	}

	// 2. 检查文件夹是否为空
	files, total, err := s.repo.ListFiles(ctx, id, db_query.NewQueryBuilder())
	if err != nil {
		return errors.StorageError(err)
	}
	if total > 0 {
		return errors.FolderNotEmpty(id, len(files))
	}

	// 3. 检查是否有子文件夹
	subFolders, total, err := s.repo.ListFolders(ctx, id, db_query.NewQueryBuilder())
	if err != nil {
		return errors.StorageError(err)
	}
	if total > 0 {
		return errors.FolderNotEmpty(id, len(subFolders))
	}

	// 4. 删除文件夹
	if err := s.repo.DeleteFolder(ctx, id); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// RenameFolder 重命名文件夹
func (s *StorageService) RenameFolder(ctx context.Context, id string, newName string) herrors.Herr {
	// 1. 检查文件夹是否存在
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if folder == nil {
		return errors.FolderNotFound(id)
	}

	// 2. 更新文件夹名
	if err := folder.SetName(newName); err != nil {
		return errors.InvalidFolderName(newName, err.Error())
	}

	if err := s.repo.SaveFolder(ctx, folder); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// MoveFolder 移动文件夹
func (s *StorageService) MoveFolder(ctx context.Context, id string, targetParentID string) herrors.Herr {
	// 1. 检查文件夹是否存在
	folder, err := s.repo.GetFolder(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if folder == nil {
		return errors.FolderNotFound(id)
	}

	// 2. 检查目标父文件夹是否存在
	if targetParentID != "0" {
		parent, err := s.repo.GetFolder(ctx, targetParentID)
		if err != nil {
			return errors.StorageError(err)
		}
		if parent == nil {
			return errors.FolderNotFound(targetParentID)
		}
		// 检查是否移动到自己的子文件夹
		if isSubFolder(ctx, s.repo, targetParentID, id) {
			return errors.InvalidMoveOperation("cannot move folder to its subfolder")
		}
	}

	// 3. 更新文件夹
	folder.ParentID = targetParentID
	folder.UpdatedAt = time.Now().Unix()
	if err := s.repo.SaveFolder(ctx, folder); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// BatchDeleteFiles 批量删除文件
func (s *StorageService) BatchDeleteFiles(ctx context.Context, fileIDs []string) herrors.Herr {
	for _, id := range fileIDs {
		if err := s.DeleteFile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// BatchMoveFiles 批量移动文件
func (s *StorageService) BatchMoveFiles(ctx context.Context, fileIDs []string, targetFolderID string) herrors.Herr {
	for _, id := range fileIDs {
		if err := s.MoveFile(ctx, id, targetFolderID); err != nil {
			return err
		}
	}
	return nil
}

// ShareFile 分享文件
func (s *StorageService) ShareFile(ctx context.Context, fileID string, expireTime int64, password string, createdBy string) (*model.FileShare, herrors.Herr) {
	// 1. 检查文件是否存在
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return nil, errors.StorageError(err)
	}
	if file == nil {
		return nil, errors.FileNotFound(fileID)
	}

	// 2. 创建分享记录
	share := &model.FileShare{
		FileID:     fileID,
		ShareCode:  generateShareCode(), // 生成分享码
		Password:   password,
		ExpireTime: expireTime,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now().Unix(),
	}

	// 3. 验证分享属性
	if err := share.Validate(); err != nil {
		return nil, herrors.NewBadReqError(err.Error())
	}

	// 4. 保存分享记录
	if err := s.repo.CreateFileShare(ctx, share); err != nil {
		return nil, errors.StorageError(err)
	}

	return share, nil
}

// RecycleFile 回收文件
func (s *StorageService) RecycleFile(ctx context.Context, id string, deletedBy string) herrors.Herr {
	// 1. 检查文件是否存在
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if file == nil {
		return errors.FileNotFound(id)
	}

	// 2. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return errors.StorageError(err)
	}

	// 3. 构建回收站路径
	recyclePath := "/recycle/" + file.Name

	// 4. 在存储中移动文件到回收站
	if err := storage.Move(ctx, file, recyclePath); err != nil {
		return errors.StorageError(err)
	}

	// 5. 更新文件状态
	file.IsDeleted = true
	file.DeletedBy = deletedBy
	file.DeletedAt = time.Now().Unix()
	file.OriginalPath = file.Path
	file.Path = recyclePath

	// 6. 更新数据库记录
	if _, err := s.repo.SaveFile(ctx, file, nil); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// RestoreFile 恢复文件
func (s *StorageService) RestoreFile(ctx context.Context, id string) herrors.Herr {
	// 1. 检查文件是否存在
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return errors.StorageError(err)
	}
	if file == nil {
		return errors.FileNotFound(id)
	}

	// 2. 检查文件是否在回收站
	if !file.IsDeleted {
		return errors.InvalidOperation("file is not in recycle bin")
	}

	// 3. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return errors.StorageError(err)
	}

	// 4. 在存储中恢复文件
	if err := storage.Move(ctx, file, file.OriginalPath); err != nil {
		return errors.StorageError(err)
	}

	// 5. 恢复文件状态
	file.IsDeleted = false
	file.DeletedBy = ""
	file.DeletedAt = 0
	file.Path = file.OriginalPath
	file.OriginalPath = ""

	// 6. 更新数据库记录
	if _, err := s.repo.SaveFile(ctx, file, nil); err != nil {
		return errors.StorageError(err)
	}

	return nil
}

// PreviewFile 预览文件
func (s *StorageService) PreviewFile(ctx context.Context, fileID string) (string, herrors.Herr) {
	// 1. 获取文件信息
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return "", errors.StorageError(err)
	}
	if file == nil {
		return "", errors.FileNotFound(fileID)
	}

	// 2. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return "", errors.StorageError(err)
	}

	// 3. 获取预览URL
	previewURL, err := storage.GetPreviewURL(ctx, file)
	if err != nil {
		return "", errors.StorageError(err)
	}

	return previewURL, nil
}

// DownloadFile 下载文件
func (s *StorageService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, string, herrors.Herr) {
	// 1. 获取文件信息
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return nil, "", errors.StorageError(err)
	}
	if file == nil {
		return nil, "", errors.FileNotFound(fileID)
	}

	// 2. 获取存储实例
	storage, err := s.storage.GetStorage(file.StorageType)
	if err != nil {
		return nil, "", errors.StorageError(err)
	}

	// 3. 下载文件
	reader, err := storage.Download(ctx, file)
	if err != nil {
		return nil, "", errors.StorageError(err)
	}

	return reader, file.Name, nil
}

// 辅助函数

// isSubFolder 检查是否是子文件夹
func isSubFolder(ctx context.Context, repo repository.IStorageRepository, parentID, childID string) bool {
	if parentID == childID {
		return true
	}

	folder, err := repo.GetFolder(ctx, parentID)
	if err != nil || folder == nil {
		return false
	}

	if folder.ParentID == "0" {
		return false
	}

	return isSubFolder(ctx, repo, folder.ParentID, childID)
}

// generateShareCode 生成分享码
func generateShareCode() string {
	// 生成8位随机字符串
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
