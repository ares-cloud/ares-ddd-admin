package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type PermissionsQueryHandler struct {
	permRepo repository.IPermissionsRepository
	pds      *service.PermissionService
}

func NewPermissionsQueryHandler(permRepo repository.IPermissionsRepository, pds *service.PermissionService) *PermissionsQueryHandler {
	return &PermissionsQueryHandler{
		permRepo: permRepo,
		pds:      pds,
	}
}

// HandleList 处理列表查询
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
	perms, err := h.permRepo.FindTreeByType(ctx, 1)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToPermissionsDtoList(perms), nil
}

func (h *PermissionsQueryHandler) HandleGetPermissionsTree(ctx context.Context) (*dto.PermissionsTreeResult, herrors.Herr) {
	// 获取所有权限并构建树
	permissions, ids, err := h.permRepo.FindAllTree(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return &dto.PermissionsTreeResult{
		Tree: dto.ToPermissionsTreeDtoList(permissions),
		Ids:  ids,
	}, nil
}

// HandleGetAllEnabled 获取所有启用状态的权限
func (h *PermissionsQueryHandler) HandleGetAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, herrors.Herr) {
	permissions, err := h.pds.FindAllEnabled(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToPermissionsDtoList(permissions), nil
}
