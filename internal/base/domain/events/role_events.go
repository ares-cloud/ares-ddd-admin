package events

import (
	"strconv"

	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	// 角色聚合根名称
	RoleAggregate = "role"
	// 版本号
	Version = "1.0"
)

// 角色事件类型定义
const (
	RoleCreated           = "role.created"
	RoleUpdated           = "role.updated"
	RoleDeleted           = "role.deleted"
	RolePermissionChanged = "role.permission.changed"
	RoleDataScopeChanged  = "role.datascope.changed"
)

// RoleEvent 角色事件
type RoleEvent struct {
	events.BaseTenantEvent
}

// NewRoleEvent 创建角色事件
func NewRoleEvent(tenantID string, roleID int64, eventName string) *RoleEvent {
	return &RoleEvent{
		BaseTenantEvent: events.NewBaseTenantEvent(
			eventName,
			Version,
			strconv.FormatInt(roleID, 10),
			RoleAggregate,
			tenantID,
		),
	}
}

func (e *RoleEvent) RoleID() string {
	return e.AggregateID()
}

// RolePermissionEvent 角色权限变更事件
type RolePermissionEvent struct {
	RoleEvent
	permissionIDs []int64
}

// NewRolePermissionEvent 创建角色权限变更事件
func NewRolePermissionEvent(tenantID string, roleID int64, permissionIDs []int64) *RolePermissionEvent {
	return &RolePermissionEvent{
		RoleEvent:     *NewRoleEvent(tenantID, roleID, RolePermissionChanged),
		permissionIDs: permissionIDs,
	}
}

func (e *RolePermissionEvent) PermissionIDs() []int64 {
	return e.permissionIDs
}
