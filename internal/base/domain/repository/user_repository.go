package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// 新增动态查询方法
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error)
	Count(ctx context.Context, qb *query.QueryBuilder) (int64, error)

	// BelongsToDepartment 检查用户是否属于指定部门
	BelongsToDepartment(ctx context.Context, userID string, deptID string) bool

	// FindByDepartment 查询部门下的用户(排除管理员)
	FindByDepartment(ctx context.Context, deptID string, excludeAdminID string) ([]*model.User, error)

	// FindUnassignedUsers 查询未分配部门的用户
	FindUnassignedUsers(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error)
}
