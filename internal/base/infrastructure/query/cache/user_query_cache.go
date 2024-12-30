package cache

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	dquery "github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// UserQueryCache 用户查询缓存装饰器
type UserQueryCache struct {
	next      *impl.UserQueryService
	decorator *dCache.CacheDecorator
}

// NewUserQueryCache 创建用户查询缓存装饰器
func NewUserQueryCache(
	next *impl.UserQueryService,
	decorator *dCache.CacheDecorator,
) *UserQueryCache {
	return &UserQueryCache{
		next:      next,
		decorator: decorator,
	}
}

// GetUser 获取用户信息(带缓存)
func (c *UserQueryCache) GetUser(ctx context.Context, id string) (*model.User, error) {
	key := keys.UserKey(id)
	var user *model.User
	err := c.decorator.Cached(ctx, key, &user, func() error {
		var err error
		user, err = c.next.GetUser(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})

	return user, err
}

// GetUserPermissions 获取用户权限(带缓存)
func (c *UserQueryCache) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := keys.UserPermKey(userID)
	var permissions []string
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.GetUserPermissions(ctx, userID)
		return err
	})
	return permissions, err
}

// GetUserRoles 获取用户角色(带缓存)
func (c *UserQueryCache) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	key := keys.UserRoleKey(userID)
	var roles []*model.Role

	// 先清除可能存在的错误类型缓存
	if err := c.decorator.InvalidateCache(ctx, key); err != nil {
		return nil, fmt.Errorf("清除旧缓存失败: %w", err)
	}

	err := c.decorator.Cached(ctx, key, &roles, func() error {
		var err error
		roles, err = c.next.GetUserRoles(ctx, userID)
		if err != nil {
			return err
		}
		return nil
	})

	return roles, err
}
func (c *UserQueryCache) GetUserRolesCode(ctx context.Context, userID string) ([]string, error) {
	roles, err := c.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(roles))
	for _, role := range roles {
		codes = append(codes, role.Code)
	}
	return codes, nil
}

// GetUserMenus 获取用户菜单(带缓存)
func (c *UserQueryCache) GetUserMenus(ctx context.Context, userID string) ([]*model.Permissions, error) {
	key := keys.UserMenuKey(userID)
	var menus []*model.Permissions
	err := c.decorator.Cached(ctx, key, &menus, func() error {
		var err error
		menus, err = c.next.GetUserMenus(ctx, userID)
		return err
	})
	return menus, err
}

// GetUserTreeMenus 获取用户菜单(带缓存)
func (c *UserQueryCache) GetUserTreeMenus(ctx context.Context, userID string) ([]*model.Permissions, error) {
	key := keys.UserMenuTreeKey(userID)
	var menus []*model.Permissions
	// 2. 使用新的缓存
	err := c.decorator.Cached(ctx, key, &menus, func() error {
		var err error
		menus, err = c.next.GetUserTreeMenus(ctx, userID)
		return err
	})

	return menus, err
}

// 列表查询不缓存,直接透传
func (c *UserQueryCache) FindUsers(ctx context.Context, qb *dquery.QueryBuilder) ([]*model.User, error) {
	return c.next.FindUsers(ctx, qb)
}

func (c *UserQueryCache) CountUsers(ctx context.Context, qb *dquery.QueryBuilder) (int64, error) {
	return c.next.CountUsers(ctx, qb)
}

func (c *UserQueryCache) FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *dquery.QueryBuilder) ([]*model.User, error) {
	return c.next.FindUsersByDepartment(ctx, deptID, excludeAdminID, qb)
}

func (c *UserQueryCache) CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *dquery.QueryBuilder) (int64, error) {
	return c.next.CountUsersByDepartment(ctx, deptID, excludeAdminID, qb)
}

func (c *UserQueryCache) FindUnassignedUsers(ctx context.Context, qb *dquery.QueryBuilder) ([]*model.User, error) {
	return c.next.FindUnassignedUsers(ctx, qb)
}

func (c *UserQueryCache) CountUnassignedUsers(ctx context.Context, qb *dquery.QueryBuilder) (int64, error) {
	return c.next.CountUnassignedUsers(ctx, qb)
}

// InvalidateUserCache 使用户基本信息缓存失效
func (c *UserQueryCache) InvalidateUserCache(ctx context.Context, userID string) error {
	// 清除用户基本信息缓存
	if err := c.decorator.InvalidateCache(ctx, keys.UserKey(userID)); err != nil {
		return err
	}

	// 同时清除用户相关的所有缓存
	return c.InvalidateUserPermissionCache(ctx, userID)
}

// InvalidateUserListCache 使用户列表缓存失效
func (c *UserQueryCache) InvalidateUserListCache(ctx context.Context) error {
	return c.decorator.InvalidateCache(ctx, keys.UserListKey())
}

// InvalidateUserPermissionCache 使用户权限相关缓存失效
func (c *UserQueryCache) InvalidateUserPermissionCache(ctx context.Context, userID string) error {
	// 清除用户权限相关的所有缓存
	keysToDelete := []string{
		keys.UserPermKey(userID),     // 权限列表缓存
		keys.UserRoleKey(userID),     // 角色列表缓存
		keys.UserMenuKey(userID),     // 菜单列表缓存
		keys.UserMenuTreeKey(userID), // 菜单树缓存
	}
	return c.decorator.InvalidateCache(ctx, keysToDelete...)
}

// InvalidateRoleListCache 使角色列表缓存失效
func (c *UserQueryCache) InvalidateRoleListCache(ctx context.Context) error {
	return c.decorator.InvalidateCache(ctx, keys.RoleListKey())
}

// WarmupUserCache 预热用户缓存
func (c *UserQueryCache) WarmupUserCache(ctx context.Context, userID string) error {
	// 1. 预热用户基本信息
	if _, err := c.GetUser(ctx, userID); err != nil {
		return fmt.Errorf("预热用户信息缓存失败: %w", err)
	}

	// 2. 预热用户权限
	if _, err := c.GetUserPermissions(ctx, userID); err != nil {
		return fmt.Errorf("预热用户权限缓存失败: %w", err)
	}

	// 3. 预热用户角色
	if _, err := c.GetUserRoles(ctx, userID); err != nil {
		return fmt.Errorf("预热用户角色缓存失败: %w", err)
	}

	// 4. 预热用户菜单树
	if _, err := c.GetUserTreeMenus(ctx, userID); err != nil {
		return fmt.Errorf("预热用户菜单树缓存失败: %w", err)
	}

	return nil
}
