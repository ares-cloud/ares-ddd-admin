package query

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// ITenantQueryService 租户查询服务接口
type ITenantQueryService interface {
	GetTenant(ctx context.Context, id string) (*dto.TenantDto, error)
	FindTenants(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.TenantDto, error)
	CountTenants(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
	GetTenantPermissions(ctx context.Context, tenantID string) ([]*dto.PermissionsDto, error)
}
