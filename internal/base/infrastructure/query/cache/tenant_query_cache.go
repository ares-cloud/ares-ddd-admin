package cache

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type TenantQueryCache struct {
	next      *impl.TenantQueryService
	decorator *dCache.CacheDecorator
}

func NewTenantQueryCache(
	next *impl.TenantQueryService,
	decorator *dCache.CacheDecorator,
) *TenantQueryCache {
	return &TenantQueryCache{
		next:      next,
		decorator: decorator,
	}
}

func (c *TenantQueryCache) GetTenant(ctx context.Context, id string) (*dto.TenantDto, error) {
	key := keys.TenantKey(id)
	var tenant *dto.TenantDto
	err := c.decorator.Cached(ctx, key, &tenant, func() error {
		var err error
		tenant, err = c.next.GetTenant(ctx, id)
		return err
	})
	return tenant, err
}

// 列表查询不缓存,直接透传
func (c *TenantQueryCache) FindTenants(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.TenantDto, error) {
	return c.next.FindTenants(ctx, qb)
}

func (c *TenantQueryCache) CountTenants(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountTenants(ctx, qb)
}

// GetTenantPermissions 获取租户权限(带缓存)
func (c *TenantQueryCache) GetTenantPermissions(ctx context.Context, tenantID string) ([]*dto.PermissionsDto, error) {
	key := keys.TenantPermissionsKey(tenantID)
	var permissions []*dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.GetTenantPermissions(ctx, tenantID)
		return err
	})
	return permissions, err
}

// InvalidateTenantPermissionCache 清除租户权限缓存
func (c *TenantQueryCache) InvalidateTenantPermissionCache(ctx context.Context, tenantID string) error {
	return c.decorator.InvalidateCache(ctx, keys.TenantPermissionsKey(tenantID))
}

// InvalidateTenantStatusCache 清除租户状态缓存
func (c *TenantQueryCache) InvalidateTenantStatusCache(ctx context.Context, tenantID string) error {
	return c.decorator.InvalidateCache(ctx, keys.TenantStatusKey(tenantID))
}

// InvalidateTenantCache 清除租户缓存
func (c *TenantQueryCache) InvalidateTenantCache(ctx context.Context, tenantID string) error {
	keys := []string{
		keys.TenantKey(tenantID),
		keys.TenantStatusKey(tenantID),
		keys.TenantPermissionsKey(tenantID),
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}
