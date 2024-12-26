package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/metrics"
	"github.com/redis/go-redis/v9"
)

const (
	deptKeyPrefix   = "dept:"        // 部门缓存前缀
	deptTreeKey     = "dept:tree"    // 部门树缓存key
	deptListKey     = "dept:list"    // 部门列表缓存key
	userDeptKey     = "user:dept:%s" // 用户部门关系缓存key
	deptCacheExpire = 24 * time.Hour // 部门缓存过期时间
	treeCacheExpire = 12 * time.Hour // 树形缓存过期时间
)

type DepartmentCache struct {
	rdb *redis.Client
}

func NewDepartmentCache(rdb *redis.Client) *DepartmentCache {
	return &DepartmentCache{rdb: rdb}
}

// SetDepartment 缓存部门信息
func (c *DepartmentCache) SetDepartment(ctx context.Context, dept *model.Department) error {
	key := c.getDeptKey(dept.ID)
	data, err := json.Marshal(dept)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, deptCacheExpire).Err()
}

// GetDepartment 获取部门缓存
func (c *DepartmentCache) GetDepartment(ctx context.Context, id string) (*model.Department, error) {
	key := c.getDeptKey(id)
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var dept model.Department
	if err := json.Unmarshal(data, &dept); err != nil {
		return nil, err
	}
	return &dept, nil
}

// DeleteDepartment 删除部门缓存
func (c *DepartmentCache) DeleteDepartment(ctx context.Context, id string) error {
	key := c.getDeptKey(id)
	return c.rdb.Del(ctx, key).Err()
}

// SetDepartmentTree 缓存部门树
func (c *DepartmentCache) SetDepartmentTree(ctx context.Context, rootID string, tree []*model.Department) error {
	key := c.getDeptTreeKey(rootID)
	data, err := json.Marshal(tree)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, treeCacheExpire).Err()
}

// GetDepartmentTree 获取部门树缓存
func (c *DepartmentCache) GetDepartmentTree(ctx context.Context, rootID string) ([]*model.Department, error) {
	key := c.getDeptTreeKey(rootID)
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var tree []*model.Department
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}
	return tree, nil
}

// DeleteDepartmentTree 删除部门树缓存
func (c *DepartmentCache) DeleteDepartmentTree(ctx context.Context, rootID string) error {
	key := c.getDeptTreeKey(rootID)
	return c.rdb.Del(ctx, key).Err()
}

// SetUserDepartments 缓存用户部门关系
func (c *DepartmentCache) SetUserDepartments(ctx context.Context, userID string, depts []*model.Department) error {
	key := fmt.Sprintf(userDeptKey, userID)
	data, err := json.Marshal(depts)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, deptCacheExpire).Err()
}

// GetUserDepartments 获取用户部门关系缓存
func (c *DepartmentCache) GetUserDepartments(ctx context.Context, userID string) ([]*model.Department, error) {
	key := fmt.Sprintf(userDeptKey, userID)
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var depts []*model.Department
	if err := json.Unmarshal(data, &depts); err != nil {
		return nil, err
	}
	return depts, nil
}

// DeleteUserDepartments 删除用户部门关系缓存
func (c *DepartmentCache) DeleteUserDepartments(ctx context.Context, userID string) error {
	key := fmt.Sprintf(userDeptKey, userID)
	return c.rdb.Del(ctx, key).Err()
}

// ClearCache 清理所有部门相关缓存
func (c *DepartmentCache) ClearCache(ctx context.Context) error {
	// 1. 清理部��树缓存
	if err := c.rdb.Del(ctx, deptTreeKey).Err(); err != nil {
		return err
	}

	// 2. 清理部门列表缓存
	if err := c.rdb.Del(ctx, deptListKey).Err(); err != nil {
		return err
	}

	// 3. 清理部门缓存
	pattern := deptKeyPrefix + "*"
	keys, err := c.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.rdb.Del(ctx, keys...).Err()
	}

	return nil
}

// 获取部门缓存key
func (c *DepartmentCache) getDeptKey(id string) string {
	return deptKeyPrefix + id
}

// 获取部门树缓存key
func (c *DepartmentCache) getDeptTreeKey(rootID string) string {
	if rootID == "" {
		return deptTreeKey
	}
	return deptTreeKey + ":" + rootID
}

// SetDepartmentWithExpire 设置部门缓存(带过期时间)
func (c *DepartmentCache) SetDepartmentWithExpire(ctx context.Context, dept *model.Department, expire time.Duration) error {
	key := c.getDeptKey(dept.ID)
	data, err := json.Marshal(dept)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, expire).Err()
}

// SetDepartmentTreeWithExpire 设置部门树缓存(带过期时间)
func (c *DepartmentCache) SetDepartmentTreeWithExpire(ctx context.Context, rootID string, tree []*model.Department, expire time.Duration) error {
	key := c.getDeptTreeKey(rootID)
	data, err := json.Marshal(tree)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, expire).Err()
}

// RefreshExpire 刷新缓存过期时间
func (c *DepartmentCache) RefreshExpire(ctx context.Context, id string) error {
	key := c.getDeptKey(id)
	return c.rdb.Expire(ctx, key, deptCacheExpire).Err()
}

// SetNX 设置缓存(不存在时)
func (c *DepartmentCache) SetNX(ctx context.Context, dept *model.Department) error {
	key := c.getDeptKey(dept.ID)
	data, err := json.Marshal(dept)
	if err != nil {
		return err
	}
	return c.rdb.SetNX(ctx, key, data, deptCacheExpire).Err()
}

// GetWithRefresh 获取缓存(自动刷新)
func (c *DepartmentCache) GetWithRefresh(ctx context.Context, id string) (*model.Department, error) {
	dept, err := c.GetDepartment(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept != nil {
		// 异步刷新过期时间
		go func() {
			if err := c.RefreshExpire(context.Background(), id); err != nil {
				metrics.CacheErrorCounter.WithLabelValues("department", "refresh").Inc()
			}
		}()
	}
	return dept, nil
}
