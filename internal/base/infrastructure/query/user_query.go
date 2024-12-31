package query

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IUserQueryService 用户查询服务接口
type IUserQueryService interface {
	// 基础查询
	GetUser(ctx context.Context, id string) (*dto.UserDto, error)
	FindUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error)
	CountUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)

	// 权限相关查询
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetUserRoles(ctx context.Context, userID string) ([]*dto.RoleDto, error)
	GetUserMenus(ctx context.Context, userID string) ([]*dto.PermissionsDto, error)
	GetUserTreeMenus(ctx context.Context, userID string) ([]*dto.PermissionsTreeDto, error)
	GetUserRolesCode(ctx context.Context, userID string) ([]string, error)

	// 部门相关查询
	FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error)
	CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error)
}
