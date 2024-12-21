package h_redis

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"log"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	rs     *redsync.Redsync
}

func NewRedisClient(opt Option) (*RedisClient, func(), error) {
	db := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	err := db.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	pool := goredis.NewPool(db)
	cleanup := func() {
		hlog.Info("closing the data resources")
		db.Close()
	}
	rs := redsync.New(pool)
	return &RedisClient{
		client: db,
		rs:     rs,
	}, cleanup, nil
}

// MutexWithUnlock 分布式锁，并发控制
func (rc *RedisClient) MutexWithUnlock(name string, options ...redsync.Option) (UnlockFunc, error) {
	mutex := rc.rs.NewMutex(name, options...)
	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	unlock := func() error {
		_, err := mutex.Unlock()
		if err != nil {
			return err
		}
		return nil
	}

	return unlock, nil
}

// SimpleMutexWithUnlock 分布式锁，并发控制
func (rc *RedisClient) SimpleMutexWithUnlock(name string) (UnlockFunc, error) {
	mutex := rc.rs.NewMutex(name)
	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	unlock := func() error {
		_, err := mutex.Unlock()
		if err != nil {
			return err
		}
		return nil
	}

	return unlock, nil
}

// GetClient 获取客户端
func (rc *RedisClient) GetClient() *redis.Client {
	return rc.client
}

// IfErrorNotNil 是否为非空错误
func IfErrorNotNil(err error) bool {
	return err != nil && !errors.Is(err, redis.Nil)
}

// IncrementCounter 自增计数器
func IncrementCounter(ctx context.Context, rdb *redis.Client, key string) (int64, error) {
	// 将 key 对应的值增加 1，初始值为1
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCounterValue 获取当前计数器值
func GetCounterValue(ctx context.Context, rdb *redis.Client, key string) (int64, error) {
	// 获取当前计数器值
	value, err := rdb.Get(ctx, key).Int64()
	if err != nil {
		// 如果 key 不存在，返回0
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return value, nil
}

// ResetCounter 更新计数器到安全值（当发生错误时）
func ResetCounter(ctx context.Context, rdb *redis.Client, key string, safeValue int64) error {
	// 设置计数器到一个安全值
	err := rdb.Set(ctx, key, safeValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
