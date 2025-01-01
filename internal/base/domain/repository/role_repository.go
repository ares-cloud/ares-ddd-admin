package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// IRoleRepository 角色仓储接口
type IRoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Role, error)
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	// ExistsById 检查数据权限是否存在
	ExistsById(ctx context.Context, id int64) (bool, error)
	// 权限相关
	AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*model.Permissions, error)
	ExistsPermission(ctx context.Context, permissionID int64) (bool, error)

	// 业务相关
	IsRoleInUse(ctx context.Context, roleID int64) (bool, error)
}
