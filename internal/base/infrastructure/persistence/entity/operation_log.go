package entity

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

const (
	OperationTypeQuery  = 1 // 查询
	OperationTypeCreate = 2 // 创建
	OperationTypeUpdate = 3 // 更新
	OperationTypeDelete = 4 // 删除
)

// OperationLog 操作日志实体
type OperationLog struct {
	database.BaseIntTime
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`
	UserID    string `json:"user_id" gorm:"type:varchar(64);index:idx_user_id;comment:用户ID"`
	Username  string `json:"username" gorm:"type:varchar(64);index:idx_username;comment:用户名"`
	TenantID  string `json:"tenant_id" gorm:"type:varchar(64);index:idx_tenant_id;comment:租户ID"`
	Method    string `json:"method" gorm:"type:varchar(10);comment:请求方法"`
	Path      string `json:"path" gorm:"type:varchar(255);comment:请求路径"`
	Query     string `json:"query" gorm:"type:text;comment:查询参数"`
	Body      string `json:"body" gorm:"type:text;comment:请求体"`
	IP        string `json:"ip" gorm:"type:varchar(64);comment:请求IP"`
	UserAgent string `json:"user_agent" gorm:"type:varchar(600);comment:用户代理"`
	Status    int    `json:"status" gorm:"type:int;comment:响应状态码"`
	Error     string `json:"error" gorm:"type:text;comment:错误信息"`
	Duration  int64  `json:"duration" gorm:"type:bigint;comment:执行时长(ms)"`
	Module    string `json:"module" gorm:"type:varchar(64);comment:模块名称"`
	Action    string `json:"action" gorm:"type:varchar(32);comment:操作类型"`
}
