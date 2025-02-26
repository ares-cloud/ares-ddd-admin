package cache

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type UserQueryCache struct {
	next      *impl.UserQueryService
	decorator *dCache.CacheDecorator
}

func NewUserQueryCache(
	next *impl.UserQueryService,
	decorator *dCache.CacheDecorator,
) *UserQueryCache {
	return &UserQueryCache{
		next:      next,
		decorator: decorator,
	}
}

func (c *UserQueryCache) GetUser(ctx context.Context, id string) (*dto.UserDto, error) {
	key := keys.UserKey(id)
	var user *dto.UserDto
	err := c.decorator.Cached(ctx, key, &user, func() error {
		var err error
		user, err = c.next.GetUser(ctx, id)
		return err
	})
	return user, err
}

func (c *UserQueryCache) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := keys.UserPermissionsKey(userID)
	var permissions []string
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.GetUserPermissions(ctx, userID)
		return err
	})
	return permissions, err
}

func (c *UserQueryCache) GetUserRoles(ctx context.Context, userID string) ([]*dto.RoleDto, error) {
	key := keys.UserRolesKey(userID)
	var roles []*dto.RoleDto
	err := c.decorator.Cached(ctx, key, &roles, func() error {
		var err error
		roles, err = c.next.GetUserRoles(ctx, userID)
		return err
	})
	return roles, err
}

func (c *UserQueryCache) GetUserTreeMenus(ctx context.Context, userID string) ([]*dto.PermissionsTreeDto, error) {
	key := keys.UserMenusKey(userID)
	var menus []*dto.PermissionsTreeDto
	err := c.decorator.Cached(ctx, key, &menus, func() error {
		var err error
		menus, err = c.next.GetUserTreeMenus(ctx, userID)
		return err
	})
	return menus, err
}

// 列表查询不缓存,直接透传
func (c *UserQueryCache) FindUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.FindUsers(ctx, qb)
}

func (c *UserQueryCache) CountUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountUsers(ctx, qb)
}
func (c *UserQueryCache) GetUserMenus(ctx context.Context, userID string) ([]*dto.PermissionsDto, error) {
	return c.next.GetUserMenus(ctx, userID)
}

func (c *UserQueryCache) FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.FindUsersByDepartment(ctx, deptID, excludeAdminID, qb)
}

func (c *UserQueryCache) CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountUsersByDepartment(ctx, deptID, excludeAdminID, qb)
}

// GetUserRolesCode 获取用户角色编码列表(带缓存)
func (c *UserQueryCache) GetUserRolesCode(ctx context.Context, userID string) ([]string, error) {
	key := keys.UserRoleCodesKey(userID)
	var roleCodes []string
	err := c.decorator.Cached(ctx, key, &roleCodes, func() error {
		var err error
		roleCodes, err = c.next.GetUserRolesCode(ctx, userID)
		return err
	})
	return roleCodes, err
}

// InvalidateUserCache 使用户缓存失效
func (c *UserQueryCache) InvalidateUserCache(ctx context.Context, userID string) error {
	keys := []string{
		keys.UserKey(userID),
		keys.UserPermissionsKey(userID),
		keys.UserRolesKey(userID),
		keys.UserMenusKey(userID),
		keys.UserRoleCodesKey(userID),
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}

// InvalidateUserListCache 使用户列表缓存失效
func (c *UserQueryCache) InvalidateUserListCache(ctx context.Context) error {
	return c.decorator.InvalidateCache(ctx, keys.UserListKey())
}

// InvalidateUserPermissionCache 清除用户权限缓存
func (c *UserQueryCache) InvalidateUserPermissionCache(ctx context.Context, userID string) error {
	return c.decorator.InvalidateCache(ctx, keys.UserPermissionsKey(userID))
}

// InvalidateUserMenuCache 清除用户菜单缓存
func (c *UserQueryCache) InvalidateUserMenuCache(ctx context.Context, userID string) error {
	return c.decorator.InvalidateCache(ctx, keys.UserMenusKey(userID))
}

// InvalidateUserDepartmentCache 清除用户部门缓存
func (c *UserQueryCache) InvalidateUserDepartmentCache(ctx context.Context, userID string) error {
	return c.decorator.InvalidateCache(ctx, keys.UserDepartmentKey(userID))
}

// InvalidateRoleListCache 使角色列表缓存
func (c *UserQueryCache) InvalidateRoleListCache(ctx context.Context) error {
	return c.decorator.InvalidateCache(ctx, keys.RoleListKey(actx.GetTenantId(ctx)))
}

// WarmupUserCache 预热用户缓存
func (c *UserQueryCache) WarmupUserCache(ctx context.Context, userID string) error {
	// 1. 预热用户基本信息
	if _, err := c.GetUser(ctx, userID); err != nil {
		return fmt.Errorf("预热用户信息失败: %w", err)
	}

	// 2. 预热用户权限
	if _, err := c.GetUserPermissions(ctx, userID); err != nil {
		return fmt.Errorf("预热用户权限失败: %w", err)
	}

	// 3. 预热用户菜单
	if _, err := c.GetUserTreeMenus(ctx, userID); err != nil {
		return fmt.Errorf("预热用户菜单失败: %w", err)
	}

	return nil
}

// InvalidateTenantUserCache 清除租户下所有用户缓存
func (c *UserQueryCache) InvalidateTenantUserCache(ctx context.Context, tenantID string) error {
	// 使用租户前缀清除所有相关缓存
	return c.decorator.InvalidateTenantTypeCache(ctx, tenantID, "user")
}
