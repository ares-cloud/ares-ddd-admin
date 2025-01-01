package keys

import "fmt"

const (
	dataPermissionPrefix = "data_permission:"
)

// DataPermissionKey 数据权限缓存key
func DataPermissionKey(tenantID string, roleID int64) string {
	return fmt.Sprintf("%s%s:role:%d", dataPermissionPrefix, tenantID, roleID)
}
