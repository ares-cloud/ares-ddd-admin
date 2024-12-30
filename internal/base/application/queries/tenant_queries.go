package queries

import "github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"

type GetTenantQuery struct {
	Id string `json:"id" query:"id"`
}

type ListTenantsQuery struct {
	db_query.Page
	Code   string `json:"code" query:"code"`
	Name   string `json:"name" query:"name"`
	Status int8   `json:"status" query:"status"`
}

// GetTenantPermissionsQuery 获取租户权限查询
type GetTenantPermissionsQuery struct {
	TenantID string `json:"tenant_id" query:"tenant_id"`
}
