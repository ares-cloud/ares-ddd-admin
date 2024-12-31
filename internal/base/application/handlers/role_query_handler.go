package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type RoleQueryHandler struct {
	roleQuery query.IRoleQueryService
	converter *converter.RoleConverter
}

func NewRoleQueryHandler(
	roleQuery query.IRoleQueryService,
	converter *converter.RoleConverter,
) *RoleQueryHandler {
	return &RoleQueryHandler{
		roleQuery: roleQuery,
		converter: converter,
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
	total, err := h.roleQuery.CountRoles(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to count roles: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 查询数据
	roles, err := h.roleQuery.FindRoles(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to find roles: %s", err)
		return nil, herrors.QueryFail(err)
	}
	return &models.PageRes[dto.RoleDto]{
		List:  roles,
		Total: total,
	}, nil
}

// HandleGet 处理获取角色查询
func (h *RoleQueryHandler) HandleGet(ctx context.Context, query queries.GetRoleQuery) (*dto.RoleDto, herrors.Herr) {
	// 查询角色
	role, err := h.roleQuery.GetRole(ctx, query.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get role: %s", err)
		return nil, herrors.QueryFail(err)
	}
	if role == nil {
		return nil, errors.RoleNotFound(query.Id)
	}

	// 查询角色权限
	perms, err := h.roleQuery.GetRolePermissions(ctx, query.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get role permissions: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO并填充权限ID
	role.PermIds = make([]int64, len(perms))
	if len(perms) > 0 {
		for i, perm := range perms {
			role.PermIds[i] = perm.ID
		}
	}

	return role, nil
}

func (h *RoleQueryHandler) HandleGetUserRoles(ctx context.Context, query queries.GetUserRolesQuery) ([]*dto.RoleDto, herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("user_id", db_query.Eq, query.UserID)

	// 查询用户角色
	roles, err := h.roleQuery.FindRoles(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get user roles: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return roles, nil
}

// HandleGetAllEnabled 获取所有启用状态的角色
func (h *RoleQueryHandler) HandleGetAllEnabled(ctx context.Context) ([]*dto.RoleDto, herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("status", db_query.Eq, 1) // 状态为启用
	qb.OrderBy("sequence", true)       // 按序号排序

	// 查询角色列表
	roles, err := h.roleQuery.FindRoles(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get enabled roles: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return roles, nil
}

// HandleGetAllDataPermission 获取所有数据权限角色
func (h *RoleQueryHandler) HandleGetAllDataPermission(ctx context.Context) ([]*dto.RoleDto, herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("type", db_query.Eq, 2)   // 类型为数据权限
	qb.Where("status", db_query.Eq, 1) // 状态为启用
	qb.OrderBy("sequence", true)       // 按序号排序

	// 查询角色列表
	roles, err := h.roleQuery.FindRoles(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get data permission roles: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return roles, nil
}
