package events

import (
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
