package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	iQuery "github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type UserQueryHandler struct {
	queryService iQuery.UserQueryService
}

func NewUserQueryHandler(queryService iQuery.UserQueryService) *UserQueryHandler {
	return &UserQueryHandler{
		queryService: queryService,
	}
}

// HandleGet 处理获取用户详情查询
func (h *UserQueryHandler) HandleGet(ctx context.Context, q queries.GetUserQuery) (*dto.UserDto, herrors.Herr) {
	user, err := h.queryService.GetUser(ctx, q.ID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToUserDto(user), nil
}

// HandleList 处理用户列表查询
func (h *UserQueryHandler) HandleList(ctx context.Context, q *queries.ListUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Username != "" {
		qb.Where("username", query.Like, "%"+q.Username+"%")
	}
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	if q.Phone != "" {
		qb.Where("phone", query.Like, "%"+q.Phone+"%")
	}
	if q.Email != "" {
		qb.Where("email", query.Like, "%"+q.Email+"%")
	}
	if q.Status != 0 {
		qb.Where("status", query.Eq, q.Status)
	}
	qb.WithPage(&q.Page)

	// 查询总数
	total, err := h.queryService.CountUsers(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 查询数据
	users, err := h.queryService.FindUsers(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return &models.PageRes[dto.UserDto]{
		List:  dto.ToUserDtoList(users),
		Total: total,
	}, nil
}

// HandleGetUserInfo 处理获取用户信息查询
func (h *UserQueryHandler) HandleGetUserInfo(ctx context.Context, q queries.GetUserInfoQuery) (*dto.UserInfoDto, herrors.Herr) {
	// 1. 获取用户基本信息
	user, err := h.queryService.GetUser(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 2. 获取用户权限
	permissions, err := h.queryService.GetUserPermissions(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 3. 获取用户角色
	roles, err := h.queryService.GetUserRoles(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	roleCodes := make([]string, 0)
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}
	// 4. 构建用户信息DTO
	return &dto.UserInfoDto{
		User:        dto.ToUserDto(user),
		Roles:       roleCodes,
		HomePage:    "User",
		Permissions: permissions,
	}, nil
}

// HandleGetUserMenus 处理获取用户菜单查询
func (h *UserQueryHandler) HandleGetUserMenus(ctx context.Context, q queries.GetUserMenusQuery) ([]*dto.PermissionsTreeDto, herrors.Herr) {
	// 获取用户菜单
	menus, err := h.queryService.GetUserTreeMenus(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为树形结构DTO
	return dto.ToPermissionsTreeDtoList(menus), nil
}

// HandleGetUserPermissions 处理获取用户权限查询
func (h *UserQueryHandler) HandleGetUserPermissions(ctx context.Context, q *queries.GetUserPermissionsQuery) ([]string, herrors.Herr) {
	permissions, err := h.queryService.GetUserPermissions(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return permissions, nil
}

// HandleGetUserRoles 处理获取用户角色查询
func (h *UserQueryHandler) HandleGetUserRoles(ctx context.Context, q *queries.GetUserRolesQuery) ([]*dto.RoleDto, herrors.Herr) {
	roles, err := h.queryService.GetUserRoles(ctx, q.UserID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToRoleDtoList(roles), nil
}

// HandleListDepartmentUsers 处理部门用户列表查询
func (h *UserQueryHandler) HandleListDepartmentUsers(ctx context.Context, q *queries.ListDepartmentUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Username != "" {
		qb.Where("username", query.Like, "%"+q.Username+"%")
	}
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	qb.WithPage(&q.Page)

	// 查询总数
	total, err := h.queryService.CountUsersByDepartment(ctx, q.DeptID, q.ExcludeAdminID, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 查询数据
	users, err := h.queryService.FindUsersByDepartment(ctx, q.DeptID, q.ExcludeAdminID, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return &models.PageRes[dto.UserDto]{
		List:  dto.ToUserDtoList(users),
		Total: total,
	}, nil
}

// 其他查询处理方法...
