package cache

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type RoleQueryCache struct {
	next      *impl.RoleQueryService
	decorator *dCache.CacheDecorator
}

func NewRoleQueryCache(
	next *impl.RoleQueryService,
	decorator *dCache.CacheDecorator,
) *RoleQueryCache {
	return &RoleQueryCache{
		next:      next,
		decorator: decorator,
	}
}

func (c *RoleQueryCache) GetRole(ctx context.Context, id int64) (*dto.RoleDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.RoleKey(tenantID, id)
	var role *dto.RoleDto
	err := c.decorator.Cached(ctx, key, &role, func() error {
		var err error
		role, err = c.next.GetRole(ctx, id)
		return err
	})
	return role, err
}

func (c *RoleQueryCache) GetRolePermissions(ctx context.Context, roleID int64) ([]*dto.PermissionsDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.RolePermissionsKey(tenantID, roleID)
	var permissions []*dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.GetRolePermissions(ctx, roleID)
		return err
	})
	return permissions, err
}

// 列表查询不缓存,直接透传
func (c *RoleQueryCache) FindRoles(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.RoleDto, error) {
	return c.next.FindRoles(ctx, qb)
}

func (c *RoleQueryCache) CountRoles(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountRoles(ctx, qb)
}

func (c *RoleQueryCache) FindByType(ctx context.Context, roleType int8) ([]*dto.RoleDto, error) {
	return c.next.FindByType(ctx, roleType)
}

// GetRoleByCode 根据编码获取角色(带缓存)
func (c *RoleQueryCache) GetRoleByCode(ctx context.Context, code string) (*dto.RoleDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.RoleCodeKey(tenantID, code)
	var role *dto.RoleDto
	err := c.decorator.Cached(ctx, key, &role, func() error {
		var err error
		role, err = c.next.GetRoleByCode(ctx, code)
		return err
	})
	return role, err
}

// InvalidateRolePermissionCache 清除角色权限缓存
func (c *RoleQueryCache) InvalidateRolePermissionCache(ctx context.Context, roleID int64) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.RolePermissionsKey(tenantID, roleID))
}

// GetRoleUsers 获取角色下的用户列表
func (c *RoleQueryCache) GetRoleUsers(ctx context.Context, roleID int64) ([]*dto.UserDto, error) {
	return c.next.GetRoleUsers(ctx, roleID)
}

// GetTenantRoles 获取租户下的角色列表
func (c *RoleQueryCache) GetTenantRoles(ctx context.Context, tenantID string) ([]*dto.RoleDto, error) {
	return c.next.GetTenantRoles(ctx, tenantID)
}

// InvalidateRoleCache 清除角色缓存
func (c *RoleQueryCache) InvalidateRoleCache(ctx context.Context, roleID int64) error {
	tenantID := actx.GetTenantId(ctx)
	keys := []string{
		keys.RoleKey(tenantID, roleID),
		keys.RolePermissionsKey(tenantID, roleID),
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}

// InvalidateRoleListCache 清除角色列表缓存
func (c *RoleQueryCache) InvalidateRoleListCache(ctx context.Context) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.RoleListKey(tenantID))
}

// InvalidateTenantRoleCache 清除租户下所有角色缓存
func (c *RoleQueryCache) InvalidateTenantRoleCache(ctx context.Context, tenantID string) error {
	// 使用租户前缀清除所有相关缓存
	return c.decorator.InvalidateTenantTypeCache(ctx, tenantID, "role")
}

// InvalidateRoleUserCache 清除角色用户列表缓存
func (c *RoleQueryCache) InvalidateRoleUserCache(ctx context.Context, roleID int64) error {
	key := keys.RoleUsersKey(actx.GetTenantId(ctx), roleID)
	return c.decorator.InvalidateCache(ctx, key)
}
