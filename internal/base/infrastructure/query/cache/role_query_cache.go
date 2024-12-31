package cache

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/keys"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	dCache "github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
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
	key := keys.RoleKey(id)
	var role *dto.RoleDto
	err := c.decorator.Cached(ctx, key, &role, func() error {
		var err error
		role, err = c.next.GetRole(ctx, id)
		return err
	})
	return role, err
}

func (c *RoleQueryCache) GetRolePermissions(ctx context.Context, roleID int64) ([]*dto.PermissionsDto, error) {
	key := keys.RolePermissionsKey(roleID)
	var perms []*dto.PermissionsDto
	err := c.decorator.Cached(ctx, key, &perms, func() error {
		var err error
		perms, err = c.next.GetRolePermissions(ctx, roleID)
		return err
	})
	return perms, err
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

// GetRoleByCode 根据编码获取角色
func (c *RoleQueryCache) GetRoleByCode(ctx context.Context, code string) (*dto.RoleDto, error) {
	key := keys.RoleCodeKey(code)
	var role *dto.RoleDto
	err := c.decorator.Cached(ctx, key, &role, func() error {
		var err error
		role, err = c.next.GetRoleByCode(ctx, code)
		return err
	})
	return role, err
}

func (c *RoleQueryCache) FindByID(ctx context.Context, id int64) (*dto.RoleDto, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.RoleKey(tenantID, id)
	// ... 其他代码
}
