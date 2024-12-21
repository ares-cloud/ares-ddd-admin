package repository

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ILoginLogRepository interface {
	Create(ctx context.Context, log *model.LoginLog) error
	FindByID(ctx context.Context, id int64) (*model.LoginLog, error)

	// 动态查询方法
	Find(ctx context.Context, tenantID string, month time.Time, qb *query.QueryBuilder) ([]*model.LoginLog, error)
	Count(ctx context.Context, tenantID string, month time.Time, qb *query.QueryBuilder) (int64, error)

	// 表管理方法
	EnsureTable(ctx context.Context, tenantID string, month time.Time) error
	GetTableName(tenantID string, month time.Time) string
}
