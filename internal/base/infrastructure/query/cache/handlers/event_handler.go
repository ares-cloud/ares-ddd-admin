package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// EventHandler 缓存事件处理器
type EventHandler struct {
	userCache     *cache.UserQueryCache
	roleCache     *cache.RoleQueryCache
	deptCache     *cache.DepartmentQueryCache
	permCache     *cache.PermissionsQueryCache
	dataPermCache *cache.DataPermissionQueryCache
	tenantCache   *cache.TenantQueryCache
}

func NewCacheEventHandler(
	userCache *cache.UserQueryCache,
	roleCache *cache.RoleQueryCache,
	deptCache *cache.DepartmentQueryCache,
	permCache *cache.PermissionsQueryCache,
	dataPermCache *cache.DataPermissionQueryCache,
	tenantCache *cache.TenantQueryCache,
) *EventHandler {
	return &EventHandler{
		userCache:     userCache,
		roleCache:     roleCache,
		deptCache:     deptCache,
		permCache:     permCache,
		dataPermCache: dataPermCache,
		tenantCache:   tenantCache,
	}
}

// Handle 处理事件
func (h *EventHandler) Handle(ctx context.Context, event pkgEvent.Event) error {
	switch e := event.(type) {
	// 用户相关事件
	case *events.UserEvent:
		return h.handleUserEvent(ctx, e)
	case *events.UserRoleEvent:
		return h.handleUserRoleEvent(ctx, e)
	case *events.UserLoginEvent:
		return h.handleUserLoginEvent(ctx, e)

	// 角色相关事件
	case *events.RoleEvent:
		return h.handleRoleEvent(ctx, e)

	// 部门相关事件
	case *events.DepartmentEvent:
		return h.handleDepartmentEvent(ctx, e)

	// 权限相关事件
	case *events.PermissionEvent:
		return h.handlePermissionEvent(ctx, e)

	// 数据权限相关事件
	case *events.DataPermissionEvent:
		return h.handleDataPermissionEvent(ctx, e)

	// 租户相关事件
	case *events.TenantEvent:
		return h.handleTenantEvent(ctx, e)
	case *events.TenantPermissionEvent:
		return h.handleTenantPermissionEvent(ctx, e)

	default:
		return nil
	}
}

// 用户相关事件处理
func (h *EventHandler) handleUserEvent(ctx context.Context, event *events.UserEvent) error {
	log.Printf("处理用户事件: %s, 用户ID=%s", event.EventName(), event.UserID())

	// 清除用户基本信息缓存
	if err := h.userCache.InvalidateUserCache(ctx, event.UserID()); err != nil {
		return fmt.Errorf("清除用户缓存失败: %w", err)
	}

	// 特定事件的额外处理
	switch event.EventName() {
	case events.UserDeleted, events.UserStatusChanged:
		// 清除用户列表缓存
		if err := h.userCache.InvalidateUserListCache(ctx); err != nil {
			return fmt.Errorf("清除用户列表缓存失败: %w", err)
		}
		// 清除部门用户列表缓存
		//if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DepartmentID()); err != nil {
		//	return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
		//}
	case events.UserPasswordChanged:
		// 密码修改只需要清除用户基本信息缓存，已在上面处理
		break
	}
	return nil
}

// 用户调动事件处理
func (h *EventHandler) handleUserTransferEvent(ctx context.Context, event *events.UserTransferEvent) error {
	log.Printf("处理用户调动事件: 用户ID=%s, 从部门%s到部门%s", event.UserID(), event.FromDept, event.ToDept)

	// 1. 清除用户部门缓存
	if err := h.userCache.InvalidateUserDepartmentCache(ctx, event.UserID()); err != nil {
		return fmt.Errorf("清除用户部门缓存失败: %w", err)
	}

	// 2. 清除原部门和新部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.FromDept); err != nil {
		return fmt.Errorf("清除原部门用户列表缓存失败: %w", err)
	}
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.ToDept); err != nil {
		return fmt.Errorf("清除新部门用户列表缓存失败: %w", err)
	}

	return nil
}

// 角色相关事件处理
func (h *EventHandler) handleRoleEvent(ctx context.Context, event *events.RoleEvent) error {
	log.Printf("处理角色事件: %s, 角色ID=%d", event.EventName(), event.RoleID)

	// 清除角色缓存
	if err := h.roleCache.InvalidateRoleCache(ctx, event.RoleID); err != nil {
		return fmt.Errorf("清除角色缓存失败: %w", err)
	}

	// 清除角色列表缓存
	if err := h.roleCache.InvalidateRoleListCache(ctx); err != nil {
		return fmt.Errorf("清除角色列表缓存失败: %w", err)
	}

	return nil
}

// 部门相关事件处理
func (h *EventHandler) handleDepartmentEvent(ctx context.Context, event *events.DepartmentEvent) error {
	log.Printf("处理部门事件: %s, 部门ID=%s", event.EventName(), event.DepartmentID())

	// 清除部门基本信息缓存
	if err := h.deptCache.InvalidateCache(ctx, event.DepartmentID()); err != nil {
		return fmt.Errorf("清除部门缓存失败: %w", err)
	}

	// 特定事件的额外处理
	switch event.EventName() {
	case events.DepartmentMoved:
		//if moveEvent, ok := event.(*events.DepartmentMovedEvent); ok {
		//
		//}
		// 清除原父部门和新父部门的子部门列表缓存
		if err := h.deptCache.InvalidateChildrenCache(ctx, event.DeptID); err != nil {
			return fmt.Errorf("清除原父部门子部门列表缓存失败: %w", err)
		}
		if err := h.deptCache.InvalidateChildrenCache(ctx, event.DeptID); err != nil {
			return fmt.Errorf("清除新父部门子部门列表缓存失败: %w", err)
		}
	case events.DepartmentDeleted:
		// 清除部门树缓存
		if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
			return fmt.Errorf("清除部门树缓存失败: %w", err)
		}
	}

	return nil
}

// 权限相关事件处理
func (h *EventHandler) handlePermissionEvent(ctx context.Context, event *events.PermissionEvent) error {
	log.Printf("处理权限事件: %s, 权限ID=%d", event.EventName(), event.PermissionID())

	// 清除权限基本信息缓存
	if err := h.permCache.InvalidatePermissionCache(ctx, event.PermissionID()); err != nil {
		return fmt.Errorf("清除权限缓存失败: %w", err)
	}

	// 特定事件的额外处理
	switch event.EventName() {
	case events.PermissionMoved:
		// 清除原父权限和新父权限的子权限列表缓存
		//if err := h.permCache.InvalidateChildrenCache(ctx, event.FromParentID); err != nil {
		//	return fmt.Errorf("清除原父权限子权限列表缓存失败: %w", err)
		//}
		//if err := h.permCache.InvalidateChildrenCache(ctx, event.ToParentID); err != nil {
		//	return fmt.Errorf("清除新父权限子权限列表缓存失败: %w", err)
		//}
	case events.PermissionEnabled:
		// 清除权限树缓存
		if err := h.permCache.InvalidatePermissionTreeCache(ctx); err != nil {
			return fmt.Errorf("清除权限树缓存失败: %w", err)
		}
	case events.ResourceAssigned, events.ResourceRemoved:
		// 清除权限资源缓存
		if err := h.permCache.InvalidatePermissionResourceCache(ctx, event.PermissionID()); err != nil {
			return fmt.Errorf("清除权限资源缓存失败: %w", err)
		}
	}

	return nil
}

// 数据权限相关事件处理
func (h *EventHandler) handleDataPermissionEvent(ctx context.Context, event *events.DataPermissionEvent) error {
	log.Printf("处理数据权限事件: %s, 角色ID=%d", event.EventName(), event.RoleID)

	// 清除角色的数据权限缓存
	if err := h.dataPermCache.InvalidateCache(ctx, event.RoleID); err != nil {
		return fmt.Errorf("清除数据权限缓存失败: %w", err)
	}

	// 如果是分配事件，还需要清除相关用户的权限缓存
	switch event.EventName() {
	case events.DataPermissionAssigned:
		// 清除角色下所有用户的权限缓存
		users, err := h.roleCache.GetRoleUsers(ctx, event.RoleID)
		if err != nil {
			return fmt.Errorf("获取角色用户列表失败: %w", err)
		}
		for _, user := range users {
			if err := h.userCache.InvalidateUserPermissionCache(ctx, user.ID); err != nil {
				return fmt.Errorf("清除用户[%s]权限缓存失败: %w", user.ID, err)
			}
		}
	}

	return nil
}

// 租户相关事件处理
func (h *EventHandler) handleTenantEvent(ctx context.Context, event *events.TenantEvent) error {
	log.Printf("处理租户事件: %s, 租户ID=%s", event.EventName(), event.TenantID())

	// 清除租户基本信息缓存
	if err := h.tenantCache.InvalidateTenantCache(ctx, event.TenantID()); err != nil {
		return fmt.Errorf("清除租户缓存失败: %w", err)
	}

	// 特定事件的额外处理
	switch event.EventName() {
	case events.TenantDeleted:
		// 清除租户下所有角色的缓存
		roles, err := h.roleCache.GetTenantRoles(ctx, event.TenantID())
		if err != nil {
			return fmt.Errorf("获取租户角色列表失败: %w", err)
		}
		for _, role := range roles {
			if err := h.roleCache.InvalidateRoleCache(ctx, role.ID); err != nil {
				return fmt.Errorf("清除角色[%d]缓存失败: %w", role.ID, err)
			}
		}
	case events.TenantLocked, events.TenantUnlocked:
		// 清除租户状态缓存
		if err := h.tenantCache.InvalidateTenantStatusCache(ctx, event.TenantID()); err != nil {
			return fmt.Errorf("清除租户状态缓存失败: %w", err)
		}
	}

	return nil
}

// 用户角色变更事件处理
func (h *EventHandler) handleUserRoleEvent(ctx context.Context, event *events.UserRoleEvent) error {
	log.Printf("处理用户角色变更事件: 用户ID=%s, 角色IDs=%v", event.UserID(), event.RoleIDs())

	// 1. 清除用户的权限缓存
	if err := h.userCache.InvalidateUserPermissionCache(ctx, event.UserID()); err != nil {
		return fmt.Errorf("清除用户权限缓存失败: %w", err)
	}

	// 2. 清除用户的菜单缓存
	if err := h.userCache.InvalidateUserMenuCache(ctx, event.UserID()); err != nil {
		return fmt.Errorf("清除用户菜单缓存失败: %w", err)
	}

	return nil
}

// 用户登录事件处理
func (h *EventHandler) handleUserLoginEvent(ctx context.Context, event *events.UserLoginEvent) error {
	log.Printf("处理用户登录事件: 用户ID=%s", event.UserID())

	// 预热用户相关缓存
	if err := h.userCache.WarmupUserCache(ctx, event.UserID()); err != nil {
		return fmt.Errorf("预热用户缓存失败: %w", err)
	}

	return nil
}

// 角色权限变更事件处理
func (h *EventHandler) handleRolePermissionEvent(ctx context.Context, event *events.RolePermissionsAssignedEvent) error {
	log.Printf("处理角色权限变更事件: 角色ID=%d", event.RoleID)

	// 1. 清除角色的权限缓存
	if err := h.roleCache.InvalidateRolePermissionCache(ctx, event.RoleID); err != nil {
		return fmt.Errorf("清除角色权限缓存失败: %w", err)
	}

	// 2. 清除该角色下所有用户的权限缓存
	users, err := h.roleCache.GetRoleUsers(ctx, event.RoleID)
	if err != nil {
		return fmt.Errorf("获取角色用户列表失败: %w", err)
	}
	for _, user := range users {
		if err := h.userCache.InvalidateUserPermissionCache(ctx, user.ID); err != nil {
			return fmt.Errorf("清除用户[%s]权限缓存失败: %w", user.ID, err)
		}
		if err := h.userCache.InvalidateUserMenuCache(ctx, user.ID); err != nil {
			return fmt.Errorf("清除用户[%s]菜单缓存失败: %w", user.ID, err)
		}
	}

	return nil
}

// 部门用户变更事件处理
func (h *EventHandler) handleDepartmentUserEvent(ctx context.Context, event *events.UserTransferredEvent) error {
	log.Printf("处理部门用户变更事件: 部门ID=%s, 用户IDs=%v", event.DepartmentID(), event)

	// 1. 清除部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DepartmentID()); err != nil {
		return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
	}

	// 2. 清除相关用户的部门缓存
	if err := h.userCache.InvalidateUserDepartmentCache(ctx, event.UserID); err != nil {
		return fmt.Errorf("清除用户[%s]部门缓存失败: %w", event.UserID, err)
	}

	return nil
}

// 租户权限变更事件处理
func (h *EventHandler) handleTenantPermissionEvent(ctx context.Context, event *events.TenantPermissionEvent) error {
	log.Printf("处理租户权限变更事件: 租户ID=%s", event.TenantID())

	// 1. 清除租户的权限缓存
	if err := h.tenantCache.InvalidateTenantPermissionCache(ctx, event.TenantID()); err != nil {
		return fmt.Errorf("清除租户权限缓存失败: %w", err)
	}

	// 2. 清除租户下所有角色的权限缓存
	roles, err := h.roleCache.GetTenantRoles(ctx, event.TenantID())
	if err != nil {
		return fmt.Errorf("获取租户角色列表失败: %w", err)
	}
	for _, role := range roles {
		if err := h.roleCache.InvalidateRolePermissionCache(ctx, role.ID); err != nil {
			return fmt.Errorf("清除角色[%d]权限缓存失败: %w", role.ID, err)
		}
	}

	return nil
}
