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

type RoleQueryHandler struct {
	roleRepo repository.IRoleRepository
}

func NewRoleQueryHandler(roleRepo repository.IRoleRepository) *RoleQueryHandler {
	return &RoleQueryHandler{
		roleRepo: roleRepo,
	}
}

// HandleList 处理列表查询
func (h *RoleQueryHandler) HandleList(ctx context.Context, q *queries.ListRolesQuery) (*models.PageRes[dto.RoleDto], herrors.Herr) {
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
	total, err := h.roleRepo.Count(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	roles, err := h.roleRepo.Find(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToRoleDtoList(roles)

	return &models.PageRes[dto.RoleDto]{
		List:  dtoList,
		Total: total,
	}, nil
}

// HandleGet 获取单个角色
func (h *RoleQueryHandler) HandleGet(ctx context.Context, query queries.GetRoleQuery) (*dto.RoleDto, herrors.Herr) {
	role, err := h.roleRepo.FindByID(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToRoleDto(role), nil
}

func (h *RoleQueryHandler) HandleGetUserRoles(ctx context.Context, query queries.GetUserRolesQuery) ([]*dto.RoleDto, herrors.Herr) {
	roles, err := h.roleRepo.FindByUserID(ctx, query.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToRoleDtoList(roles), nil
}

// HandleGetAllEnabled 获取所有启用状态的角色
func (h *RoleQueryHandler) HandleGetAllEnabled(ctx context.Context) ([]*dto.RoleDto, herrors.Herr) {
	roles, err := h.roleRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToRoleDtoList(roles), nil
}
