package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ITenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Tenant, error)
	FindByCode(ctx context.Context, code string) (*model.Tenant, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// 动态查询方法
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.Tenant, error)
	Count(ctx context.Context, qb *query.QueryBuilder) (int64, error)

	// 权限相关方法
	AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error
	GetPermissions(ctx context.Context, tenantID string) ([]*model.Permissions, error)
	HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error)
}
