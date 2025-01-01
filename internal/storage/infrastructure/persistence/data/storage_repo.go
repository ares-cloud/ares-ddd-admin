package data

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/repository"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type StorageRepo struct {
	db database.IDataBase
}

func NewStorageRepo(db database.IDataBase) repository.IStorageRepos {
	return &StorageRepo{db: db}
}

// CreateFile 创建文件
func (r *StorageRepo) CreateFile(ctx context.Context, file *entity.File) error {
	return r.db.DB(ctx).Create(file).Error
}

// UpdateFile 更新文件
func (r *StorageRepo) UpdateFile(ctx context.Context, file *entity.File) error {
	return r.db.DB(ctx).Save(file).Error
}

// DeleteFile 删除文件
func (r *StorageRepo) DeleteFile(ctx context.Context, id string) error {
	return r.db.DB(ctx).Delete(&entity.File{}, id).Error
}

// GetFile 获取文件
func (r *StorageRepo) GetFile(ctx context.Context, id string) (*entity.File, error) {
	var file entity.File
	if err := r.db.DB(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// ListFiles 查询文件列表
func (r *StorageRepo) ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*entity.File, int64, error) {
	db := r.db.DB(ctx).Model(&entity.File{})

	// 添加文件夹条件
	if folderID != "0" && folderID != "" {
		db = db.Where("folder_id = ?", folderID)
	}

	// 添加查询条件
	if qb != nil {
		db = qb.Build(db)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	var files []*entity.File
	if err := db.Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// CreateFolder 创建文件夹
func (r *StorageRepo) CreateFolder(ctx context.Context, folder *entity.Folder) error {
	return r.db.DB(ctx).Create(folder).Error
}

// UpdateFolder 更新文件夹
func (r *StorageRepo) UpdateFolder(ctx context.Context, folder *entity.Folder) error {
	return r.db.DB(ctx).Updates(folder).Error
}

// DeleteFolder 删除文件夹
func (r *StorageRepo) DeleteFolder(ctx context.Context, id string) error {
	return r.db.DB(ctx).Delete(&entity.Folder{}, id).Error
}

// GetFolder 获取文件夹
func (r *StorageRepo) GetFolder(ctx context.Context, id string) (*entity.Folder, error) {
	var folder entity.Folder
	if err := r.db.DB(ctx).First(&folder, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get folder error: %v", err)
	}
	return &folder, nil
}

// ListFolders 查询文件夹列表
func (r *StorageRepo) ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*entity.Folder, int64, error) {
	db := r.db.DB(ctx).Model(&entity.Folder{})

	// 添加父文件夹条件
	if parentID != "" {
		db = db.Where("parent_id = ?", parentID)
	}

	// 添加查询条件
	if qb != nil {
		db = qb.Build(db)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	var folders []*entity.Folder
	if err := db.Find(&folders).Error; err != nil {
		return nil, 0, err
	}

	return folders, total, nil
}

// GetFileShare 获取文件分享信息
func (r *StorageRepo) GetFileShare(ctx context.Context, shareCode string) (*entity.FileShare, error) {
	var share entity.FileShare
	if err := r.db.DB(ctx).Where("share_code = ?", shareCode).First(&share).Error; err != nil {
		return nil, err
	}
	return &share, nil
}

// CreateFileShare 创建文件分享记录
func (r *StorageRepo) CreateFileShare(ctx context.Context, share *entity.FileShare) error {
	return r.db.DB(ctx).Create(share).Error
}

// GetExpiredRecycleFiles 获取过期的回收站文件
func (r *StorageRepo) GetExpiredRecycleFiles(ctx context.Context, expireTime time.Time) ([]*entity.File, error) {
	var files []*entity.File
	err := r.db.DB(ctx).
		Where("is_deleted = ?", true).
		Where("deleted_at < ?", expireTime.Unix()).
		Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
