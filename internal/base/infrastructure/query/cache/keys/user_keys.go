package keys

import "fmt"

const (
	// 缓存键前缀
	userKeyPrefix = "user:"
	roleKeyPrefix = "role:"
	permKeyPrefix = "perm:"
	menuKeyPrefix = "menu:"
)

// 用户相关缓存键生成函数
func UserKey(userID string) string {
	return fmt.Sprintf("%s%s", userKeyPrefix, userID)
}

func UserPermKey(userID string) string {
	return fmt.Sprintf("%s%s:permissions", userKeyPrefix, userID)
}

func UserRoleKey(userID string) string {
	return fmt.Sprintf("%s%s:roles", userKeyPrefix, userID)
}

func UserMenuKey(userID string) string {
	return fmt.Sprintf("%s%s:menus", userKeyPrefix, userID)
}
func UserMenuTreeKey(userID string) string {
	return fmt.Sprintf("%s%s:menus_tree", userKeyPrefix, userID)
}

// 列表缓存键
func UserListKey() string {
	return fmt.Sprintf("%slist", userKeyPrefix)
}

// 角色相关缓存键
func RoleListKey() string {
	return fmt.Sprintf("%slist", roleKeyPrefix)
}
