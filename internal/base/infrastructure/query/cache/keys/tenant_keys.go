package keys

import "fmt"

// TenantKey 生成租户缓存key
func TenantKey(tenantID string) string {
	return fmt.Sprintf("tenant:%s", tenantID)
}

// TenantPermKey 生成租户权限缓存key
func TenantPermKey(tenantID string) string {
	return fmt.Sprintf("tenant:perm:%s", tenantID)
}

// TenantListKey 生成租户列表缓存key
func TenantListKey() string {
	return "tenant:list"
}

// TenantPermissionsKey 租户权限缓存key
func TenantPermissionsKey(tenantID string) string {
	return fmt.Sprintf("tenant:permissions:%s", tenantID)
}

// TenantKeys 生成租户相关的所有缓存key
func TenantKeys(tenantID string) []string {
	return []string{
		TenantKey(tenantID),
		TenantPermissionsKey(tenantID),
	}
}

// 租户相关的缓存键
func TenantStatusKey(tenantID string) string {
	return fmt.Sprintf("tenant:status:%s", tenantID)
}
