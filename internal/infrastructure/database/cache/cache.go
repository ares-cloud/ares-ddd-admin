package cache

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	// 缓存过期时间
	DefaultExpiration = 24 * time.Hour
	// 缓存更新时间
	DefaultWait = 10 * time.Second
	// 缓存随机过期时间范围
	DefaultRandomExpiration = 60 * time.Second
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (string, error)
	// Set 设置缓存
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	// DeletePrefix 删除前缀的缓存
	DeletePrefix(ctx context.Context, prefix string) error
	// Fetch 获取缓存,不存在则加载
	Fetch(ctx context.Context, key string, expiration time.Duration, fetch func() (string, error)) (string, error)
	// TagAsDeleted 标记缓存已删除
	TagAsDeleted(ctx context.Context, key string) error
}

// NewCache 创建缓存实例
func NewCache(rdb *h_redis.RedisClient, rc *rockscache.Client) Cache {
	return &cacheImpl{
		client: rc,
		rdb:    rdb.GetClient(),
	}
}

type cacheImpl struct {
	client *rockscache.Client
	rdb    *redis.Client
}

func (c *cacheImpl) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *cacheImpl) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *cacheImpl) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *cacheImpl) DeletePrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var result []string
		var err error
		result, cursor, err = c.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, result...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return c.rdb.Del(ctx, keys...).Err()
	}
	return nil
}

func (c *cacheImpl) Fetch(ctx context.Context, key string, expiration time.Duration, fetch func() (string, error)) (string, error) {
	return c.client.Fetch2(ctx, key, expiration, fetch)
}

func (c *cacheImpl) TagAsDeleted(ctx context.Context, key string) error {
	return c.client.TagAsDeleted2(ctx, key)
}
