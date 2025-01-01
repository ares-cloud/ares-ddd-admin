package repository

import (
	"context"
	"io"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IStorageRepos 存储仓储接口(数据库操作)
type IStorageRepos interface {
	// 文件相关
	CreateFile(ctx context.Context, file *entity.File) error
	UpdateFile(ctx context.Context, file *entity.File) error
	DeleteFile(ctx context.Context, id string) error
	GetFile(ctx context.Context, id string) (*entity.File, error)
	ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*entity.File, int64, error)

	// 文件夹相关
	CreateFolder(ctx context.Context, folder *entity.Folder) error
	UpdateFolder(ctx context.Context, folder *entity.Folder) error
	DeleteFolder(ctx context.Context, id string) error
	GetFolder(ctx context.Context, id string) (*entity.Folder, error)
	ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*entity.Folder, int64, error)

	// 文件分享相关
	GetFileShare(ctx context.Context, shareCode string) (*entity.FileShare, error)
	CreateFileShare(ctx context.Context, share *entity.FileShare) error

	// 其他
	GetExpiredRecycleFiles(ctx context.Context, expireTime time.Time) ([]*entity.File, error)
}

type storageRepository struct {
	db     database.IDataBase
	repo   IStorageRepos
	mapper *mapper.StorageMapper
}

func NewStorageRepository(db database.IDataBase, repo IStorageRepos) repository.IStorageRepository {
	// 同步表
	if err := db.DB(context.Background()).AutoMigrate(&entity.Folder{}, &entity.File{}, &entity.FileShare{}); err != nil {
		hlog.Fatalf("sync sys storage tables to db error: %v", err)
	}
	return &storageRepository{
		db:     db,
		repo:   repo,
		mapper: &mapper.StorageMapper{},
	}
}

// GetFile 获取文件
func (r *storageRepository) GetFile(ctx context.Context, id string) (*model.File, error) {
	file, err := r.repo.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToFileDomain(file), nil
}

// SaveFile 保存文件
func (r *storageRepository) SaveFile(ctx context.Context, file *model.File, reader io.Reader) (*model.File, error) {

	// 2. 转换为实体
	fileEntity := r.mapper.ToFileEntity(file) // 1. 生成ID
	if file.ID == "" {
		fileEntity.ID = r.db.GenStringId()
		// 3. 保存到数据库
		if err := r.repo.CreateFile(ctx, fileEntity); err != nil {
			return nil, err
		}
	} else {
		if err := r.repo.UpdateFile(ctx, fileEntity); err != nil {
			return nil, err
		}
	}
	return r.mapper.ToFileDomain(fileEntity), nil
}

// DeleteFile 删除文件
func (r *storageRepository) DeleteFile(ctx context.Context, id string) error {
	return r.repo.DeleteFile(ctx, id)
}

// ListFiles 查询文件列表
func (r *storageRepository) ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*model.File, int64, error) {
	files, total, err := r.repo.ListFiles(ctx, folderID, qb)
	if err != nil {
		return nil, 0, err
	}
	return r.mapper.ToFileDomainList(files), total, nil
}

// GetFolder 获取文件夹
func (r *storageRepository) GetFolder(ctx context.Context, id string) (*model.Folder, error) {
	folder, err := r.repo.GetFolder(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToFolderDomain(folder), nil
}

// SaveFolder 保存文件夹
func (r *storageRepository) SaveFolder(ctx context.Context, folder *model.Folder) error {
	// 1. 生成ID
	if folder.ID == "" {
		folder.ID = r.db.GenStringId()
	}

	// 2. 转换为实体
	folderEntity := r.mapper.ToFolderEntity(folder)

	// 3. 保存到数据库
	if err := r.repo.CreateFolder(ctx, folderEntity); err != nil {
		return err
	}

	return nil
}

// DeleteFolder 删除文件夹
func (r *storageRepository) DeleteFolder(ctx context.Context, id string) error {
	return r.repo.DeleteFolder(ctx, id)
}

// ListFolders 查询文件夹列表
func (r *storageRepository) ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*model.Folder, int64, error) {
	folders, total, err := r.repo.ListFolders(ctx, parentID, qb)
	if err != nil {
		return nil, 0, err
	}
	return r.mapper.ToFolderDomainList(folders), total, nil
}

// CreateFileShare 创建文件分享记录
func (r *storageRepository) CreateFileShare(ctx context.Context, share *model.FileShare) error {
	// 1. 生成ID
	if share.ID == "" {
		share.ID = r.db.GenStringId()
	}

	// 2. 转换为实体
	shareEntity := r.mapper.ToFileShareEntity(share)

	// 3. 保存到数据库
	return r.repo.CreateFileShare(ctx, shareEntity)
}
