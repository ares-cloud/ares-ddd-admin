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
	tenantID string
}

func (e *TenantEvent) TenantID() string {
	return e.tenantID
}

// NewTenantEvent 创建租户事件
func NewTenantEvent(tenantID string, eventName string) *TenantEvent {
	return &TenantEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		tenantID:  tenantID,
	}
}

// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct {
	*TenantEvent
}

func NewTenantCreatedEvent(tenantID string) *TenantCreatedEvent {
	return &TenantCreatedEvent{
		TenantEvent: NewTenantEvent(tenantID, TenantCreated),
	}
}

// TenantPermissionEvent 租户权限变更事件
type TenantPermissionEvent struct {
	*TenantEvent
	permissionIDs []int64
}

func NewTenantPermissionEvent(tenantID string, permissionIDs []int64) *TenantPermissionEvent {
	return &TenantPermissionEvent{
		TenantEvent:   NewTenantEvent(tenantID, TenantUpdated),
		permissionIDs: permissionIDs,
	}
}

func (e *TenantPermissionEvent) PermissionIDs() []int64 {
	return e.permissionIDs
}

// TenantUpdatedEvent 租户更新事件
type TenantUpdatedEvent struct {
	*TenantEvent
}

func NewTenantUpdatedEvent(tenantID string) *TenantUpdatedEvent {
	return &TenantUpdatedEvent{
		TenantEvent: NewTenantEvent(tenantID, TenantUpdated),
	}
}

// TenantDeletedEvent 租户删除事件
type TenantDeletedEvent struct {
	*TenantEvent
}

func NewTenantDeletedEvent(tenantID string) *TenantDeletedEvent {
	return &TenantDeletedEvent{
		TenantEvent: NewTenantEvent(tenantID, TenantUpdated),
	}
}

// TenantLockedEvent 租户锁定事件
type TenantLockedEvent struct {
	*TenantEvent
	Reason string
}

// NewTenantLockedEvent 创建租户锁定事件
func NewTenantLockedEvent(tenantID string, reason string) *TenantLockedEvent {
	return &TenantLockedEvent{
		TenantEvent: NewTenantEvent(tenantID, TenantUpdated),
		Reason:      reason,
	}
}

// TenantUnlockedEvent 租户解锁事件
type TenantUnlockedEvent struct {
	*TenantEvent
}

// NewTenantUnlockedEvent 创建租户解锁事件
func NewTenantUnlockedEvent(tenantID string) *TenantUnlockedEvent {
	return &TenantUnlockedEvent{
		TenantEvent: NewTenantEvent(tenantID, TenantUpdated),
	}
}
