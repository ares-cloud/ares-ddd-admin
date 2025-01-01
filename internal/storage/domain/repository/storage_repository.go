package repository

import (
	"context"
	"io"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type IStorageRepository interface {
	// 文件相关
	GetFile(ctx context.Context, id string) (*model.File, error)
	SaveFile(ctx context.Context, file *model.File, reader io.Reader) (*model.File, error)
	DeleteFile(ctx context.Context, id string) error
	ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*model.File, int64, error)

	// 文件夹相关
	GetFolder(ctx context.Context, id string) (*model.Folder, error)
	SaveFolder(ctx context.Context, folder *model.Folder) error
	DeleteFolder(ctx context.Context, id string) error
	ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*model.Folder, int64, error)

	// 文件分享相关
	CreateFileShare(ctx context.Context, share *model.FileShare) error
}
