package queries

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// ListDepartmentsQuery 部门列表查询
type ListDepartmentsQuery struct {
	query.Page
	Name   string `json:"name" validate:"omitempty,max=50" label:"部门名称"`    // 部门名称
	Code   string `json:"code" validate:"omitempty,max=50" label:"部门编码"`    // 部门编码
	Status *int8  `json:"status" validate:"omitempty,oneof=0 1" label:"状态"` // 部门状态
}

func (q *ListDepartmentsQuery) Validate() herrors.Herr {
	return validator.Validate(q)
}

// GetDepartmentQuery 获取部门查询
type GetDepartmentQuery struct {
	ID string `json:"id" validate:"required" label:"部门ID"` // 部门ID
}

func (q *GetDepartmentQuery) Validate() herrors.Herr {
	return validator.Validate(q)
}

// GetDepartmentTreeQuery 获取部门树查询
type GetDepartmentTreeQuery struct {
	ParentID string `json:"parentId" validate:"omitempty" label:"父部门ID"` // 父部门ID,为空则查���全部
}

func (q *GetDepartmentTreeQuery) Validate() herrors.Herr {
	return validator.Validate(q)
}

// GetUserDepartmentsQuery 获取用户部门查询
type GetUserDepartmentsQuery struct {
	UserID string `json:"userId" validate:"required" label:"用户ID"` // 用户ID
}

func (q *GetUserDepartmentsQuery) Validate() herrors.Herr {
	return validator.Validate(q)
}
