package keys

import (
	"fmt"
)

const (
	// 角色缓存key前缀
	rolePrefix     = "role:"
	roleCodePrefix = "role:code:"
	rolePermPrefix = "role:perm:"
)

// RoleKey 生成角色缓存key
func RoleKey(tenantID string, id int64) string {
	return fmt.Sprintf("%s%s:%d", rolePrefix, tenantID, id)
}

// RoleCodeKey 生成角色编码缓存key
func RoleCodeKey(tenantID string, code string) string {
	return fmt.Sprintf("%s%s:code:%s", rolePrefix, tenantID, code)
}

// RolePermKey 生成角色权限缓存key
func RolePermKey(tenantID string, roleID int64) string {
	return fmt.Sprintf("%s%s:perm:%d", rolePrefix, tenantID, roleID)
}

// RolePermissionsKey 角色权限缓存key
func RolePermissionsKey(tenantID string, roleID int64) string {
	return fmt.Sprintf("%s%s:permissions:%d", rolePrefix, tenantID, roleID)
}

// RoleListKey 角色列表缓存key
func RoleListKey(tenantID string) string {
	return fmt.Sprintf("%s%s:list", rolePrefix, tenantID)
}

// RoleKeys 生成角色相关的所有缓存key
func RoleKeys(tenantID string, roleID int64) []string {
	return []string{
		RoleKey(tenantID, roleID),
		RolePermKey(tenantID, roleID),
		RolePermissionsKey(tenantID, roleID),
		RoleListKey(tenantID),
	}
}

// RoleUsersKey 角色用户缓存key
func RoleUsersKey(tenantID string, roleID int64) string {
	return fmt.Sprintf("%s%s:users:%d", rolePrefix, tenantID, roleID)
}
