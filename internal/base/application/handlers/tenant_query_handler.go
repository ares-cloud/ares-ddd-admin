package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type TenantQueryHandler struct {
	tenantRepo repository.ITenantRepository
	permRepo   repository.IPermissionsRepository
}

func NewTenantQueryHandler(tenantRepo repository.ITenantRepository, permRepo repository.IPermissionsRepository) *TenantQueryHandler {
	return &TenantQueryHandler{
		tenantRepo: tenantRepo,
		permRepo:   permRepo,
	}
}

func (h *TenantQueryHandler) HandleList(ctx context.Context, q *queries.ListTenantsQuery) (*models.PageRes[dto.TenantDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()

	if q.Code != "" {
		qb.Where("code", query.Like, "%"+q.Code+"%")
	}
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	if q.Status != 0 {
		qb.Where("status", query.Eq, q.Status)
	}

	// 设置分页
	qb.WithPage(&q.Page)

	// 获取总数
	total, err := h.tenantRepo.Count(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	tenants, err := h.tenantRepo.Find(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToTenantDtoList(tenants)

	return &models.PageRes[dto.TenantDto]{
		List:  dtoList,
		Total: total,
	}, nil
}

func (h *TenantQueryHandler) HandleGet(ctx context.Context, query queries.GetTenantQuery) (*dto.TenantDto, herrors.Herr) {
	tenant, err := h.tenantRepo.FindByID(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToTenantDto(tenant), nil
}

func (h *TenantQueryHandler) HandleGetPermissions(ctx context.Context, query queries.GetTenantPermissionsQuery) ([]*dto.PermissionsDto, herrors.Herr) {
	// 查找租户
	tenant, err := h.tenantRepo.FindByID(ctx, query.TenantID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 获取租户权限
	permissions, err := h.tenantRepo.GetPermissions(ctx, tenant.ID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	return dto.ToPermissionsDtoList(permissions), nil
}
