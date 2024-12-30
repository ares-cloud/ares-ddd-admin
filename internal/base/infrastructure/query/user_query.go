package query

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

// UserQueryService 用户查询服务接口
type UserQueryService interface {
	// 基础查询
	GetUser(ctx context.Context, id string) (*model.User, error)
	FindUsers(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error)
	CountUsers(ctx context.Context, qb *query.QueryBuilder) (int64, error)

	// 权限相关查询
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error)
	GetUserRolesCode(ctx context.Context, userID string) ([]string, error)
	GetUserMenus(ctx context.Context, userID string) ([]*model.Permissions, error)
	GetUserTreeMenus(ctx context.Context, userID string) ([]*model.Permissions, error)

	// 部门相关查询
	FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *query.QueryBuilder) ([]*model.User, error)
	CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *query.QueryBuilder) (int64, error)
	FindUnassignedUsers(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error)
	CountUnassignedUsers(ctx context.Context, qb *query.QueryBuilder) (int64, error)
}
