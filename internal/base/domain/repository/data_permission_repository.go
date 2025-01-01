package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type IDataPermissionRepository interface {
	// GetByRoleID 获取角色的数据权限
	GetByRoleID(ctx context.Context, roleID int64) (*model.DataPermission, error)

	// GetByRoleIDs 批量获取角色的数据权限
	GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.DataPermission, error)

	// Save 保存数据权限(创建或更新)
	Save(ctx context.Context, perm *model.DataPermission) error

	// DeleteByRoleID 根据角色ID删除数据权限
	DeleteByRoleID(ctx context.Context, roleID int64) error
}
