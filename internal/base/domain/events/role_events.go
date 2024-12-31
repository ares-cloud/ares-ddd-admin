package events

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	RoleCreated            = "role.created"
	RoleUpdated            = "role.updated"
	RoleDeleted            = "role.deleted"
	RolePermissionsChanged = "role.permissions.changed"
)

// RoleEvent 角色事件基类
type RoleEvent struct {
	events.BaseEvent
	RoleID   int64  `json:"role_id"`
	TenantID string `json:"tenant_id"`
}

// NewRoleEvent 创建角色事件
func NewRoleEvent(tenantID string, roleID int64, eventType string) *RoleEvent {
	return &RoleEvent{
		BaseEvent: events.NewBaseEvent(eventType),
		RoleID:    roleID,
		TenantID:  tenantID,
	}
}

// RoleCreatedEvent 角色创建事件
type RoleCreatedEvent struct {
	*RoleEvent
	Role *model.Role `json:"role"`
}

func NewRoleCreatedEvent(role *model.Role) *RoleCreatedEvent {
	return &RoleCreatedEvent{
		RoleEvent: NewRoleEvent(role.TenantID, role.ID, RoleCreated),
		Role:      role,
	}
}

// RoleUpdatedEvent 角色更新事件
type RoleUpdatedEvent struct {
	*RoleEvent
	Role *model.Role `json:"role"`
}

func NewRoleUpdatedEvent(role *model.Role) *RoleUpdatedEvent {
	return &RoleUpdatedEvent{
		RoleEvent: NewRoleEvent(role.TenantID, role.ID, RoleUpdated),
		Role:      role,
	}
}

// RoleDeletedEvent 角色删除事件
type RoleDeletedEvent struct {
	*RoleEvent
}

func NewRoleDeletedEvent(roleID int64) *RoleDeletedEvent {
	return &RoleDeletedEvent{
		RoleEvent: NewRoleEvent("", roleID, RoleDeleted),
	}
}

// RolePermissionsAssignedEvent 角色权限分配事件
type RolePermissionsAssignedEvent struct {
	*RoleEvent
	PermissionIDs []int64 `json:"permission_ids"`
}

func NewRolePermissionsAssignedEvent(roleID int64, permissionIDs []int64) *RolePermissionsAssignedEvent {
	return &RolePermissionsAssignedEvent{
		RoleEvent:     NewRoleEvent("", roleID, RolePermissionsChanged),
		PermissionIDs: permissionIDs,
	}
}
