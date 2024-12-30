package cache

import (
	"context"
	"encoding/json"
	"fmt"
)

// CacheDecorator 缓存装饰器
type CacheDecorator struct {
	cache Cache
}

// NewCacheDecorator 创建缓存装饰器
func NewCacheDecorator(cache Cache) *CacheDecorator {
	return &CacheDecorator{cache: cache}
}

// Cached 缓存装饰方法
func (d *CacheDecorator) Cached(ctx context.Context, key string, result interface{}, fn func() error) error {
	if key == "" {
		return ErrInvalidKey
	}

	// 尝试从缓存获取
	data, err := d.cache.Get(ctx, key)
	if err == nil {
		// 缓存命中,反序列化
		if err := json.Unmarshal([]byte(data), result); err != nil {
			return fmt.Errorf("unmarshal cache data failed: %w", err)
		}
		return nil
	}

	// 执行原始方法
	if err := fn(); err != nil {
		if err == ErrNotFound {
			// 如果是未找到记录,设置空值缓存,避免缓存穿透
			d.cache.Set(ctx, key, "null", DefaultExpiration)
			return err
		}
		return err
	}

	// 写入缓存
	marshal, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result failed: %w", err)
	}
	return d.cache.Set(ctx, key, string(marshal), DefaultExpiration)
}

// InvalidateCache 使缓存失效
func (d *CacheDecorator) InvalidateCache(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		if err := d.cache.TagAsDeleted(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// InvalidatePrefix 使前缀的缓存失效
func (d *CacheDecorator) InvalidatePrefix(ctx context.Context, prefix string) error {
	return d.cache.DeletePrefix(ctx, prefix)
}

// InvalidateTenantCache 清理租户所有缓存
func (d *CacheDecorator) InvalidateTenantCache(ctx context.Context, tenantID string) error {
	prefix := fmt.Sprintf("tenant:%s", tenantID)
	return d.InvalidatePrefix(ctx, prefix)
}

// InvalidateTenantTypeCache 清理租户特定类型的缓存
func (d *CacheDecorator) InvalidateTenantTypeCache(ctx context.Context, tenantID string, typePrefix string) error {
	prefix := fmt.Sprintf("tenant:%s:%s", tenantID, typePrefix)
	return d.InvalidatePrefix(ctx, prefix)
}
