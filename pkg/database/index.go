package database

import (
	"context"
	"gorm.io/gorm"
)

type IDataBase interface {
	// InTx 下面2个方法配合使用，在InTx方法中执行ORM操作的时候需要使用DB方法获取db！
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	// InIndependentTx 开启独立事物
	InIndependentTx(ctx context.Context, fn func(ctx context.Context) error) error
	// DB 获取链接
	DB(ctx context.Context) *gorm.DB
	// GenStringId 主键生成
	GenStringId() string
	// GenInt64Id 生成int类型的id
	GenInt64Id() int64
}
