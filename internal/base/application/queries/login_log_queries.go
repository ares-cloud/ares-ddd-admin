package queries

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ListLoginLogsQuery struct {
	query.Page
	Month     string `json:"month" query:"month"`           // 查询月份(格式:202403)
	Username  string `json:"username" query:"username"`     // 用户名
	IP        string `json:"ip" query:"ip"`                 // 登录IP
	Status    int8   `json:"status" query:"status"`         // 登录状态
	StartTime int64  `json:"start_time" query:"start_time"` // 开始时间
	EndTime   int64  `json:"end_time" query:"end_time"`     // 结束时间
}
