package cache

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type PermissionsQueryCache struct {
	next      *impl.PermissionsQueryService
	decorator *dCache.CacheDecorator
}

func NewPermissionsQueryCache(
	next *impl.PermissionsQueryService,
	decorator *dCache.CacheDecorator,
) *PermissionsQueryCache {
	return &PermissionsQueryCache{
		next:      next,
		decorator: decorator,
	}
}

// FindByID 根据ID查询权限(带缓存)
func (c *PermissionsQueryCache) FindByID(ctx context.Context, id int64) (*dto.PermissionsDto, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.PermissionKey(tenantID, id)
	var permission *dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &permission, func() error {
		var err error
		permission, err = c.next.FindByID(ctx, id)
		return err
	})
	if err != nil {
		if herrors.IsHError(err) {
			return nil, herrors.TohError(err)
		}
		return nil, herrors.NewServerHError(err)
	}
	return permission, nil
}

// FindTreeByType 查询权限树(带缓存)
func (c *PermissionsQueryCache) FindTreeByType(ctx context.Context, permType int8) ([]*dto.PermissionsDto, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.PermissionTreeKey(tenantID, permType)
	var permissions []*dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.FindTreeByType(ctx, permType)
		return err
	})
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return permissions, nil
}

// FindAllEnabled 查询所有启用的权限(带缓存)
func (c *PermissionsQueryCache) FindAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.PermissionEnabledKey(tenantID)
	var permissions []*dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &permissions, func() error {
		var err error
		permissions, err = c.next.FindAllEnabled(ctx)
		return err
	})
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return permissions, nil
}

// GetSimplePermissionsTree 获取简化的权限树(带缓存)
func (c *PermissionsQueryCache) GetSimplePermissionsTree(ctx context.Context) (*dto.PermissionsTreeResult, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.PermissionSimpleTreeKey(tenantID)
	var permissionsTree *dto.PermissionsTreeResult
	err := c.decorator.Cached(ctx, key, &permissionsTree, func() error {
		var err error
		permissionsTree, err = c.next.GetSimplePermissionsTree(ctx)
		return err
	})
	if err != nil {
		if herrors.IsHError(err) {
			return nil, herrors.TohError(err)
		}
		return nil, herrors.NewServerHError(err)
	}
	return permissionsTree, nil
}

// Find 查询权限列表(不缓存)
func (c *PermissionsQueryCache) Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.PermissionsDto, int64, herrors.Herr) {
	return c.next.Find(ctx, qb)
}
func (c *PermissionsQueryCache) GetPermissionRoles(ctx context.Context, permID int64) ([]*dto.RoleDto, error) {
	return c.next.GetPermissionRoles(ctx, permID)
}

// InvalidatePermissionCache 使权限缓存失效
func (c *PermissionsQueryCache) InvalidatePermissionCache(ctx context.Context, id int64) error {
	tenantID := actx.GetTenantId(ctx)
	keys := []string{
		keys.PermissionKey(tenantID, id),
		keys.PermissionTreeKey(tenantID, "*"),
		keys.PermissionListKey(tenantID),
		keys.PermissionEnabledKey(tenantID),
		keys.PermissionSimpleTreeKey(tenantID),
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}

// WarmupPermissionCache 预热权限缓存
func (c *PermissionsQueryCache) WarmupPermissionCache(ctx context.Context, id int64) error {
	// 1. 预热权限详情
	if _, err := c.FindByID(ctx, id); err != nil {
		return fmt.Errorf("预热权限详情缓存失败: %w", err)
	}

	// 2. 预热权限树
	if _, err := c.FindTreeByType(ctx, 1); err != nil { // 预热菜单类型权限树
		return fmt.Errorf("预热权限树缓存失败: %w", err)
	}

	return nil
}

// InvalidateChildrenCache 清除权限子节点列表缓存
func (c *PermissionsQueryCache) InvalidateChildrenCache(ctx context.Context, parentID int64) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.PermissionChildrenKey(tenantID, parentID))
}

// InvalidatePermissionTreeCache 清除权限树缓存
func (c *PermissionsQueryCache) InvalidatePermissionTreeCache(ctx context.Context) error {
	tenantID := actx.GetTenantId(ctx)
	keys := []string{
		keys.PermissionTreeKey(tenantID, "*"),  // 所有类型的权限树
		keys.PermissionSimpleTreeKey(tenantID), // 简化的权限树
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}

// InvalidatePermissionResourceCache 清除权限资源缓存
func (c *PermissionsQueryCache) InvalidatePermissionResourceCache(ctx context.Context, permID int64) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.PermissionResourceKey(tenantID, permID))
}
