package queries

import "github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

type GetPermissionsQuery struct {
	Id int64 `json:"id" query:"id"` // 权限ID
}

type ListPermissionsQuery struct {
	query.Page
	Code   string `json:"code" query:"code"`    // 编码
	Name   string `json:"name" query:"name"`    // 名称
	Status int8   `json:"status" query:"email"` // 角色状态（禁用、启用）
}

type GetPermissionsTreeQuery struct {
	Type string `json:"type" query:"type"` // 权限类型
}
