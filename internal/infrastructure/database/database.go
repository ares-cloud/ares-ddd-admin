package database

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/snowflake_id"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	snowflake_id.NewSnowIdGen,
	NewHdbClient,
	NewDataBase,
)

func NewDataBase(ig snowflake_id.IIdGenerate, conf *configs.Data) (database.IDataBase, func(), error) {
	return database.NewData(ig, conf)
}

func NewHdbClient(conf *configs.Data) (*h_redis.RedisClient, error) {
	return h_redis.NewRedisClient(h_redis.Option{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.Db),
		Timeout:  conf.Redis.WriteTimeout,
	})
}
