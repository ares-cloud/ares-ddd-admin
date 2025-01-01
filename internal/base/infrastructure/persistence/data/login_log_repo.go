package data

import (
	"context"
	"fmt"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type loginLogRepo struct {
	db database.IDataBase
}

func NewLoginLogRepo(db database.IDataBase) repository.ILoginLogRepo {
	return &loginLogRepo{
		db: db,
	}
}

func (r *loginLogRepo) Create(ctx context.Context, log *entity.LoginLog) error {
	t := time.Unix(log.LoginTime, 0)
	if err := r.EnsureTable(ctx, log.TenantID, t); err != nil {
		return err
	}
	return r.db.DB(ctx).Table(r.GetTableName(log.TenantID, t)).Create(log).Error
}

func (r *loginLogRepo) FindByID(ctx context.Context, id int64) (*entity.LoginLog, error) {
	var entity entity.LoginLog
	err := r.db.DB(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *loginLogRepo) Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*entity.LoginLog, error) {
	// 确保表存在
	if err := r.EnsureTable(ctx, tenantID, month); err != nil {
		return nil, err
	}

	var entities []*entity.LoginLog
	db := r.db.DB(ctx).Table(r.GetTableName(tenantID, month))

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	// 添加排序
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	}

	// 添加分页
	if limit, offset := qb.BuildLimit(); limit != "" {
		db = db.Limit(offset[1]).Offset(offset[0])
	}

	// 执行查询
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *loginLogRepo) Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error) {
	// 确保表存在
	if err := r.EnsureTable(ctx, tenantID, month); err != nil {
		return 0, err
	}

	var count int64
	db := r.db.DB(ctx).Table(r.GetTableName(tenantID, month))

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *loginLogRepo) EnsureTable(ctx context.Context, tenantID string, month time.Time) error {
	tableName := r.GetTableName(tenantID, month)

	// 检查表是否存在
	if r.db.DB(ctx).Migrator().HasTable(tableName) {
		return nil
	}

	// 创建一个临时结构体
	type LoginLogTable struct {
		entity.LoginLog
	}

	// 使用 GORM 自动迁移创建表
	return r.db.DB(ctx).Table(tableName).AutoMigrate(&LoginLogTable{})
}

func (r *loginLogRepo) GetTableName(tenantID string, month time.Time) string {
	return fmt.Sprintf("sys_login_log_%s_%s", tenantID, month.Format("200601"))
}
