package entity

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

// LoginLog 登录日志实体
type LoginLog struct {
	database.BaseIntTime
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`
	UserID    string `json:"user_id" gorm:"type:varchar(64);index:idx_user_id;comment:用户ID"`
	Username  string `json:"username" gorm:"type:varchar(64);index:idx_username;comment:用户名"`
	TenantID  string `json:"tenant_id" gorm:"type:varchar(64);index:idx_tenant_id;comment:租户ID"`
	LoginType int8   `json:"login_type" gorm:"type:smallint;default:1;comment:登录类型(1:管理端 2:前台)"`
	IP        string `json:"ip" gorm:"type:varchar(64);comment:登录IP"`
	Location  string `json:"location" gorm:"type:varchar(128);comment:登录地点"`
	Device    string `json:"device" gorm:"type:varchar(128);comment:登录设备"`
	OS        string `json:"os" gorm:"type:varchar(64);comment:操作系统"`
	Browser   string `json:"browser" gorm:"type:varchar(600);comment:浏览器"`
	Status    int8   `json:"status" gorm:"type:smallint;default:1;comment:登录状态(1:成功 2:失败)"`
	Message   string `json:"message" gorm:"type:varchar(255);comment:登录消息"`
	LoginTime int64  `json:"login_time" gorm:"index:idx_login_time;comment:登录时间"`
}
