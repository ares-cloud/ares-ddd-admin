package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/service"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type UserQueryHandler struct {
	userRepo repository.IUserRepository
	uds      *service.UserService
}

func NewUserQueryHandler(userRepo repository.IUserRepository, uds *service.UserService) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo: userRepo,
		uds:      uds,
	}
}

// HandleList 处理列表查询
func (h *UserQueryHandler) HandleList(ctx context.Context, q *queries.ListUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()

	// 添加查询条件
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

	// 设置分页
	qb.WithPage(&q.Page)

	// 获取总数
	total, err := h.userRepo.Count(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	users, err := h.userRepo.Find(ctx, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToUserDtoList(users)

	return &models.PageRes[dto.UserDto]{
		List:  dtoList,
		Total: total,
	}, nil
}

// HandleGet 获取单个用户
func (h *UserQueryHandler) HandleGet(ctx context.Context, query queries.GetUserQuery) (*dto.UserDto, herrors.Herr) {
	user, err := h.userRepo.FindByID(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToUserDto(user), nil
}

// HandleGetUserInfo 获取用户信息
func (h *UserQueryHandler) HandleGetUserInfo(ctx context.Context, query queries.GetUserInfoQuery) (*dto.UserInfoDto, herrors.Herr) {
	// 获取用户信息
	user, err := h.userRepo.FindByID(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 获取用户所有权限
	permissions, err := h.uds.GetUserPermissions(ctx, query.Id)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
	}
	return dto.ToUserInfoDto(user, permissions, roles), nil
}

// HandleGetUserMenus 获取用户菜单树
func (h *UserQueryHandler) HandleGetUserMenus(ctx context.Context, query queries.GetUserMenusQuery) ([]*dto.PermissionsDto, herrors.Herr) {
	menus, err := h.uds.GetUserMenus(ctx, query.Id) // 1表示菜单类型
	if err != nil {
		return nil, herrors.QueryFail(err)
	}
	return dto.ToPermissionsDtoList(menus), nil
}
