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

type DepartmentQueryCache struct {
	next      *impl.DepartmentQueryService
	decorator *dCache.CacheDecorator
}

func NewDepartmentQueryCache(
	next *impl.DepartmentQueryService,
	decorator *dCache.CacheDecorator,
) *DepartmentQueryCache {
	return &DepartmentQueryCache{
		next:      next,
		decorator: decorator,
	}
}

// GetDepartment 获取部门详情(带缓存)
func (c *DepartmentQueryCache) GetDepartment(ctx context.Context, id string) (*dto.DepartmentDto, error) {
	key := keys.DepartmentKey(id)
	var dept *dto.DepartmentDto
	err := c.decorator.Cached(ctx, key, &dept, func() error {
		var err error
		dept, err = c.next.GetDepartment(ctx, id)
		return err
	})
	return dept, err
}

// GetDepartmentTree 获取部门树(带缓存)
func (c *DepartmentQueryCache) GetDepartmentTree(ctx context.Context, parentID string) ([]*dto.DepartmentTreeDto, error) {
	key := keys.DepartmentTreeKey(parentID)
	var tree []*dto.DepartmentTreeDto
	err := c.decorator.Cached(ctx, key, &tree, func() error {
		var err error
		tree, err = c.next.GetDepartmentTree(ctx, parentID)
		return err
	})
	return tree, err
}

// GetUserDepartments 获取用户部门(带缓存)
func (c *DepartmentQueryCache) GetUserDepartments(ctx context.Context, userID string) ([]*dto.DepartmentDto, error) {
	key := keys.UserDepartmentsKey(userID)
	var depts []*dto.DepartmentDto
	err := c.decorator.Cached(ctx, key, &depts, func() error {
		var err error
		depts, err = c.next.GetUserDepartments(ctx, userID)
		return err
	})
	return depts, err
}

// 列表查询不缓存,直接透传
func (c *DepartmentQueryCache) FindDepartments(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.DepartmentDto, error) {
	return c.next.FindDepartments(ctx, qb)
}

func (c *DepartmentQueryCache) CountDepartments(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountDepartments(ctx, qb)
}

func (c *DepartmentQueryCache) GetDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.GetDepartmentUsers(ctx, deptID, excludeAdminID, qb)
}

func (c *DepartmentQueryCache) CountDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountDepartmentUsers(ctx, deptID, excludeAdminID, qb)
}

func (c *DepartmentQueryCache) GetUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.GetUnassignedUsers(ctx, qb)
}

func (c *DepartmentQueryCache) CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountUnassignedUsers(ctx, qb)
}

// InvalidateDepartmentCache 使部门缓存失效
func (c *DepartmentQueryCache) InvalidateDepartmentCache(ctx context.Context, deptID string) error {
	keys := []string{
		keys.DepartmentKey(deptID),
		keys.DepartmentTreeKey(""),
	}
	return c.decorator.InvalidateCache(ctx, keys...)
}

// InvalidateUserDepartmentsCache 使用户部门缓存失效
func (c *DepartmentQueryCache) InvalidateUserDepartmentsCache(ctx context.Context, userID string) error {
	return c.decorator.InvalidateCache(ctx, keys.UserDepartmentsKey(userID))
}

func (c *DepartmentQueryCache) FindByID(ctx context.Context, id string) (*dto.DepartmentDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.DepartmentKey(tenantID, id)
	var dept *dto.DepartmentDto
	err := c.decorator.Cached(ctx, key, &dept, func() error {
		var err error
		dept, err = c.next.FindByID(ctx, id)
		return err
	})
	return dept, err
}
