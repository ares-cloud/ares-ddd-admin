package database

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database/cache"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/snowflake_id"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/dtm-labs/rockscache"
	"github.com/google/wire"
	"time"
)

var ProviderSet = wire.NewSet(
	snowflake_id.NewSnowIdGen,
	NewHdbClient,
	NewDataBase,
	NewRc,
	cache.NewCache,
	cache.NewCacheDecorator,
)

func NewDataBase(ig snowflake_id.IIdGenerate, conf *configs.Data) (database.IDataBase, func(), error) {
	return database.NewData(ig, conf)
}

func NewHdbClient(conf *configs.Data) (*h_redis.RedisClient, func(), error) {
	return h_redis.NewRedisClient(h_redis.Option{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.Db),
		Timeout:  conf.Redis.WriteTimeout,
	})
}

func NewRc(rdb *h_redis.RedisClient) *rockscache.Client {
	// 强一致性缓存，当一个key被标记删除，其他请求线程会被锁住轮询直到新的key生成，适合各种同步的拉取, 如果弱一致可能导致拉取还是老数据，毫无意义
	options := rockscache.NewDefaultOptions()
	options.StrongConsistency = true
	options.Delay = time.Millisecond * 1
	Rc := rockscache.NewClient(rdb.GetClient(), options)
	Rc.Options.StrongConsistency = true
	return Rc
}
