package keys

import "fmt"

const (
	permissionKeyPrefix = "permission:"
)

// PermissionKey 权限详情缓存key
func PermissionKey(tenantID string, id int64) string {
	return fmt.Sprintf("%s%s:detail:%d", permissionKeyPrefix, tenantID, id)
}

// PermissionTreeKey 权限树缓存key
func PermissionTreeKey(tenantID string, permType interface{}) string {
	return fmt.Sprintf("%s%s:tree:%v", permissionKeyPrefix, tenantID, permType)
}

// PermissionListKey 权限列表缓存key
func PermissionListKey(tenantID string) string {
	return fmt.Sprintf("%s%s:list", permissionKeyPrefix, tenantID)
}

// PermissionEnabledKey 启用权限列表缓存key
func PermissionEnabledKey(tenantID string) string {
	return fmt.Sprintf("%s%s:enabled", permissionKeyPrefix, tenantID)
}

// PermissionSimpleTreeKey 简化权限树缓存key
func PermissionSimpleTreeKey(tenantID string) string {
	return fmt.Sprintf("%s%s:simple:tree", permissionKeyPrefix, tenantID)
}
