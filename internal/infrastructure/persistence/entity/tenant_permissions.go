package entity

// TenantPermissions 租户权限关联
type TenantPermissions struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`
	TenantID     string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`
	PermissionID int64  `json:"permission_id" gorm:"index;comment:权限ID"`
}

// TableName 定义表名
func (t TenantPermissions) TableName() string {
	return "sys_tenant_permissions"
}

// GetPrimaryKey ， 定义表主键 base repo 会使用，非 gorm 原生接口
// 参数：
// 返回值：
//
//	string ：表主键
func (TenantPermissions) GetPrimaryKey() string {
	return "id"
}
