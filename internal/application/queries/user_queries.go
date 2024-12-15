package queries

import "github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

type GetUserQuery struct {
	Id string `json:"id" query:"id"` // 用户ID
}

type ListUsersQuery struct {
	query.Page
	Username string `json:"username" query:"username"` // 登录用户名
	Name     string `json:"name" query:"name"`         // 用户姓名
	Phone    string `json:"phone" query:"phone"`       // 用户手机号
	Email    string `json:"email" query:"email"`       // 用户邮箱
	Status   int8   `json:"status" query:"status"`     // 角色状态（禁用、启用）
}
