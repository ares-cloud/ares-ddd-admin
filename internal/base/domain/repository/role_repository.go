package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type IRoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Role, error)
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	FindByUserID(ctx context.Context, userID string) ([]*model.Role, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Role, error)

	// 新增动态查询方法
	Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*model.Role, error)
	Count(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)

	// FindAllEnabled 获取所有启用状态的角色
	FindAllEnabled(ctx context.Context) ([]*model.Role, error)

	// FindByType 根据角色类型查询角色列表
	FindByType(ctx context.Context, roleType int8) ([]*model.Role, error)

	// GetPermissionsByRoleID 获取角色的权限ID列表
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]int64, error)

	// FindByPermissionID 根据权限ID查找角色
	FindByPermissionID(ctx context.Context, permissionID int64) ([]*model.Role, error)

	// UpdatePermissions 更新角色权限
	UpdatePermissions(ctx context.Context, roleID int64, permIDs []int64) error

	// GetRolePermissions 获取角色权限
	GetRolePermissions(ctx context.Context, roleID int64) ([]*model.Permissions, error)

	// GetRoleDataPermission 获取角色数据权限
	GetRoleDataPermission(ctx context.Context, roleID int64) (*model.DataPermission, error)
}
