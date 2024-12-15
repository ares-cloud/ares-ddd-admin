package h_redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 生成唯一ID
func generateUniqueID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// AcquireLockWithUnlock 获取锁
func AcquireLockWithUnlock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (func() (bool, error), error) {
	uniqueID, err := generateUniqueID() // 生成锁的唯一标识
	if err != nil {
		return nil, err
	}

	success, err := rdb.SetNX(ctx, key, uniqueID, ttl).Result() // 尝试获取锁
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, fmt.Errorf("failed to acquire lock")
	}

	// 返回一个释放锁的方法
	unlock := func() (bool, error) {
		luaScript := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
        `
		result, err := rdb.Eval(ctx, luaScript, []string{key}, uniqueID).Int()
		if err != nil {
			return false, err
		}
		return result == 1, nil
	}

	return unlock, nil
}
