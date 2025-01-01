package repository

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type ILoginLogRepo interface {
	Create(ctx context.Context, log *entity.LoginLog) error
	// 动态查询方法
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*entity.LoginLog, error)
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}

type loginLogRepository struct {
	db     ILoginLogRepo
	mapper *mapper.LoginLogMapper
}

func NewLoginLogRepository(db ILoginLogRepo) repository.ILoginLogRepository {
	return &loginLogRepository{
		db:     db,
		mapper: &mapper.LoginLogMapper{},
	}
}

func (r *loginLogRepository) Create(ctx context.Context, log *model.LoginLog) error {
	toEntity := r.mapper.ToEntity(log)
	return r.db.Create(ctx, toEntity)
}
