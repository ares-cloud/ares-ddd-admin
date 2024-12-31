package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type TenantQueryHandler struct {
	queryService query.ITenantQueryService
}

func NewTenantQueryHandler(queryService query.ITenantQueryService) *TenantQueryHandler {
	return &TenantQueryHandler{
		queryService: queryService,
	}
}

func (h *TenantQueryHandler) HandleList(ctx context.Context, q *queries.ListTenantsQuery) (*models.PageRes[dto.TenantDto], herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()

	if q.Code != "" {
		qb.Where("code", db_query.Like, "%"+q.Code+"%")
	}
	if q.Name != "" {
		qb.Where("name", db_query.Like, "%"+q.Name+"%")
	}
	if q.Status != 0 {
		qb.Where("status", db_query.Eq, q.Status)
	}

	// 设置分页
	qb.WithPage(&q.Page)

	// 查询数据
	total, err := h.queryService.CountTenants(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	tenants, err := h.queryService.FindTenants(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return &models.PageRes[dto.TenantDto]{
		List:  tenants,
		Total: total,
	}, nil
}

func (h *TenantQueryHandler) HandleGet(ctx context.Context, query queries.GetTenantQuery) (*dto.TenantDto, herrors.Herr) {
	tenant, err := h.queryService.GetTenant(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return tenant, nil
}

func (h *TenantQueryHandler) HandleGetPermissions(ctx context.Context, query queries.GetTenantPermissionsQuery) ([]*dto.PermissionsDto, herrors.Herr) {
	// 查找租户
	_, err := h.queryService.GetTenant(ctx, query.TenantID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 获取租户权限
	permissions, err := h.queryService.GetTenantPermissions(ctx, query.TenantID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	return permissions, nil
}
