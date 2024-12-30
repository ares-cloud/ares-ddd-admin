package events

// CacheInvalidateEvent 缓存失效事件
type CacheInvalidateEvent struct {
	TenantID string   // 租户ID
	UserIDs  []string // 受影响的用户ID
	RoleIDs  []int64  // 受影响的角色ID
}
