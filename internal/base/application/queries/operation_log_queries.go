package queries

import "github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

// ListOperationLogQuery 查询操作日志列表
type ListOperationLogQuery struct {
	query.Page
	Month     string `json:"month" query:"month"`           // 查询月份(格式:202403)
	TenantID  string `json:"tenant_id" query:"tenant_id"`   // 租户ID
	Username  string `json:"username" query:"username"`     // 用户名
	Module    string `json:"module" query:"module"`         // 模块
	Action    string `json:"action" query:"action"`         // 操作类型
	StartTime int64  `json:"start_time" query:"start_time"` // 开始时间
	EndTime   int64  `json:"end_time" query:"end_time"`     // 结束时间
}
