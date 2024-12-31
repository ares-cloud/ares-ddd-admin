package keys

import (
	"fmt"
	"strconv"
)

const (
	// 角色缓存key前缀
	rolePrefix     = "role:"
	roleCodePrefix = "role:code:"
	rolePermPrefix = "role:perm:"
)

// RoleKey 生成角色缓存key
func RoleKey(id int64) string {
	return rolePrefix + strconv.FormatInt(id, 10)
}

// RoleCodeKey 生成角色编码缓存key
func RoleCodeKey(code string) string {
	return fmt.Sprintf("role:code:%s", code)
}

// RolePermKey 生成角色权限缓存key
func RolePermKey(roleID int64) string {
	return fmt.Sprintf("%s%d", rolePermPrefix, roleID)
}

// RolePermissionsKey 角色权限缓存key
func RolePermissionsKey(roleID int64) string {
	return fmt.Sprintf("role:permissions:%d", roleID)
}

// RoleListKey 角色列表缓存key
func RoleListKey() string {
	return "role:list"
}

// RoleKeys 生成角色相关的所有缓存key
func RoleKeys(roleID int64) []string {
	return []string{
		RoleKey(roleID),
		RolePermKey(roleID),
		RolePermissionsKey(roleID),
		RoleListKey(),
	}
}

// RoleUsersKey 角色用户缓存key
func RoleUsersKey(roleID int64) string {
	return fmt.Sprintf("role:users:%d", roleID)
}
