package queries

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type GetRoleQuery struct {
	Id int64 `json:"id" query:"id"` // 角色id
}

type ListRolesQuery struct {
	query.Page
	Code   string `json:"code" query:"code"`    // 编码
	Name   string `json:"name" query:"name"`    // 名称
	Status int8   `json:"status" query:"email"` // 角色状态（禁用、启用）
}

type GetUserRolesQuery struct {
	UserID string `json:"user_id" query:"user_id"` // 用户ID
}