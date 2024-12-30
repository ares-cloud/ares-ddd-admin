package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type UserCommandHandler struct {
	userCommandService *service.UserCommandService
	queryService       query.UserQueryService
	roleService        repository.IRoleRepository
}

func NewUserCommandHandler(
	userCommandService *service.UserCommandService,
	queryService query.UserQueryService,
	roleService repository.IRoleRepository,
) *UserCommandHandler {
	return &UserCommandHandler{
		userCommandService: userCommandService,
		queryService:       queryService,
		roleService:        roleService,
	}
}

// 错误转换函数
func convertUserError(err error) herrors.Herr {
	if err == nil {
		return nil
	}

	switch err {
	case errors.ErrUserNotFound:
		return herrors.ErrRecordNotFount
	case errors.ErrUsernameExists:
		return herrors.DataIsExist
	case errors.ErrInvalidCredentials:
		return herrors.NewBadReqError("invalid username or password")
	case errors.ErrUserDisabled:
		return herrors.NewBadReqError("user is disabled")
	case errors.ErrPasswordMismatch:
		return herrors.NewBadReqError("password mismatch")
	default:
		return herrors.NewServerHError(err)
	}
}

// HandleCreate 处理创建用户请求
func (h *UserCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateUserCommand) herrors.Herr {
	// 创建用户领域模型
	user := model.NewUser(cmd.Username, cmd.Name, cmd.Password)
	user.Phone = cmd.Phone
	user.Email = cmd.Email
	user.InvitationCode = cmd.InvitationCode

	// 加密密码
	if err := user.HashPassword(); err != nil {
		return herrors.CreateFail(err)
	}

	// 创建用户
	if err := h.userCommandService.CreateUser(ctx, user); err != nil {
		return convertUserError(err)
	}

	// 分配角色
	if len(cmd.RoleIDs) > 0 {
		if err := h.userCommandService.AssignRoles(ctx, user.ID, cmd.RoleIDs); err != nil {
			return convertUserError(err)
		}
	}

	return nil
}

// HandleUpdate 处理更新用户请求
func (h *UserCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateUserCommand) herrors.Herr {
	// 获取现有用户
	user, err := h.queryService.GetUser(ctx, cmd.ID)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	// 更新基本信息
	user.UpdateBasicInfo(cmd.Name, cmd.Phone, cmd.Email, cmd.FaceURL, cmd.Remark)

	// 更新状态
	if cmd.Status != 0 {
		if err := user.UpdateStatus(cmd.Status); err != nil {
			return herrors.UpdateFail(err)
		}
	}

	// 更新角色
	if len(cmd.RoleIDs) > 0 {
		//if err := h.userCommandService.AssignRoles(ctx, user.ID, cmd.RoleIDs); err != nil {
		//	return herrors.UpdateFail(err)
		//}
		roles, err := h.roleService.FindByIDs(ctx, cmd.RoleIDs)
		if err != nil {
			return herrors.UpdateFail(err)
		}
		user.Roles = roles
	}

	// 保存更新
	if err := h.userCommandService.UpdateUser(ctx, user); err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}

// HandleDelete 处理删除用户请求
func (h *UserCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteUserCommand) herrors.Herr {
	// 删除用户
	if err := h.userCommandService.DeleteUser(ctx, cmd.ID); err != nil {
		if err == errors.ErrUserNotFound {
			return herrors.ErrRecordNotFount
		}
		return herrors.DeleteFail(err)
	}
	return nil
}

// HandleUpdateStatus 处理更新用户状态请求
func (h *UserCommandHandler) HandleUpdateStatus(ctx context.Context, cmd commands.UpdateUserStatusCommand) herrors.Herr {
	// 获取用户
	user, err := h.queryService.GetUser(ctx, cmd.ID)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	// 更新状态
	if err = user.UpdateStatus(cmd.Status); err != nil {
		return herrors.UpdateFail(err)
	}

	// 保存更新
	if err = h.userCommandService.UpdateUser(ctx, user); err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}
