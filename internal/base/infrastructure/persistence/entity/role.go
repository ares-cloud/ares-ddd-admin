package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// RoleType 角色类型
type RoleType int8

const (
	RoleTypeResource RoleType = 1 // 资源类型角色(默认)
	RoleTypeData     RoleType = 2 // 数据权限角色
)

// Role 基于角色的访问控制 (RBAC) 的角色管理
type Role struct {
	database.BaseModel
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`        // 唯一ID
	Code        string `json:"code" gorm:"size:32;index;comment:角色代码（唯一）"`             // 角色代码（唯一）
	Name        string `json:"name" gorm:"size:128;index;comment:角色显示名称"`              // 角色显示名称
	Type        int8   `json:"type" gorm:"type:int8;default:1;comment:角色类型"`           // 角色类型(1:资源角色 2:数据权限角色)
	Localize    string `json:"localize" gorm:"size:128;comment:国际化key;"`               // 国际化key
	Description string `json:"description" gorm:"size:1024;comment:角色的详细信息"`           // 角色的详细信息
	Sequence    int    `json:"sequence" gorm:"index;comment:排序顺序"`                     // 排序顺序
	Status      int8   `json:"status" gorm:"column:status;default:1;comment:状态，启用,禁用"` // 用户状态（激活、冻结）
	TenantID    string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`            // 租户ID
}

// TableName 定义数据库中角色表的名称
func (a Role) TableName() string {
	return "sys_role"
}

// GetPrimaryKey 定义表主键
func (a Role) GetPrimaryKey() string {
	return "id"
}

// IsResourceRole 是否为资源类型角色
func (a Role) IsResourceRole() bool {
	return a.Type == int8(RoleTypeResource)
}

// IsDataRole 是否为数据权限角色
func (a Role) IsDataRole() bool {
	return a.Type == int8(RoleTypeData)
}
