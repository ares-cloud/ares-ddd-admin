package handlers

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/handlers"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type HandlerEvent struct {
	queryCache *handlers.EventHandler
	uh         *UserEventHandler
	eventBus   pkgEvent.IEventBus
}

func NewHandlerEvent(eventBus pkgEvent.IEventBus, queryCache *handlers.EventHandler, uh *UserEventHandler) *HandlerEvent {
	return &HandlerEvent{
		queryCache: queryCache,
		uh:         uh,
		eventBus:   eventBus,
	}
}

func (h *HandlerEvent) Register() {
	// 注册用户相关事件
	h.eventBus.Subscribe(events.UserCreated, h.uh)
	h.eventBus.Subscribe(events.UserUpdated, h.uh)
	h.eventBus.Subscribe(events.UserDeleted, h.uh)
	h.eventBus.Subscribe(events.UserRoleChanged, h.uh)

	// 注册缓存相关事件
	// 用户事件
	h.eventBus.Subscribe(events.UserLoggedIn, h.queryCache)
	h.eventBus.Subscribe(events.UserCreated, h.queryCache)
	h.eventBus.Subscribe(events.UserUpdated, h.queryCache)
	h.eventBus.Subscribe(events.UserDeleted, h.queryCache)
	h.eventBus.Subscribe(events.UserRoleChanged, h.queryCache)

	// 角色事件
	h.eventBus.Subscribe(events.RoleCreated, h.queryCache)
	h.eventBus.Subscribe(events.RoleUpdated, h.queryCache)
	h.eventBus.Subscribe(events.RoleDeleted, h.queryCache)
	h.eventBus.Subscribe(events.RolePermissionsChanged, h.queryCache)

	// 部门事件
	h.eventBus.Subscribe(events.DepartmentCreated, h.queryCache)
	h.eventBus.Subscribe(events.DepartmentUpdated, h.queryCache)
	h.eventBus.Subscribe(events.DepartmentDeleted, h.queryCache)
	h.eventBus.Subscribe(events.DepartmentMoved, h.queryCache)
	h.eventBus.Subscribe(events.UserAssigned, h.queryCache)
	h.eventBus.Subscribe(events.UserRemoved, h.queryCache)
	h.eventBus.Subscribe(events.UserTransferred, h.queryCache)

	// 权限事件
	h.eventBus.Subscribe(events.PermissionCreated, h.queryCache)
	h.eventBus.Subscribe(events.PermissionUpdated, h.queryCache)
	h.eventBus.Subscribe(events.PermissionDeleted, h.queryCache)
	h.eventBus.Subscribe(events.PermissionStatusChange, h.queryCache)

	// 数据权限事件
	h.eventBus.Subscribe(events.DataPermissionAssigned, h.queryCache)
	h.eventBus.Subscribe(events.DataPermissionRemoved, h.queryCache)

	// 租户事件
	h.eventBus.Subscribe(events.TenantCreated, h.queryCache)
	h.eventBus.Subscribe(events.TenantUpdated, h.queryCache)
	h.eventBus.Subscribe(events.TenantDeleted, h.queryCache)
	h.eventBus.Subscribe(events.TenantLocked, h.queryCache)
	h.eventBus.Subscribe(events.TenantUnlocked, h.queryCache)
}
