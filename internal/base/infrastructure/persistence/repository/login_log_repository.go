package repository

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type loginLogRepository struct {
	db     database.IDataBase
	mapper *mapper.LoginLogMapper
}

func NewLoginLogRepository(db database.IDataBase) repository.ILoginLogRepository {
	return &loginLogRepository{
		db:     db,
		mapper: &mapper.LoginLogMapper{},
	}
}

func (r *loginLogRepository) Create(ctx context.Context, log *model.LoginLog) error {
	// 确保表存在
	t := time.Unix(log.LoginTime, 0)
	if err := r.EnsureTable(ctx, log.TenantID, t); err != nil {
		return err
	}

	// 转换为实体并保存
	entity := r.mapper.ToEntity(log)
	return r.db.DB(ctx).Table(r.GetTableName(log.TenantID, t)).Create(entity).Error
}

func (r *loginLogRepository) FindByID(ctx context.Context, id int64) (*model.LoginLog, error) {
	var entity entity.LoginLog
	err := r.db.DB(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(&entity), nil
}

func (r *loginLogRepository) Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*model.LoginLog, error) {
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

	return r.mapper.ToDomainList(entities), nil
}

func (r *loginLogRepository) Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error) {
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

func (r *loginLogRepository) EnsureTable(ctx context.Context, tenantID string, month time.Time) error {
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

func (r *loginLogRepository) GetTableName(tenantID string, month time.Time) string {
	return fmt.Sprintf("sys_login_log_%s_%s", tenantID, month.Format("200601"))
}
