package query

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IOperationLogQuery 操作日志查询接口
type IOperationLogQuery interface {
	// Find 查询操作日志列表
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*dto.OperationLogDto, error)
	// Count 统计操作日志数量
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}
