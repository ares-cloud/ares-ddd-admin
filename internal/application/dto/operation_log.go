package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
)

// OperationLogDto 操作日志DTO
type OperationLogDto struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`    // 操作人ID
	Username  string `json:"username"`   // 操作人用户名
	TenantID  string `json:"tenant_id"`  // 租户ID
	Method    string `json:"method"`     // 请求方法
	Path      string `json:"path"`       // 请求路径
	Query     string `json:"query"`      // 查询参数
	Body      string `json:"body"`       // 请求体
	IP        string `json:"ip"`         // 请求IP
	UserAgent string `json:"user_agent"` // 用户代理
	Status    int    `json:"status"`     // 响应状态码
	Error     string `json:"error"`      // 错误信息
	Duration  int64  `json:"duration"`   // 执行时长(ms)
	Module    string `json:"module"`     // 模块名称
	Action    string `json:"action"`     // 操作类型
	CreatedAt int64  `json:"created_at"` // 创建时间
}

// ToOperationLogDto 转换为DTO
func ToOperationLogDto(model *model.OperationLog) *OperationLogDto {
	return &OperationLogDto{
		ID:        model.ID,
		UserID:    model.UserID,
		Username:  model.Username,
		TenantID:  model.TenantID,
		Method:    model.Method,
		Path:      model.Path,
		Query:     model.Query,
		Body:      model.Body,
		IP:        model.IP,
		UserAgent: model.UserAgent,
		Status:    model.Status,
		Error:     model.Error,
		Duration:  model.Duration,
		Module:    model.Module,
		Action:    model.Action,
		CreatedAt: model.CreatedAt,
	}
}

// ToOperationLogDtoList 转换为DTO列表
func ToOperationLogDtoList(models []*model.OperationLog) []*OperationLogDto {
	dtos := make([]*OperationLogDto, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, ToOperationLogDto(m))
	}
	return dtos
}
