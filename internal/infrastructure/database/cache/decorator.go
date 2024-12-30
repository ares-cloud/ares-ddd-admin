package cache

import (
	"context"
	"encoding/json"
	"errors"
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
	if result == nil {
		return ErrInvalidResult
	}

	fetch, err := d.cache.Fetch(ctx, key, DefaultExpiration, func() (string, error) {
		// 执行原始方法
		if err := fn(); err != nil {
			if errors.Is(err, ErrNotFound) {
				return "", nil
			}
			return "", err
		}
		marshal, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("marshal result failed: %w", err)
		}
		return string(marshal), nil
	})
	if err != nil {
		return err
	}
	if fetch != "" {
		if err = json.Unmarshal([]byte(fetch), result); err != nil {
			// 如果反序列化失败，可能是类型不匹配，清除缓存
			d.InvalidateCache(ctx, key)
			return fmt.Errorf("unmarshal result failed: %w", err)
		}
	}
	return nil
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
