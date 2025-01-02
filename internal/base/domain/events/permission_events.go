package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 权限事件类型定义
const (
	PermissionCreated      = "permission.created"
	PermissionUpdated      = "permission.updated"
	PermissionDeleted      = "permission.deleted"
	PermissionStatusChange = "permission.status.change"
)

// PermissionEvent 权限事件基类
type PermissionEvent struct {
	events.BaseEvent
	TenantID string `json:"tenant_id"`
	PermID   int64  `json:"perm_id"`
}

// NewPermissionEvent 创建权限事件
func NewPermissionEvent(tenantID string, permID int64, eventName string) *PermissionEvent {
	return &PermissionEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		TenantID:  tenantID,
		PermID:    permID,
	}
}
