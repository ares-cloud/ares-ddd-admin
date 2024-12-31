package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type IPermissionsRepository interface {
	// 基础CRUD方法
	Create(ctx context.Context, permissions *model.Permissions) error
	Update(ctx context.Context, permissions *model.Permissions) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Permissions, error)

	// 业务查询方法
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
