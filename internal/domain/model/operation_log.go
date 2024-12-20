package model

import "time"

// OperationType 操作类型
type OperationType int8

const (
	OperationTypeQuery  OperationType = 1 // 查询
	OperationTypeCreate OperationType = 2 // 创建
	OperationTypeUpdate OperationType = 3 // 更新
	OperationTypeDelete OperationType = 4 // 删除
)

// OperationLog 操作日志领域模型
type OperationLog struct {
	ID        int64  // ID
	UserID    string // 操作人ID
	Username  string // 操作人用户名
	TenantID  string // 租户ID
	Method    string // 请求方法
	Path      string // 请求路径
	Query     string // 查询参数
	Body      string // 请求体
	IP        string // 请求IP
	UserAgent string // 用户代理
	Status    int    // 响应状态码
	Error     string // 错误信息
	Duration  int64  // 执行时长(ms)
	Module    string // 模块名称
	Action    string // 操作类型
	CreatedAt int64  // 创建时间
}

// NewOperationLog 创建操作日志
func NewOperationLog(userID, username, tenantID string) *OperationLog {
	return &OperationLog{
		UserID:    userID,
		Username:  username,
		TenantID:  tenantID,
		CreatedAt: time.Now().Unix(),
	}
}

// SetRequestInfo 设置请求信息
func (l *OperationLog) SetRequestInfo(method, path, query, body, ip, userAgent string) {
	l.Method = method
	l.Path = path
	l.Query = query
	l.Body = body
	l.IP = ip
	l.UserAgent = userAgent
}

// SetResponseInfo 设置响应信息
func (l *OperationLog) SetResponseInfo(status int, errMsg string, duration int64) {
	l.Status = status
	l.Error = errMsg
	l.Duration = duration
}

// SetModuleInfo 设置模块信息
func (l *OperationLog) SetModuleInfo(module, action string) {
	l.Module = module
	l.Action = action
}
