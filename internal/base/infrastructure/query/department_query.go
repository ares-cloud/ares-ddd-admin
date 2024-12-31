package query

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IDepartmentQueryService 部门查询服务接口
type IDepartmentQueryService interface {
	// 基础查询
	GetDepartment(ctx context.Context, id string) (*dto.DepartmentDto, error)
	FindDepartments(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.DepartmentDto, error)
	CountDepartments(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)

	// 树形结构查询
	GetDepartmentTree(ctx context.Context, parentID string) ([]*dto.DepartmentTreeDto, error)

	// 用户部门查询
	GetUserDepartments(ctx context.Context, userID string) ([]*dto.DepartmentDto, error)
	GetDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error)
	CountDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error)

	// 未分配用户查询
	GetUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error)
	CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
}
