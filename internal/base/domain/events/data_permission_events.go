package events

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	DataPermissionAssigned = "data_permission.assigned"
	DataPermissionRemoved  = "data_permission.removed"
)

// DataPermissionEvent 数据权限事件基类
type DataPermissionEvent struct {
	events.BaseEvent
	Permission *model.DataPermission `json:"permission"`
	TenantID   string                `json:"tenantID"`
}

// NewDataPermissionEvent 创建数据权限事件
func NewDataPermissionEvent(tenantID string, permission *model.DataPermission, eventType string) *DataPermissionEvent {
	return &DataPermissionEvent{
		BaseEvent:  events.NewBaseEvent(eventType),
		Permission: permission,
		TenantID:   tenantID,
	}
}
