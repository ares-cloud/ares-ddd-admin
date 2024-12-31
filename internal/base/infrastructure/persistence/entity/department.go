package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// Department 部门数据库实体
type Department struct {
	database.BaseModel
	ID          string `json:"id" gorm:"primaryKey;size:32;comment:部门ID"`     // 部门ID
	TenantID    string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`   // 租户ID
	ParentID    string `json:"parent_id" gorm:"size:32;index;comment:父部门ID"`  // 父部门ID
	Code        string `json:"code" gorm:"size:50;uniqueIndex;comment:部门编码"`  // 部门编码
	Name        string `json:"name" gorm:"size:100;comment:部门名称"`             // 部门名称
	Sequence    int32  `json:"sequence" gorm:"default:0;comment:显示顺序"`        // 显示顺序
	AdminID     string `json:"admin_id" gorm:"size:32;index;comment:管理员ID"`   // 管理员ID
	Leader      string `json:"leader" gorm:"size:50;comment:负责人"`             // 负责人
	Phone       string `json:"phone" gorm:"size:20;comment:联系电话"`             // 联系电话
	Email       string `json:"email" gorm:"size:100;comment:邮箱"`              // 邮箱
	Status      int8   `json:"status" gorm:"default:1;comment:部门状态(0停用 1启用)"` // 部门状态
	Description string `json:"description" gorm:"size:200;comment:描述"`        // 描述
}

// TableName 表名
func (Department) TableName() string {
	return "sys_department"
}

// GetPrimaryKey ， 定义表主键 base repo 会使用，非 gorm 原生接口
// 参数：
// 返回值：
//
//	string ：表主键
func (Department) GetPrimaryKey() string {
	return "id"
}
