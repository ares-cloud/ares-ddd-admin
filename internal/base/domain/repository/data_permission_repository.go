package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type IDataPermissionRepository interface {
	// Create 创建数据权限
	Create(ctx context.Context, perm *model.DataPermission) error

	// Update 更新数据权限
	Update(ctx context.Context, perm *model.DataPermission) error

	// Delete 删除数据权限
	Delete(ctx context.Context, id string) error

	// GetByRoleID 获取角色的数据权限
	GetByRoleID(ctx context.Context, roleID int64) (*model.DataPermission, error)

	// GetByRoleIDs 批量获取角色的数据权限
	GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.DataPermission, error)
}
