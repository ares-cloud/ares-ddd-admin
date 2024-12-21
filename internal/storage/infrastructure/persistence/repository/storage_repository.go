package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type storageRepository struct {
	db     database.IDataBase
	mapper *mapper.StorageMapper
}

func NewStorageRepository(db database.IDataBase) repository.IStorageRepository {
	// 同步表
	if err := db.DB(context.Background()).AutoMigrate(&entity.Folder{}, &entity.File{}, &entity.FileShare{}); err != nil {
		hlog.Fatalf("sync sys storage tables to db error: %v", err)
	}
	return &storageRepository{
		db:     db,
		mapper: &mapper.StorageMapper{},
	}
}

// genStringId 生成字符串ID
func (r *storageRepository) genStringId() string {
	return r.db.GenStringId()
}

// CreateFile 创建文件
func (r *storageRepository) CreateFile(ctx context.Context, file *model.File) error {
	// 生成ID
	file.ID = r.genStringId()
	return r.db.DB(ctx).Create(r.mapper.ToFileEntity(file)).Error
}

// UpdateFile 更新文件
func (r *storageRepository) UpdateFile(ctx context.Context, file *model.File) error {
	return r.db.DB(ctx).Updates(r.mapper.ToFileEntity(file)).Error
}

// DeleteFile 删除文件
func (r *storageRepository) DeleteFile(ctx context.Context, id string) error {
	return r.db.DB(ctx).Delete(&entity.File{}, id).Error
}

// GetFile 获取文件
func (r *storageRepository) GetFile(ctx context.Context, id string) (*model.File, error) {
	var ent entity.File
	if err := r.db.DB(ctx).First(&ent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToFileDomain(&ent), nil
}

// ListFiles 查询文件列表
func (r *storageRepository) ListFiles(ctx context.Context, folderID string, qb *query.QueryBuilder) ([]*model.File, int64, error) {
	db := r.db.DB(ctx).Model(&entity.File{})

	// 添加文件夹条件
	if folderID != "0" {
		db = db.Where("folder_id = ?", folderID)
	}

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序和分页
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	}
	if limit, offset := qb.BuildLimit(); limit != "" {
		db = db.Limit(offset[1]).Offset(offset[0])
	}

	// 查询数据
	var entities []*entity.File
	if err := db.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return r.mapper.ToFileDomainList(entities), total, nil
}

// CreateFolder 创建文件夹
func (r *storageRepository) CreateFolder(ctx context.Context, folder *model.Folder) error {
	// 生成ID
	folder.ID = r.genStringId()
	return r.db.DB(ctx).Create(r.mapper.ToFolderEntity(folder)).Error
}

// UpdateFolder 更新文件夹
func (r *storageRepository) UpdateFolder(ctx context.Context, folder *model.Folder) error {
	return r.db.DB(ctx).Updates(r.mapper.ToFolderEntity(folder)).Error
}

// DeleteFolder 删除文件夹
func (r *storageRepository) DeleteFolder(ctx context.Context, id string) error {
	return r.db.DB(ctx).Delete(&entity.Folder{}, id).Error
}

// GetFolder 获取文件夹
func (r *storageRepository) GetFolder(ctx context.Context, id string) (*model.Folder, error) {
	var ent entity.Folder
	if err := r.db.DB(ctx).First(&ent, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get folder error: %v", err)
	}
	return r.mapper.ToFolderDomain(&ent), nil
}

// ListFolders 查询文件夹列表
func (r *storageRepository) ListFolders(ctx context.Context, parentID string, qb *query.QueryBuilder) ([]*model.Folder, int64, error) {
	db := r.db.DB(ctx).Model(&entity.Folder{})

	// 添加父文件夹条件
	if parentID != "0" {
		db = db.Where("parent_id = ?", parentID)
	}

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count folders error: %v", err)
	}

	// 添加排序和分页
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	}
	if limit, offset := qb.BuildLimit(); limit != "" {
		db = db.Limit(offset[1]).Offset(offset[0])
	}

	// 查询数据
	var entities []*entity.Folder
	if err := db.Find(&entities).Error; err != nil {
		return nil, 0, fmt.Errorf("list folders error: %v", err)
	}

	return r.mapper.ToFolderDomainList(entities), total, nil
}

// GetExpiredRecycleFiles 获取过期的回收站文件
func (r *storageRepository) GetExpiredRecycleFiles(ctx context.Context, expireTime time.Time) ([]*model.File, error) {
	var entities []*entity.File

	// 查询条件：
	// 1. 已删除的文件
	// 2. 删除时间早于过期时间
	err := r.db.DB(ctx).
		Where("is_deleted = ?", true).
		Where("deleted_at < ?", expireTime.Unix()).
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return r.mapper.ToFileDomainList(entities), nil
}

// GetFileShare 获取文件分享信息
func (r *storageRepository) GetFileShare(ctx context.Context, shareCode string) (*model.FileShare, error) {
	var ent entity.FileShare
	if err := r.db.DB(ctx).Where("share_code = ?", shareCode).First(&ent).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToFileShareDomain(&ent), nil
}

// CreateFileShare 创建文件分享记录
func (r *storageRepository) CreateFileShare(ctx context.Context, share *model.FileShare) error {
	// 生成ID
	share.ID = r.genStringId()
	return r.db.DB(ctx).Create(r.mapper.ToFileShareEntity(share)).Error
}

// GetFolderTree 获取文件夹树形结构
func (r *storageRepository) GetFolderTree(ctx context.Context, tenantID string) ([]*model.FolderTree, error) {
	// 1. 获取所有文件夹
	qb := query.NewQueryBuilder().Where("tenant_id", query.Eq, tenantID)
	folders, _, err := r.ListFolders(ctx, "0", qb)
	if err != nil {
		return nil, err
	}

	// 2. 构建文件夹映射
	folderMap := make(map[string]*model.FolderTree)
	var roots []*model.FolderTree

	// 3. 构建树形结构
	for _, folder := range folders {
		node := &model.FolderTree{
			Folder:   folder,
			Children: make([]*model.FolderTree, 0),
		}
		folderMap[folder.ID] = node

		if folder.ParentID == "0" {
			roots = append(roots, node)
		} else {
			if parent, ok := folderMap[folder.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots, nil
}
