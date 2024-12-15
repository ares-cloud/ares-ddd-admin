package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type PermissionsQueryHandler struct {
	permRepo repository.IPermissionsRepository
}

func NewPermissionsQueryHandler(permRepo repository.IPermissionsRepository) *PermissionsQueryHandler {
	return &PermissionsQueryHandler{
		permRepo: permRepo,
	}
}

// HandleList 处理列表查询
func (h *PermissionsQueryHandler) HandleList(ctx context.Context, q *queries.ListPermissionsQuery) (*models.PageRes[dto.PermissionsDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()

	// 添加查询条件
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
	total, err := h.permRepo.Count(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	perms, err := h.permRepo.Find(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToPermissionsDtoList(perms)

	return &models.PageRes[dto.PermissionsDto]{
		List:  dtoList,
		Total: total,
	}, nil
}

// HandleGet 保留其他原有方法
func (h *PermissionsQueryHandler) HandleGet(ctx context.Context, query queries.GetPermissionsQuery) (*dto.PermissionsDto, herrors.Herr) {
	perm, err := h.permRepo.FindByID(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToPermissionsDto(perm), nil
}

func (h *PermissionsQueryHandler) HandleGetTree(ctx context.Context, query queries.GetPermissionsTreeQuery) ([]*dto.PermissionsDto, herrors.Herr) {
	//perms, err := h.permRepo.FindTree(ctx, query.Type)
	//if err != nil {
	//	return nil, herrors.QueryFail(err)
	//}
	//return dto.ToPermissionsDtoList(perms), nil
	return nil, nil
}
