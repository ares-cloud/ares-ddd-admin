package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	TenantCreated  = "tenant.created"
	TenantUpdated  = "tenant.updated"
	TenantDeleted  = "tenant.deleted"
	TenantLocked   = "tenant.locked"
	TenantUnlocked = "tenant.unlocked"
)

// TenantEvent 租户事件基类
type TenantEvent struct {
	events.BaseEvent
	TenantID string `json:"tenant_id"`
}

// NewTenantEvent 创建租户事件
func NewTenantEvent(tenantID string, eventName string) *TenantEvent {
	return &TenantEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		TenantID:  tenantID,
	}
}

// TenantPermissionEvent 租户权限变更事件
type TenantPermissionEvent struct {
	*TenantEvent
	PermissionIDs []int64
}

func NewTenantPermissionEvent(tenantID string, permissionIDs []int64) *TenantPermissionEvent {
	return &TenantPermissionEvent{
		TenantEvent:   NewTenantEvent(tenantID, TenantUpdated),
		PermissionIDs: permissionIDs,
	}
}
