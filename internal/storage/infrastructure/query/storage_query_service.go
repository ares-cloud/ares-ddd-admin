package query

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IStorageQueryService 存储查询服务接口
type IStorageQueryService interface {
	// ListFolders 查询文件夹列表
	ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*dto.FolderDto, int64, error)
	// ListFiles 查询文件列表
	ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*dto.FileDto, int64, error)
	// GetShareFile 获取分享文件
	GetShareFile(ctx context.Context, shareCode string, password string) (*dto.FileDto, error)
}
