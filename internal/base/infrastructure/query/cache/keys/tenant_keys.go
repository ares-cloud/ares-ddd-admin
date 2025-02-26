package keys

import "fmt"

const (
	tenantPrefix = "tenant:"
)

// TenantKey 生成租户缓存key
func TenantKey(tenantID string) string {
	return fmt.Sprintf("%s%s", tenantPrefix, tenantID)
}

// TenantPermKey 生成租户权限缓存key
func TenantPermKey(tenantID string) string {
	return fmt.Sprintf("tenant:perm:%s", tenantID)
}

// TenantListKey 生成租户列表缓存key
func TenantListKey() string {
	return fmt.Sprintf("%slist", tenantPrefix)
}

// TenantPermissionsKey 租户权限缓存key
func TenantPermissionsKey(tenantID string) string {
	return fmt.Sprintf("%s%s:permissions", tenantPrefix, tenantID)
}

// TenantKeys 生成租户相关的所有缓存key
func TenantKeys(tenantID string) []string {
	return []string{
		TenantKey(tenantID),
		TenantPermissionsKey(tenantID),
		TenantStatusKey(tenantID),
	}
}

// TenantStatusKey 租户状态缓存key
func TenantStatusKey(tenantID string) string {
	return fmt.Sprintf("%s%s:status", tenantPrefix, tenantID)
}
