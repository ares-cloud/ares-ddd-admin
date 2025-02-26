package handlers

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"

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

	// 角色相关事件
	case *events.RoleEvent:
		return h.handleRoleEvent(ctx, e)
	case *events.RolePermissionsAssignedEvent:
		return h.handleRolePermissionsAssignedEvent(ctx, e)

	// 部门相关事件
	case *events.DepartmentEvent:
		return h.handleDepartmentEvent(ctx, e)
	case *events.DepartmentMovedEvent:
		return h.handleDepartmentMovedEvent(ctx, e)
	case *events.UserRemovedEvent:
		return h.handleUserRemovedEvent(ctx, e)
	case *events.UserAssignedEvent:
		return h.handleUserAssignedEvent(ctx, e)
	case *events.UserTransferredEvent:
		return h.handleUserTransferredEvent(ctx, e)

	// todo 权限变更清除关联用户和角色
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
	hlog.CtxDebugf(ctx, "处理用户事件: %s, 用户ID=%s", event.EventName(), event.UserID)

	switch event.EventName() {
	case events.UserCreated:
		// 1. 清除用户列表缓存
		if err := h.userCache.InvalidateUserListCache(ctx); err != nil {
			return fmt.Errorf("清除用户列表缓存失败: %w", err)
		}
		// 2. 清除部门用户列表缓存
		depts, err := h.userCache.GetUserDepartments(ctx, event.UserID)
		if err != nil {
			return fmt.Errorf("获取用户部门列表失败: %w", err)
		}
		for _, dept := range depts {
			if err := h.deptCache.InvalidateDepartmentUserCache(ctx, dept.ID); err != nil {
				return fmt.Errorf("清除部门[%s]用户列表缓存失败: %w", dept.ID, err)
			}
		}

	case events.UserUpdated:
		// 1. 清除用户相关缓存
		if err := h.userCache.InvalidateUserCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户缓存失败: %w", err)
		}
		// 2. 清除用户列表缓存
		if err := h.userCache.InvalidateUserListCache(ctx); err != nil {
			return fmt.Errorf("清除用户列表缓存失败: %w", err)
		}
		// 3. 清除用户权限和菜单缓存
		if err := h.userCache.InvalidateUserPermissionCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户权限缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserMenuCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户菜单缓存失败: %w", err)
		}
		// 4. 清除部门用户列表缓存
		depts, err := h.userCache.GetUserDepartments(ctx, event.UserID)
		if err != nil {
			return fmt.Errorf("获取用户部门列表失败: %w", err)
		}
		for _, dept := range depts {
			if err := h.deptCache.InvalidateDepartmentUserCache(ctx, dept.ID); err != nil {
				return fmt.Errorf("清除部门[%s]用户列表缓存失败: %w", dept.ID, err)
			}
		}

	case events.UserDeleted:
		// 1. 清除用户所有相关缓存
		if err := h.userCache.InvalidateUserCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserListCache(ctx); err != nil {
			return fmt.Errorf("清除用户列表缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserPermissionCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户权限缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserMenuCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户菜单缓存失败: %w", err)
		}
		// 2. 清除部门用户列表缓存
		depts, err := h.userCache.GetUserDepartments(ctx, event.UserID)
		if err != nil {
			return fmt.Errorf("获取用户部门列表失败: %w", err)
		}
		for _, dept := range depts {
			if err := h.deptCache.InvalidateDepartmentUserCache(ctx, dept.ID); err != nil {
				return fmt.Errorf("清除部门[%s]用户列表缓存失败: %w", dept.ID, err)
			}
		}
		// 3. 清除角色用户列表缓存
		roles, err := h.userCache.GetUserRoles(ctx, event.UserID)
		if err != nil {
			return fmt.Errorf("获取用户角色列表失败: %w", err)
		}
		for _, role := range roles {
			if err := h.roleCache.InvalidateRoleUserCache(ctx, role.ID); err != nil {
				return fmt.Errorf("清除角色[%d]用户列表缓存失败: %w", role.ID, err)
			}
		}

	case events.UserRoleChanged:
		// 1. 清除用户权限相关缓存
		if err := h.userCache.InvalidateUserCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserPermissionCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户权限缓存失败: %w", err)
		}
		if err := h.userCache.InvalidateUserMenuCache(ctx, event.UserID); err != nil {
			return fmt.Errorf("清除用户菜单缓存失败: %w", err)
		}
		// 2. 清除角色用户列表缓存
		roles, err := h.userCache.GetUserRoles(ctx, event.UserID)
		if err != nil {
			return fmt.Errorf("获取用户角色列表失败: %w", err)
		}
		for _, role := range roles {
			if err := h.roleCache.InvalidateRoleUserCache(ctx, role.ID); err != nil {
				return fmt.Errorf("清除角色[%d]用户列表缓存失败: %w", role.ID, err)
			}
		}

	case events.UserLoggedIn:
		// 预热用户缓存
		if err := h.userCache.WarmupUserCache(ctx, event.UserID); err != nil {
			hlog.CtxDebugf(ctx, "用户登录预热缓存失败: %v", err)
		}
	}

	return nil
}

// 角色相关事件处理
func (h *EventHandler) handleRoleEvent(ctx context.Context, event *events.RoleEvent) error {
	hlog.CtxDebugf(ctx, "处理角色事件: %s, 角色ID=%d", event.EventName(), event.RoleID)

	switch event.EventName() {
	case events.RoleCreated:
		// 清除角色列表缓存
		if err := h.roleCache.InvalidateRoleListCache(ctx); err != nil {
			return fmt.Errorf("清除角色列表缓存失败: %w", err)
		}

	case events.RoleUpdated:
		// 1. 清除角色缓存
		if err := h.roleCache.InvalidateRoleCache(ctx, event.RoleID); err != nil {
			return fmt.Errorf("清除角色缓存失败: %w", err)
		}
		// 2. 清除角色列表缓存
		if err := h.roleCache.InvalidateRoleListCache(ctx); err != nil {
			return fmt.Errorf("清除角色列表缓存失败: %w", err)
		}
		// 3. 清除该角色下所有用户的权限和菜单缓存
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

	case events.RoleDeleted:
		// 1. 清除角色缓存
		if err := h.roleCache.InvalidateRoleCache(ctx, event.RoleID); err != nil {
			return fmt.Errorf("清除角色缓存失败: %w", err)
		}
		// 2. 清除角色列表缓存
		if err := h.roleCache.InvalidateRoleListCache(ctx); err != nil {
			return fmt.Errorf("清除角色列表缓存失败: %w", err)
		}
		// 3. 清除角色权限缓存
		if err := h.roleCache.InvalidateRolePermissionCache(ctx, event.RoleID); err != nil {
			return fmt.Errorf("清除角色权限缓存失败: %w", err)
		}
		// 4. 清除该角色下所有用户的权限和菜单缓存
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
	}

	return nil
}

// 角色相关事件处理
func (h *EventHandler) handleRolePermissionsAssignedEvent(ctx context.Context, event *events.RolePermissionsAssignedEvent) error {
	hlog.CtxDebugf(ctx, "处理角色权限分配事件: 角色ID=%d", event.RoleID)

	// 1. 清除角色的权限缓存
	if err := h.roleCache.InvalidateRolePermissionCache(ctx, event.RoleID); err != nil {
		return fmt.Errorf("清除角色权限缓存失败: %w", err)
	}

	// 2. 清除该角色下所有用户的权限和菜单缓存
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

// 部门相关事件处理
func (h *EventHandler) handleDepartmentMovedEvent(ctx context.Context, event *events.DepartmentMovedEvent) error {
	hlog.CtxDebugf(ctx, "处理部门移动事件: 部门ID=%s, 原父部门ID=%s, 新父部门ID=%s",
		event.DeptID, event.FromParentID, event.ToParentID)

	// 1. 清除被移动部门的缓存
	if err := h.deptCache.InvalidateCache(ctx, event.DeptID); err != nil {
		return fmt.Errorf("清除部门缓存失败: %w", err)
	}

	// 2. 清除原父部门的子部门列表缓存
	if err := h.deptCache.InvalidateChildrenCache(ctx, event.FromParentID); err != nil {
		return fmt.Errorf("清除原父部门子部门列表缓存失败: %w", err)
	}

	// 3. 清除新父部门的子部门列表缓存
	if err := h.deptCache.InvalidateChildrenCache(ctx, event.ToParentID); err != nil {
		return fmt.Errorf("清除新父部门子部门列表缓存失败: %w", err)
	}

	// 4. 清除部门树缓存
	if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
		return fmt.Errorf("清除部门树缓存失败: %w", err)
	}

	return nil
}

func (h *EventHandler) handleUserRemovedEvent(ctx context.Context, event *events.UserRemovedEvent) error {
	hlog.CtxDebugf(ctx, "处理用户移除事件: 部门ID=%s, 用户ID=%s", event.DeptID, event.UserIDs)

	// 1. 清除部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DeptID); err != nil {
		return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
	}

	// 2. 清除用户的部门缓存
	for _, v := range event.UserIDs {
		if err := h.userCache.InvalidateUserDepartmentCache(ctx, v); err != nil {
			return fmt.Errorf("清除用户部门缓存失败: %w", err)
		}
	}

	return nil
}

func (h *EventHandler) handleUserAssignedEvent(ctx context.Context, event *events.UserAssignedEvent) error {
	hlog.CtxDebugf(ctx, "处理用户分配事件: 部门ID=%s, 用户ID=%s", event.DeptID, event.UserIDs)

	// 1. 清除部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DeptID); err != nil {
		return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
	}

	// 2. 清除用户的部门缓存
	for _, v := range event.UserIDs {
		if err := h.userCache.InvalidateUserDepartmentCache(ctx, v); err != nil {
			return fmt.Errorf("清除用户部门缓存失败: %w", err)
		}
	}
	return nil
}

func (h *EventHandler) handleUserTransferredEvent(ctx context.Context, event *events.UserTransferredEvent) error {
	hlog.CtxDebugf(ctx, "处理用户调动事件: 用户ID=%s, 原部门ID=%s, 新部门ID=%s",
		event.UserID, event.FromDeptID, event.ToDeptID)

	// 1. 清除原部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.FromDeptID); err != nil {
		return fmt.Errorf("清除原部门用户列表缓存失败: %w", err)
	}

	// 2. 清除新部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.ToDeptID); err != nil {
		return fmt.Errorf("清除新部门用户列表缓存失败: %w", err)
	}

	// 3. 清除用户的部门缓存
	if err := h.userCache.InvalidateUserDepartmentCache(ctx, event.UserID); err != nil {
		return fmt.Errorf("清除用户部门缓存失败: %w", err)
	}

	return nil
}

// 部门相关事件处理
func (h *EventHandler) handleDepartmentEvent(ctx context.Context, event *events.DepartmentEvent) error {
	hlog.CtxDebugf(ctx, "处理部门事件: %s, 部门ID=%s", event.EventName(), event.DeptID)

	switch event.EventName() {
	case events.DepartmentCreated:
		// 1. 清除部门树缓存
		if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
			return fmt.Errorf("清除部门树缓存失败: %w", err)
		}
		// 2. 清除父部门的子部门列表缓存
		if err := h.deptCache.InvalidateChildrenCache(ctx, ""); err != nil {
			return fmt.Errorf("清除根部门子部门列表缓存失败: %w", err)
		}

	case events.DepartmentUpdated:
		// 1. 清除部门基本信息缓存
		if err := h.deptCache.InvalidateCache(ctx, event.DeptID); err != nil {
			return fmt.Errorf("清除部门缓存失败: %w", err)
		}
		// 2. 清除部门树缓存
		if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
			return fmt.Errorf("清除部门树缓存失败: %w", err)
		}

	case events.DepartmentDeleted:
		// 1. 清除部门基本信息缓存
		if err := h.deptCache.InvalidateCache(ctx, event.DeptID); err != nil {
			return fmt.Errorf("清除部门缓存失败: %w", err)
		}
		// 2. 清除部门树缓存
		if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
			return fmt.Errorf("清除部门树缓存失败: %w", err)
		}
		// 3. 清除父部门的子部门列表缓存
		if err := h.deptCache.InvalidateChildrenCache(ctx, ""); err != nil {
			return fmt.Errorf("清除根部门子部门列表缓存失败: %w", err)
		}
		// 4. 清除部门用户列表缓存
		if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DeptID); err != nil {
			return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
		}
	}

	return nil
}

// 数据权限相关事件处理
func (h *EventHandler) handleDataPermissionEvent(ctx context.Context, event *events.DataPermissionEvent) error {
	hlog.CtxDebugf(ctx, "处理数据权限事件: %s, 角色ID=%d", event.EventName(), event.Permission.RoleID)

	// 清除角色的数据权限缓存
	if err := h.dataPermCache.InvalidateCache(ctx, event.Permission.RoleID); err != nil {
		return fmt.Errorf("清除数据权限缓存失败: %w", err)
	}

	// 如果是分配事件，还需要清除相关用户的权限缓存
	switch event.EventName() {
	case events.DataPermissionAssigned:
		// 清除角色下所有用户的权限缓存
		users, err := h.roleCache.GetRoleUsers(ctx, event.Permission.RoleID)
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

// 部门用户变更事件处理
func (h *EventHandler) handleDepartmentUserEvent(ctx context.Context, event *events.UserTransferredEvent) error {
	hlog.CtxDebugf(ctx, "处理部门用户变更事件: 部门ID=%s, 用户ID=%s", event.DeptID, event.UserID)

	// 1. 清除部门的用户列表缓存
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, event.DeptID); err != nil {
		return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
	}

	// 2. 清除相关用户的部门缓存
	if err := h.userCache.InvalidateUserDepartmentCache(ctx, event.UserID); err != nil {
		return fmt.Errorf("清除用户[%s]部门缓存失败: %w", event.UserID, err)
	}

	return nil
}

// 租户相关事件处理
func (h *EventHandler) handleTenantEvent(ctx context.Context, event *events.TenantEvent) error {
	hlog.CtxDebugf(ctx, "处理租户事件: %s, 租户ID=%s", event.EventName(), event.TenantID)

	// 根据事件类型处理
	switch event.EventName() {
	case events.TenantCreated:
		// 租户创建时清除租户列表缓存
		//if err := h.tenantCache.InvalidateTenantListCache(ctx); err != nil {
		//	return fmt.Errorf("清除租户列表缓存失败: %w", err)
		//}

	case events.TenantUpdated:
		// 租户更新时清除租户相关缓存
		if err := h.tenantCache.InvalidateTenantCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户缓存失败: %w", err)
		}
	case events.TenantDeleted:
		// 租户删除时清除所有相关缓存
		if err := h.tenantCache.InvalidateTenantCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户缓存失败: %w", err)
		}
		// 清除租户下所有用户缓存
		if err := h.userCache.InvalidateTenantUserCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户用户缓存失败: %w", err)
		}
		// 清除租户下所有角色缓存
		if err := h.roleCache.InvalidateTenantRoleCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户角色缓存失败: %w", err)
		}
		// 清除租户下所有部门缓存
		if err := h.deptCache.InvalidateTenantDepartmentCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户部门缓存失败: %w", err)
		}

	case events.TenantLocked, events.TenantUnlocked:
		// 租户锁定状态变更时清除租户状态缓存
		if err := h.tenantCache.InvalidateTenantStatusCache(ctx, event.TenantID); err != nil {
			return fmt.Errorf("清除租户状态缓存失败: %w", err)
		}
	}

	return nil
}

// 角色权限变更事件处理
func (h *EventHandler) handleRolePermissionEvent(ctx context.Context, event *events.RolePermissionsAssignedEvent) error {
	hlog.CtxDebugf(ctx, "处理角色权限变更事件: 角色ID=%d", event.RoleID)

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

// 租户权限变更事件处理
func (h *EventHandler) handleTenantPermissionEvent(ctx context.Context, event *events.TenantPermissionEvent) error {
	hlog.CtxDebugf(ctx, "处理租户权限变更事件: 租户ID=%s", event.TenantID)

	// 1. 清除租户的权限缓存
	if err := h.tenantCache.InvalidateTenantPermissionCache(ctx, event.TenantID); err != nil {
		return fmt.Errorf("清除租户权限缓存失败: %w", err)
	}

	// 2. 清除租户下所有角色的权限缓存
	roles, err := h.roleCache.GetTenantRoles(ctx, event.TenantID)
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

// 权限相关事件处理
func (h *EventHandler) handlePermissionEvent(ctx context.Context, event *events.PermissionEvent) error {
	hlog.CtxDebugf(ctx, "处理权限事件: %s, 权限ID=%d", event.EventName(), event.PermID)

	switch event.EventName() {
	case events.PermissionCreated:
		// 直接使得权限缓存失效
		if err := h.permCache.InvalidatePermissionCache(ctx, 0); err != nil {
			return fmt.Errorf("清除权限缓存失败: %w", err)
		}
		// 清除权限树缓存
		if err := h.permCache.InvalidatePermissionTreeCache(ctx); err != nil {
			return fmt.Errorf("清除权限树缓存失败: %w", err)
		}

	case events.PermissionUpdated, events.PermissionStatusChange:
		// 1. 清除权限基本信息缓存
		if err := h.permCache.InvalidatePermissionCache(ctx, event.PermID); err != nil {
			return fmt.Errorf("清除权限缓存失败: %w", err)
		}
		// 2. 清除权限树缓存
		if err := h.permCache.InvalidatePermissionTreeCache(ctx); err != nil {
			return fmt.Errorf("清除权限树缓存失败: %w", err)
		}
		// 3. 清除拥有该权限的角色的权限缓存
		roles, err := h.permCache.GetPermissionRoles(ctx, event.PermID)
		if err != nil {
			return fmt.Errorf("获取权限角色列表失败: %w", err)
		}
		for _, role := range roles {
			// 清除角色权限缓存
			if err := h.roleCache.InvalidateRolePermissionCache(ctx, role.ID); err != nil {
				return fmt.Errorf("清除角色[%d]权限缓存失败: %w", role.ID, err)
			}
			// 清除该角色下所有用户的权限和菜单缓存
			users, err := h.roleCache.GetRoleUsers(ctx, role.ID)
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
		}

	case events.PermissionDeleted:
		// 1. 清除权限基本信息缓存
		if err := h.permCache.InvalidatePermissionCache(ctx, event.PermID); err != nil {
			return fmt.Errorf("清除权限缓存失败: %w", err)
		}
		// 2. 清除权限树缓存
		if err := h.permCache.InvalidatePermissionTreeCache(ctx); err != nil {
			return fmt.Errorf("清除权限树缓存失败: %w", err)
		}
		// 3. 清除父权限的子权限列表缓存
		if err := h.permCache.InvalidateChildrenCache(ctx, 0); err != nil {
			return fmt.Errorf("清除根权限子权限列表缓存失败: %w", err)
		}
		// 4. 清除拥有该权限的角色的权限缓存
		roles, err := h.permCache.GetPermissionRoles(ctx, event.PermID)
		if err != nil {
			return fmt.Errorf("获取权限角色列表失败: %w", err)
		}
		for _, role := range roles {
			// 清除角色权限缓存
			if err := h.roleCache.InvalidateRolePermissionCache(ctx, role.ID); err != nil {
				return fmt.Errorf("清除角色[%d]权限缓存失败: %w", role.ID, err)
			}
			// 清除该角色下所有用户的权限和菜单缓存
			users, err := h.roleCache.GetRoleUsers(ctx, role.ID)
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
		}
	}

	return nil
}
