package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type ITenantRepository interface {
	// 基础操作
	Create(ctx context.Context, tenant *model.Tenant) error
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id string) error

	// 业务规则验证
	FindByID(ctx context.Context, id string) (*model.Tenant, error)
	FindByCode(ctx context.Context, code string) (*model.Tenant, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// 权限相关
	AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error
	GetPermissions(ctx context.Context, tenantID string) ([]*model.Permissions, error)
	HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error)

	// 锁定相关
	Lock(ctx context.Context, tenantID string, reason string) error
	Unlock(ctx context.Context, tenantID string) error
}
