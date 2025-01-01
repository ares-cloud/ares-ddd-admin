package keys

import "fmt"

const (
	userPrefix = "sys_user_cache:"
)

// UserKey 用户缓存key
func UserKey(userID string) string {
	return fmt.Sprintf("%s:%s", userPrefix, userID)
}

// UserPermissionsKey 用户权限缓存key
func UserPermissionsKey(userID string) string {
	return fmt.Sprintf("user:permissions:%s", userID)
}

// UserRolesKey 用户角色缓存key
func UserRolesKey(userID string) string {
	return fmt.Sprintf("%s:roles:%s", userPrefix, userID)
}

// UserMenusKey 用户菜单缓存key
func UserMenusKey(userID string) string {
	return fmt.Sprintf("user:menus:%s", userID)
}

// UserRoleCodesKey 用户角色编码缓存key
func UserRoleCodesKey(userID string) string {
	return fmt.Sprintf("%s:role:codes:%s", userPrefix, userID)
}

// UserDepartmentKey 用户部门缓存key
func UserDepartmentKey(userID string) string {
	return fmt.Sprintf("user:department:%s", userID)
}

// UserKeys 生成用户相关的所有缓存key
func UserKeys(userID string) []string {
	return []string{
		UserKey(userID),
		UserPermissionsKey(userID),
		UserRolesKey(userID),
		UserMenusKey(userID),
		UserRoleCodesKey(userID),
	}
}

// UserListKey 用户列表缓存key
func UserListKey() string {
	return fmt.Sprintf("%s:list", userPrefix)
}
