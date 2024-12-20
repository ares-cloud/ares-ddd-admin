package repository

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type IOperationLogRepository interface {
	Create(ctx context.Context, log *model.OperationLog) error
	FindByID(ctx context.Context, id int64) (*model.OperationLog, error)
	Find(ctx context.Context, tenantID string, month time.Time, qb *query.QueryBuilder) ([]*model.OperationLog, error)
	Count(ctx context.Context, tenantID string, month time.Time, qb *query.QueryBuilder) (int64, error)
	EnsureTable(ctx context.Context, tenantID string, month time.Time) error
	GetTableName(tenantID string, month time.Time) string
}
