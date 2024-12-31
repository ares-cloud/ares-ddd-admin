package impl

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type TenantQueryService struct {
	tenantRepo           repository.ISysTenantRepo
	userRepo             repository.ISysUserRepo
	converter            *converter.TenantConverter
	permissionsRepo      repository.IPermissionsRepo
	permissionsConverter *converter.PermissionsConverter
}

func NewTenantQueryService(
	tenantRepo repository.ISysTenantRepo,
	userRepo repository.ISysUserRepo,
	permissionsRepo repository.IPermissionsRepo,
	converter *converter.TenantConverter,
	permissionsConverter *converter.PermissionsConverter,
) *TenantQueryService {
	return &TenantQueryService{
		tenantRepo:           tenantRepo,
		userRepo:             userRepo,
		converter:            converter,
		permissionsRepo:      permissionsRepo,
		permissionsConverter: permissionsConverter,
	}
}

func (t *TenantQueryService) GetTenant(ctx context.Context, id string) (*dto.TenantDto, error) {
	tenant, err := t.tenantRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, nil
	}

	// 获取管理员用户
	adminUser, err := t.userRepo.FindById(ctx, tenant.AdminUserID)
	if err != nil {
		return nil, err
	}

	return t.converter.ToDTO(tenant, adminUser), nil
}

func (t *TenantQueryService) FindTenants(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.TenantDto, error) {
	tenants, err := t.tenantRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}
	return t.converter.ToDTOList(tenants), nil
}

func (t *TenantQueryService) CountTenants(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return t.tenantRepo.Count(ctx, qb)
}

func (t *TenantQueryService) GetTenantPermissions(ctx context.Context, tenantID string) ([]*dto.PermissionsDto, error) {
	// 1. 获取租户的角色
	roles, err := t.tenantRepo.GetTenantRoles(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return []*dto.PermissionsDto{}, nil
	}

	// 2. 获取角色ID列表
	roleIDs := make([]int64, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	// 3. 获取角色对应的权限
	permissions, err := t.permissionsRepo.GetByRoles(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	// 4. 转换为DTO
	return t.permissionsConverter.ToDTOList(permissions), nil
}
