package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// Tenant 租户实体
type Tenant struct {
	database.BaseModel
	ID          string `json:"id" gorm:"primaryKey;size:32;comment:租户ID"`
	Code        string `json:"code" gorm:"size:32;uniqueIndex;comment:租户编码"`
	Name        string `json:"name" gorm:"size:128;comment:租户名称"`
	AdminUserID string `json:"admin_user_id" gorm:"size:32;comment:管理员用户ID"`
	Status      int8   `json:"status" gorm:"default:1;comment:状态(1:启用 2:禁用)"`
	IsDefault   int8   `json:"is_default" gorm:"default:2;comment:是否默认租户(1:是 2:否)"`
	ExpireTime  int64  `json:"expire_time" gorm:"comment:过期时间"`
	Description string `json:"description" gorm:"size:512;comment:描述"`
}

// TableName 定义表名
func (t Tenant) TableName() string {
	return "sys_tenant"
}

// GetPrimaryKey 获取主键字段名
func (t Tenant) GetPrimaryKey() string {
	return "id"
}
