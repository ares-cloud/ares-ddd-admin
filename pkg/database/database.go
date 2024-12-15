package database

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/plugin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/snowflake_id"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type contextTxKey struct{}

// Data ， 数据相关类
type Data struct {
	db *gorm.DB
	ig snowflake_id.IIdGenerate
}

// NewData ， 创建 data
// 参数：
//
//	logger ： 日志
//
// 返回值：
//
//	*Data ：desc
//	func() ：desc
//	error ：desc
func NewData(ig snowflake_id.IIdGenerate, cof *configs.Data) (IDataBase, func(), error) {
	var err error
	var db *gorm.DB
	if cof.DataBase.Driver == "pgsql" {
		db, err = gorm.Open(postgres.Open(cof.DataBase.Source), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
			QueryFields:                              true, // 查询使用列
			Logger:                                   logger.Default.LogMode(logger.LogLevel(cof.DataBase.LogLevel)),
		})
	} else {
		db, err = gorm.Open(mysql.Open(cof.DataBase.Source), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
			QueryFields:                              true, // 查询使用列
			Logger:                                   logger.Default.LogMode(logger.LogLevel(cof.DataBase.LogLevel)),
		})
	}
	if err != nil {
		hlog.Fatalf("failed opening connection to mysql: %v", err)
	}
	//迁移基础表
	err = db.Use(plugin.NewTenantPlugin())
	// 获取底层的 SQL 连接池
	sqlDB, err := db.DB()
	if err != nil {
		hlog.Fatalf("failed opening connection to mysql: %v", err)
	}
	cleanup := func() {
		hlog.Info("closing the data resources")
		sqlDB.Close()
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(int(cof.DataBase.MaxIdleConns)) // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(int(cof.DataBase.MaxOpenConns)) // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)                   // 设置连接的最大存活时间
	sqlDB.SetConnMaxIdleTime(20 * time.Minute)            // 设置空闲连接的最大存活时间
	return &Data{db: db, ig: ig}, cleanup, nil
}
func (d Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// 检查 ctx 中是否已有事务
	if _, ok := ctx.Value(contextTxKey{}).(*gorm.DB); ok {
		// 已有事务，直接执行 fn
		return fn(ctx)
	}
	// 未找到事务，创建新的事务
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将 tx 放入到 ctx 中
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d Data) InIndependentTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将 tx 放入到 ctx 中
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d Data) DB(ctx context.Context) *gorm.DB {
	// 从ctx中获取tx
	txKey := ctx.Value(contextTxKey{})
	tx, ok := txKey.(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

func (d Data) GenStringId() string {
	return d.ig.GenStringId()
}
func (d Data) GenInt64Id() int64 {
	//TODO implement me
	panic("implement me")
}
