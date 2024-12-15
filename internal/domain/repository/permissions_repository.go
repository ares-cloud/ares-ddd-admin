package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type IPermissionsRepository interface {
	Create(ctx context.Context, permissions *model.Permissions) error
	Update(ctx context.Context, permissions *model.Permissions) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Permissions, error)
	FindByCode(ctx context.Context, code string) (*model.Permissions, error)
	FindByRoleID(ctx context.Context, roleID int64) ([]*model.Permissions, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindByType(ctx context.Context, permType int8) ([]*model.Permissions, error)

	// 新增动态查询方法
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.Permissions, error)
	Count(ctx context.Context, qb *query.QueryBuilder) (int64, error)

	// 构建权限树相关方法
	FindAllTree(ctx context.Context) ([]*model.Permissions, error)                                         // 构建所有权限树
	FindTreeByType(ctx context.Context, permType int8) ([]*model.Permissions, error)                       // 根据类型构建权限树
	FindTreeByQuery(ctx context.Context, qb *query.QueryBuilder) ([]*model.Permissions, error)             // 根据查询条件构建权限树
	FindTreeByUserAndType(ctx context.Context, userID string, permType int8) ([]*model.Permissions, error) // 根据用户和类型构建权限树
}
