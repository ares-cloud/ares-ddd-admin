package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type UserCommandHandler struct {
	userService *service.UserCommandService
}

func NewUserCommandHandler(
	userService *service.UserCommandService,
) *UserCommandHandler {
	return &UserCommandHandler{
		userService: userService,
	}
}

// HandleCreate 处理创建用户请求
func (h *UserCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateUserCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 创建用户领域模型
	user := model.NewUser(actx.GetTenantId(ctx), cmd.Username, cmd.Password)
	user.Phone = cmd.Phone
	user.Email = cmd.Email
	user.Nickname = cmd.Nickname
	user.Avatar = cmd.Avatar

	// 加密密码
	if err := user.HashPassword(); err != nil {
		hlog.CtxErrorf(ctx, "failed to hash password: %s", err)
		return herrors.NewServerHError(err)
	}

	// 创建用户
	if hr := h.userService.CreateUser(ctx, user); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to create user: %s", hr)
		return hr
	}

	// 分配角色
	if len(cmd.RoleIDs) > 0 {
		if hr := h.userService.AssignRoles(ctx, user.ID, cmd.RoleIDs); herrors.HaveError(hr) {
			hlog.CtxErrorf(ctx, "failed to assign roles: %s", hr)
			return hr
		}
	}

	return nil
}

// HandleUpdate 处理更新用户请求
func (h *UserCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateUserCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 获取现有用户
	user, hr := h.userService.GetUser(ctx, cmd.ID)
	if herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to get user: %s", hr)
		return hr
	}

	// 更新基本信息
	user.UpdateBasicInfo(cmd.Name, cmd.Nickname, cmd.Phone, cmd.Email, cmd.Avatar, "")

	// 更新状态
	if cmd.Status != 0 {
		if hr := user.UpdateStatus(cmd.Status); herrors.HaveError(hr) {
			return hr
		}
	}

	// 更新角色
	if len(cmd.RoleIDs) > 0 {
		if hr := h.userService.AssignRoles(ctx, user.ID, cmd.RoleIDs); herrors.HaveError(hr) {
			hlog.CtxErrorf(ctx, "failed to assign roles: %s", hr)
			return hr
		}
	}

	// 保存更新
	if hr := h.userService.UpdateUser(ctx, user); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to update user: %s", hr)
		return hr
	}

	return nil
}

// HandleDelete 处理删除用户请求
func (h *UserCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteUserCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 删除用户
	if hr := h.userService.DeleteUser(ctx, cmd.ID); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to delete user: %s", hr)
		return hr
	}
	return nil
}

// HandleUpdateStatus 处理更新用户状态请求
func (h *UserCommandHandler) HandleUpdateStatus(ctx context.Context, cmd commands.UpdateUserStatusCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 获取用户
	user, hr := h.userService.GetUser(ctx, cmd.ID)
	if herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to get user: %s", hr)
		return hr
	}

	// 更新状态
	if hr := user.UpdateStatus(cmd.Status); herrors.HaveError(hr) {
		return hr
	}

	// 保存更新
	if hr := h.userService.UpdateUser(ctx, user); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to update user status: %s", hr)
		return hr
	}

	return nil
}
