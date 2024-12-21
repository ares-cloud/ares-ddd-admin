package repository

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/domain/model"
)

type IMetricsRepository interface {
	// 系统指标
	SaveSystemMetrics(ctx context.Context, metrics *model.SystemMetrics) error
	GetSystemMetrics(ctx context.Context, tenantID string, start, end time.Time) ([]*model.SystemMetrics, error)

	// 数据库指标
	SaveDatabaseMetrics(ctx context.Context, metrics *model.DatabaseMetrics) error
	GetDatabaseMetrics(ctx context.Context, tenantID string, start, end time.Time) ([]*model.DatabaseMetrics, error)

	// Redis指标
	SaveRedisMetrics(ctx context.Context, metrics *model.RedisMetrics) error
	GetRedisMetrics(ctx context.Context, tenantID string, start, end time.Time) ([]*model.RedisMetrics, error)

	// API指标
	SaveAPIMetrics(ctx context.Context, metrics *model.APIMetrics) error
	GetAPIMetrics(ctx context.Context, tenantID string, start, end time.Time) ([]*model.APIMetrics, error)
}
