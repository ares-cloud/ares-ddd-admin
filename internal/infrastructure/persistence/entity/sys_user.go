package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// SysUser 系统用户
type SysUser struct {
	database.BaseModel
	ID             string `json:"id" gorm:"primaryKey;size:32;comment:用户ID"`
	TenantID       string `json:"tenant_id" gorm:"size:32;index;comment:租户ID"`
	Username       string `json:"username" gorm:"size:32;uniqueIndex;comment:用户名"`
	FaceURL        string `json:"face_url" gorm:"size:255;comment:头像"`
	Name           string `json:"name" gorm:"size:128;comment:姓名"`
	Password       string `json:"password" gorm:"size:128;comment:密码"`
	Phone          string `json:"phone" gorm:"size:32;comment:手机号"`
	Email          string `json:"email" gorm:"size:128;comment:邮箱"`
	Remark         string `json:"remark" gorm:"size:512;comment:备注"`
	InvitationCode string `json:"invitation_code" gorm:"size:32;comment:邀请码"`
	Status         int8   `json:"status" gorm:"column:status;default:1;comment:状态,1启用,2禁用"`
}

// TableName 定义数据库中用户表的名称
func (a SysUser) TableName() string {
	return "sys_user"
}

// GetPrimaryKey ， 定义表主键 base repo 会使用，非 gorm 原生接口
// 参数：
// 返回值：
//
//	string ：表主键
func (a SysUser) GetPrimaryKey() string {
	return "id"
}
