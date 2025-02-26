package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type PermissionsQueryHandler struct {
	permQuery query.IPermissionsQuery
}

func NewPermissionsQueryHandler(
	permQuery query.IPermissionsQuery,
) *PermissionsQueryHandler {
	return &PermissionsQueryHandler{
		permQuery: permQuery,
	}
}

func (h *PermissionsQueryHandler) HandleList(ctx context.Context, q *queries.ListPermissionsQuery) (*models.PageRes[dto.PermissionsDto], herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()

	// 添加查询条件
	if q.Code != "" {
		qb.Where("code", db_query.Like, "%"+q.Code+"%")
	}
	if q.Name != "" {
		qb.Where("name", db_query.Like, "%"+q.Name+"%")
	}
	if q.Status != 0 {
		qb.Where("status", db_query.Eq, q.Status)
	}
	// 排序
	qb.OrderBy("sequence", true)
	// 设置分页
	qb.WithPage(&q.Page)

	// 查询数据
	perms, total, err := h.permQuery.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return &models.PageRes[dto.PermissionsDto]{
		List:  perms,
		Total: total,
	}, nil
}

func (h *PermissionsQueryHandler) HandleGet(ctx context.Context, query queries.GetPermissionsQuery) (*dto.PermissionsDto, herrors.Herr) {
	perm, err := h.permQuery.FindByID(ctx, query.Id)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

func (h *PermissionsQueryHandler) HandleGetTree(ctx context.Context, query queries.GetPermissionsTreeQuery) ([]*dto.PermissionsDto, herrors.Herr) {
	perms, err := h.permQuery.FindTreeByType(ctx, query.Type)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func (h *PermissionsQueryHandler) HandleGetAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, herrors.Herr) {
	permissions, err := h.permQuery.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (h *PermissionsQueryHandler) HandleGetPermissionsTree(ctx context.Context) (*dto.PermissionsTreeResult, herrors.Herr) {
	return h.permQuery.GetSimplePermissionsTree(ctx)
}
