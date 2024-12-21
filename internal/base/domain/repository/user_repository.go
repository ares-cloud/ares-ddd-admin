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
}
