package query

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// IPermissionsQuery 权限查询接口
type IPermissionsQuery interface {
	// FindByID 根据ID查询权限
	FindByID(ctx context.Context, id int64) (*dto.PermissionsDto, herrors.Herr)
	// Find 查询权限列表
	Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.PermissionsDto, int64, herrors.Herr)
	// FindTreeByType 查询权限树
	FindTreeByType(ctx context.Context, permType int8) ([]*dto.PermissionsDto, herrors.Herr)
	// FindAllEnabled 查询所有启用的权限
	FindAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, herrors.Herr)
	// GetSimplePermissionsTree 获取简化的权限树
	GetSimplePermissionsTree(ctx context.Context) (*dto.PermissionsTreeResult, herrors.Herr)
}
