package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// IUserRepository 用户仓储接口 - 仅包含命令操作
type IUserRepository interface {
	// 基础操作
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error

	// 用于业务规则验证
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// 角色分配
	AssignRoles(ctx context.Context, userID string, roleIDs []int64) error

	// 部门相关
	BelongsToDepartment(ctx context.Context, userID string, deptID string) (bool, error)
}
