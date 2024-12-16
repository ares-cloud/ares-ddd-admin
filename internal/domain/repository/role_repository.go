package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
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
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.Role, error)
	Count(ctx context.Context, qb *query.QueryBuilder) (int64, error)

	// FindAllEnabled 获取所有启用状态的角色
	FindAllEnabled(ctx context.Context) ([]*model.Role, error)
}
