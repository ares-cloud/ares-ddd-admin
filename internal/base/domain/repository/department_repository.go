package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// IDepartmentRepository 部门仓储接口
type IDepartmentRepository interface {
	// 基础操作
	FindByID(ctx context.Context, id string) (*model.Department, error)
	Create(ctx context.Context, dept *model.Department) error
	Update(ctx context.Context, dept *model.Department) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindByCode(ctx context.Context, code string) (*model.Department, error)

	// 树形结构操作
	GetByParentID(ctx context.Context, parentID string) ([]*model.Department, error)
	GetTreeByParentID(ctx context.Context, parentID string) ([]*model.Department, error)

	// 用户部门操作
	AssignUsers(ctx context.Context, deptID string, userIDs []string) error
	RemoveUsers(ctx context.Context, deptID string, userIDs []string) error
	TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error
}
