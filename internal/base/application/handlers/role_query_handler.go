package handlers

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type RoleQueryHandler struct {
	roleRepo repository.IRoleRepository
}

func NewRoleQueryHandler(
	roleRepo repository.IRoleRepository,
) *RoleQueryHandler {
	return &RoleQueryHandler{
		roleRepo: roleRepo,
	}
}

// HandleList 处理列表查询
func (h *RoleQueryHandler) HandleList(ctx context.Context, q *queries.ListRolesQuery) (*models.PageRes[dto.RoleDto], herrors.Herr) {
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
	if q.Type != 0 {
		qb.Where("type", db_query.Eq, q.Type)
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

// HandleGet 处理获取角色查询
func (h *RoleQueryHandler) HandleGet(ctx context.Context, query queries.GetRoleQuery) (*dto.RoleDto, herrors.Herr) {
	// 1. 查询角色基本信息
	role, err := h.roleRepo.FindByID(ctx, query.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get role: %s", err)
		return nil, herrors.QueryFail(err)
	}
	if role == nil {
		return nil, herrors.QueryFail(fmt.Errorf("role not found: %d", query.Id))
	}

	// 2. 查询角色权限
	permIDs, err := h.roleRepo.GetPermissionsByRoleID(ctx, query.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get role permissions: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 3. 转换为DTO并填充权限ID
	roleDto := dto.ToRoleDto(role)
	roleDto.PermIds = permIDs

	return roleDto, nil
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

// HandleGetAllDataPermission 处理获取所有数据权限角色
func (h *RoleQueryHandler) HandleGetAllDataPermission(ctx context.Context) ([]*dto.RoleDto, herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("type", db_query.Eq, int8(model.RoleTypeData))
	qb.Where("status", db_query.Eq, 1) // 只查询启用状态的角色
	qb.OrderBy("sequence", false)      // 按sequence排序

	// 查询角色
	roles, err := h.roleRepo.Find(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get data permission roles: %s", err)
		return nil, herrors.NewServerHError(err)
	}

	// 使用已有的DTO转换方法
	return dto.ToRoleDtoList(roles), nil
}
