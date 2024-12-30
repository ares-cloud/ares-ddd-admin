package queries

import "github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	ID string
}

// ListUsersQuery 用户列表查询
type ListUsersQuery struct {
	query.Page
	Username string
	Name     string
	Phone    string
	Email    string
	Status   int
}

// GetUserPermissionsQuery 获取用户权限查询
type GetUserPermissionsQuery struct {
	UserID string
}

// GetUserMenusQuery 获取用户菜单查询
type GetUserMenusQuery struct {
	UserID string
}

// GetUserInfoQuery 查询用户信息
type GetUserInfoQuery struct {
	UserID string
}

// ListDepartmentUsersQuery 部门用户列表查询
type ListDepartmentUsersQuery struct {
	query.Page
	DeptID         string
	ExcludeAdminID string
	Username       string
	Name           string
}

// ListUnassignedUsersQuery 未分配部门用户列表查询
type ListUnassignedUsersQuery struct {
	query.Page
	Username string
	Name     string
}
