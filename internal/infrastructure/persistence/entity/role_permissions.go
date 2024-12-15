package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// RolePermissions 基于角色的访问控制 (RBAC) 的角色菜单权限
type RolePermissions struct {
	database.BaseModel
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"` // 唯一ID
	RoleID       int64  `json:"role_id" gorm:"index;comment:角色ID"`               // 来源于 角色ID
	PermissionID int64  `json:"permission_id" gorm:"index;comment:权限ID"`         // 来源于 权限ID
	TenantID     string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`     // 租户ID
}

// TableName 定义数据库中角色菜单表的名称
func (a *RolePermissions) TableName() string {
	return "sys_role_permissions"
}
