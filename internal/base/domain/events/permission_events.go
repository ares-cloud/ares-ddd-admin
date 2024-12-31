package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 权限事件类型定义
const (
	PermissionCreated  = "permission.created"
	PermissionUpdated  = "permission.updated"
	PermissionDeleted  = "permission.deleted"
	PermissionDisabled = "permission.disabled"
	PermissionEnabled  = "permission.enabled"
	PermissionMoved    = "permission.moved"
	ResourceAssigned   = "permission.resource.assigned"
	ResourceRemoved    = "permission.resource.removed"
)

// PermissionEvent 权限事件基类
type PermissionEvent struct {
	events.BaseEvent
	tenantID string
	permID   int64
}

// NewPermissionEvent 创建权限事件
func NewPermissionEvent(tenantID string, permID int64, eventName string) *PermissionEvent {
	return &PermissionEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		tenantID:  tenantID,
		permID:    permID,
	}
}

func (e *PermissionEvent) PermissionID() int64 {
	return e.permID
}

func (e *PermissionEvent) TenantID() string {
	return e.tenantID
}

// PermissionCreatedEvent 权限创建事件
type PermissionCreatedEvent struct {
	PermissionEvent
}

// NewPermissionCreatedEvent 创建权限创建事件
func NewPermissionCreatedEvent(tenantID string, permID int64) *PermissionCreatedEvent {
	return &PermissionCreatedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, PermissionCreated),
	}
}

// PermissionUpdatedEvent 权限更新事件
type PermissionUpdatedEvent struct {
	PermissionEvent
}

// NewPermissionUpdatedEvent 创建权限更新事件
func NewPermissionUpdatedEvent(tenantID string, permID int64) *PermissionUpdatedEvent {
	return &PermissionUpdatedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, PermissionUpdated),
	}
}

// PermissionDeletedEvent 权限删除事件
type PermissionDeletedEvent struct {
	PermissionEvent
}

// NewPermissionDeletedEvent 创建权限删除事件
func NewPermissionDeletedEvent(tenantID string, permID int64) *PermissionDeletedEvent {
	return &PermissionDeletedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, PermissionDeleted),
	}
}

// PermissionMovedEvent 权限移动事件
type PermissionMovedEvent struct {
	PermissionEvent
	FromParentID int64
	ToParentID   int64
}

// NewPermissionMovedEvent 创建权限移动事件
func NewPermissionMovedEvent(tenantID string, permID int64, fromParentID, toParentID int64) *PermissionMovedEvent {
	return &PermissionMovedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, PermissionMoved),
		FromParentID:    fromParentID,
		ToParentID:      toParentID,
	}
}

// PermissionStatusChangedEvent 权限状态变更事件
type PermissionStatusChangedEvent struct {
	PermissionEvent
	OldStatus int8
	NewStatus int8
}

// NewPermissionStatusChangedEvent 创建权限状态变更事件
func NewPermissionStatusChangedEvent(tenantID string, permID int64, oldStatus, newStatus int8) *PermissionStatusChangedEvent {
	return &PermissionStatusChangedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, PermissionUpdated),
		OldStatus:       oldStatus,
		NewStatus:       newStatus,
	}
}

// ResourceAssignedEvent 资源分配事件
type ResourceAssignedEvent struct {
	PermissionEvent
	Method string
	Path   string
}

// NewResourceAssignedEvent 创建资源分配事件
func NewResourceAssignedEvent(tenantID string, permID int64, method, path string) *ResourceAssignedEvent {
	return &ResourceAssignedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, ResourceAssigned),
		Method:          method,
		Path:            path,
	}
}

// ResourceRemovedEvent 资源移除事件
type ResourceRemovedEvent struct {
	PermissionEvent
	Method string
	Path   string
}

// NewResourceRemovedEvent 创建资源移除事件
func NewResourceRemovedEvent(tenantID string, permID int64, method, path string) *ResourceRemovedEvent {
	return &ResourceRemovedEvent{
		PermissionEvent: *NewPermissionEvent(tenantID, permID, ResourceRemoved),
		Method:          method,
		Path:            path,
	}
}
