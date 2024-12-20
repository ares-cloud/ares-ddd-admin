package oplog

import (
	"context"
	"time"
)

// OperationLog 操作日志结构
type OperationLog struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`    // 操作人ID
	Username  string    `json:"username"`   // 操作人用户名
	TenantID  string    `json:"tenant_id"`  // 租户ID
	Method    string    `json:"method"`     // 请求方法
	Path      string    `json:"path"`       // 请求路径
	Query     string    `json:"query"`      // 查询参数
	Body      string    `json:"body"`       // 请求体
	IP        string    `json:"ip"`         // 请求IP
	UserAgent string    `json:"user_agent"` // 用户代理
	Status    int       `json:"status"`     // 响应状态码
	Error     string    `json:"error"`      // 错误信息
	Duration  int64     `json:"duration"`   // 执行时长(ms)
	CreatedAt time.Time `json:"created_at"` // 创建时间
	Module    string    `json:"module"`     // 模块名称
	Action    string    `json:"action"`     // 操作类型
}

// LogOption 日志选项
type LogOption struct {
	IncludeBody bool   // 是否记录请求体
	Module      string // 模块名称
	Action      string // 操作类型
}

type IDbOperationLogWrite interface {
	Save(ctx context.Context, data *OperationLog, tableName string) error
	ChickTableExist(ctx context.Context, tableName string) (bool, error)
}
