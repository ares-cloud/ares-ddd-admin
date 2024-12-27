package entity

// DataPermission 数据权限实体
type DataPermission struct {
	ID       string `gorm:"column:id;primary_key"`
	RoleID   int64  `gorm:"column:role_id"`
	Scope    int8   `gorm:"column:scope"`
	DeptIDs  string `gorm:"column:dept_ids"` // JSON数组字符串
	TenantID string `gorm:"column:tenant_id"`
}

// TableName 表名
func (DataPermission) TableName() string {
	return "sys_data_permission"
}
