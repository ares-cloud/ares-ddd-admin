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
	RoleID   int64  `json:"role_id"`
	TenantID string `json:"tenant_id"`
}

// NewDataPermissionEvent 创建数据权限事件
func NewDataPermissionEvent(tenantID string, roleID int64, eventType string) *DataPermissionEvent {
	return &DataPermissionEvent{
		BaseEvent: events.NewBaseEvent(eventType),
		RoleID:    roleID,
		TenantID:  tenantID,
	}
}

// DataPermissionAssignedEvent 数据权限分配事件
type DataPermissionAssignedEvent struct {
	*DataPermissionEvent
	Permission *model.DataPermission `json:"permission"`
}

func NewDataPermissionAssignedEvent(permission *model.DataPermission) *DataPermissionAssignedEvent {
	return &DataPermissionAssignedEvent{
		DataPermissionEvent: NewDataPermissionEvent(permission.TenantID, permission.RoleID, DataPermissionAssigned),
		Permission:          permission,
	}
}

// DataPermissionRemovedEvent 数据权限移除事件
type DataPermissionRemovedEvent struct {
	*DataPermissionEvent
}

func NewDataPermissionRemovedEvent(tenantID string, roleID int64) *DataPermissionRemovedEvent {
	return &DataPermissionRemovedEvent{
		DataPermissionEvent: NewDataPermissionEvent(tenantID, roleID, DataPermissionRemoved),
	}
}
