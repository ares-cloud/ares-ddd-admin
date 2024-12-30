package repository

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// IDepartmentRepository 部门仓储接口
type IDepartmentRepository interface {
	// Create 创建部门
	Create(ctx context.Context, dept *model.Department) error

	// Update 更新部门
	Update(ctx context.Context, dept *model.Department) error

	// Delete 删除部门
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取部门
	GetByID(ctx context.Context, id string) (*model.Department, error)

	// GetByCode 根据编码获取部门
	GetByCode(ctx context.Context, code string) (*model.Department, error)

	// GetByParentID 获取子部门
	GetByParentID(ctx context.Context, parentID string) ([]*model.Department, error)

	// List 查询部门列表
	List(ctx context.Context, query *ListDepartmentQuery) ([]*model.Department, error)

	// GetUserDepartments 获取用户部门
	GetUserDepartments(ctx context.Context, userID string) ([]*model.Department, error)

	// GetAllUserIDs 获取所有用户ID
	GetAllUserIDs(ctx context.Context) ([]string, error)

	// AssignUsers 分配用户到部门
	AssignUsers(ctx context.Context, deptID string, userIDs []string) error

	// RemoveUsers 从部门移除用户
	RemoveUsers(ctx context.Context, deptID string, userIDs []string) error

	// Find 查询部门列表
	Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*model.Department, error)

	// Count 查询总数
	Count(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
}

// ListDepartmentQuery 部门列表查询参数
type ListDepartmentQuery struct {
	Name     string `json:"name"`     // 部门名称
	Code     string `json:"code"`     // 部门编码
	Status   *int8  `json:"status"`   // 部门状态
	ParentID string `json:"parentId"` // 父部门ID
}
