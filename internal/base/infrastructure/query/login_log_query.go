package query

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// ILoginLogQuery 登录日志查询接口
type ILoginLogQuery interface {
	// Find 查询登录日志列表
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*dto.LoginLogDto, error)
	// Count 统计登录日志数量
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}
