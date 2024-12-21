package entity

// SysUserRole 基于角色的访问控制 (RBAC) 的用户角色
type SysUserRole struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"` // 唯一ID
	UserID   string `json:"user_id" gorm:"index;comment:来源于 User.ID"`        // 来源于 User.ID
	RoleID   int64  `json:"role_id" gorm:"index;comment:来源于 Role.ID"`        // 来源于 Role.ID
	TenantID string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`     // 租户ID
}

// TableName 定义数据库中用户角色表的名称
func (a *SysUserRole) TableName() string {
	return "sys_user_role"
}
