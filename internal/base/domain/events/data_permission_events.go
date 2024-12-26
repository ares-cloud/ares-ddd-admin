package events

// DataPermissionAssignedEvent 数据权限分配事件
type DataPermissionAssignedEvent struct {
	ID     string `json:"id"`
	RoleID string `json:"role_id"`
	Type   int8   `json:"type"`
}

// UserDeptAssignedEvent 用户部门分配事件
type UserDeptAssignedEvent struct {
	UserID string `json:"user_id"`
	DeptID string `json:"dept_id"`
}
