package repository

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type IStorageRepository interface {
	// File operations
	CreateFile(ctx context.Context, file *model.File) error
	UpdateFile(ctx context.Context, file *model.File) error
	DeleteFile(ctx context.Context, id string) error
	GetFile(ctx context.Context, id string) (*model.File, error)
	ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*model.File, int64, error)

	// Folder operations
	CreateFolder(ctx context.Context, folder *model.Folder) error
	UpdateFolder(ctx context.Context, folder *model.Folder) error
	DeleteFolder(ctx context.Context, id string) error
	GetFolder(ctx context.Context, id string) (*model.Folder, error)
	ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*model.Folder, int64, error)

	// Share operations
	GetFileShare(ctx context.Context, shareCode string) (*model.FileShare, error)
	CreateFileShare(ctx context.Context, share *model.FileShare) error

	// GetFolderTree 获取文件夹树形结构
	GetFolderTree(ctx context.Context, tenantID string) ([]*model.FolderTree, error)

	// GetExpiredRecycleFiles 获取过期的回收站文件
	GetExpiredRecycleFiles(ctx context.Context, expireTime time.Time) ([]*model.File, error)
}
