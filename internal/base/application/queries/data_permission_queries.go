package queries

// GetDataPermissionQuery 获取数据权限查询
type GetDataPermissionQuery struct {
	RoleID int64 `json:"roleId" validate:"required" label:"角色ID"`
}
