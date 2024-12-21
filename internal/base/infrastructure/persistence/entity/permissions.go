package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// Permissions 菜单管理用于基于角色的访问控制 (RBAC)
type Permissions struct {
	database.BaseModel
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`                    // 唯一ID
	Code        string `json:"code" gorm:"size:32;index;comment:菜单代码（每个层级唯一）;"`                    // 菜单代码（每个层级唯一）
	Name        string `json:"name" gorm:"size:128;index;comment:菜单显示名称;"`                         // 菜单显示名称
	Localize    string `json:"localize" gorm:"size:128;comment:国际化key;"`                           // 国际化key
	Icon        string `json:"Icon" gorm:"size:50; comment:图标"`                                    // 图标
	Description string `json:"description" gorm:"size:1024;comment:菜单的详细信息;"`                      // 菜单的详细信息
	Sequence    int    `json:"sequence" gorm:"index;comment:排序顺序（降序排列）;"`                          // 排序顺序（降序排列）
	Type        int8   `json:"type" gorm:"column:type;default:1;comment:菜单类型(1、页面、2、按钮、3、api接口);"` // 菜单类型（页面、按钮）
	Path        string `json:"path" gorm:"size:255;comment:菜单的访问路径;"`                              // 菜单的访问路径
	Properties  string `json:"properties" gorm:"type:text;comment:菜单的属性(JSON格式);"`                 // 菜单的属性(JSON格式）
	Status      int8   `json:"status" gorm:"column:status;default:1;comment:状态,1启用,2禁用"`           // 用户状态（激活、冻结）// 菜单状态（启用、禁用）
	ParentID    int64  `json:"parent_id" gorm:"index;comment:父级ID;"`                               // 父级ID
	ParentPath  string `json:"parent_path" gorm:"size:255;index;comment:父级路径（用 . 分隔）;"`            // 父级路径（用 . 分隔）
}

func (a Permissions) TableName() string {
	return "sys_permissions"
}

// GetPrimaryKey ， 定义表主键 base repo 会使用，非 gorm 原生接口
// 参数：
// 返回值：
//
//	string ：表主键
func (a Permissions) GetPrimaryKey() string {
	return "id"
}
