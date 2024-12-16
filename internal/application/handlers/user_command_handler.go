package handlers

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type UserCommandHandler struct {
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

func NewUserCommandHandler(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) *UserCommandHandler {
	return &UserCommandHandler{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (h *UserCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateUserCommand) herrors.Herr {
	user := model.NewUser(cmd.Username, cmd.Name, cmd.Password)
	user.Phone = cmd.Phone
	user.Email = cmd.Email
	user.InvitationCode = cmd.InvitationCode

	if err := user.HashPassword(); err != nil {
		hlog.CtxErrorf(ctx, "hash password: %v", err)
		return herrors.CreateFail(err)
	}

	if len(cmd.RoleIDs) > 0 {
		roles := make([]*model.Role, 0, len(cmd.RoleIDs))
		for _, roleID := range cmd.RoleIDs {
			role, err := h.roleRepo.FindByID(ctx, roleID)
			if err != nil {
				return herrors.CreateFail(err)
			}
			roles = append(roles, role)
		}
		user.AssignRoles(roles)
	}

	err := h.userRepo.Create(ctx, user)
	if err != nil {
		hlog.CtxErrorf(ctx, "create user err: %v", err)
		return herrors.CreateFail(err)
	}
	return nil
}

// HandleUpdate 处理更新用户请求
func (h *UserCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateUserCommand) herrors.Herr {
	// 查找现有用户
	user, err := h.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	// 检查用户名是否已存在
	if cmd.Username != user.Username {
		exists, err := h.userRepo.ExistsByUsername(ctx, cmd.Username)
		if err != nil {
			return herrors.UpdateFail(err)
		}
		if exists {
			return herrors.UpdateFail(errors.New("username already exists"))
		}
	}

	// 更新基本信息
	user.UpdateBasicInfo(cmd.Name, cmd.Phone, cmd.Email, cmd.FaceURL, cmd.Remark)

	// 如果需要更新状态
	if cmd.Status != 0 {
		if err := user.UpdateStatus(cmd.Status); err != nil {
			return herrors.UpdateFail(err)
		}
	}

	// 如果提供了角色列表，更新角色
	if len(cmd.RoleIDs) > 0 {
		roles, err := h.roleRepo.FindByIDs(ctx, cmd.RoleIDs)
		if err != nil {
			return herrors.UpdateFail(err)
		}
		user.AssignRoles(roles)
	} else {
		user.AssignRoles(nil)
	}

	// 保存更新
	if err := h.userRepo.Update(ctx, user); err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}

// HandleDelete 处理删除用户请求
func (h *UserCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteUserCommand) herrors.Herr {
	// 查找用户
	_, err := h.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return herrors.ErrRecordNotFount
	}

	// 删除用户
	if err := h.userRepo.Delete(ctx, cmd.ID); err != nil {
		return herrors.DeleteFail(err)
	}

	return nil
}

// HandleUpdateStatus 处理更新用户状态请求
func (h *UserCommandHandler) HandleUpdateStatus(ctx context.Context, cmd commands.UpdateUserStatusCommand) herrors.Herr {
	// 查找用户
	user, err := h.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	// 更新状态
	if err = user.UpdateStatus(cmd.Status); err != nil {
		return herrors.UpdateFail(err)
	}

	// 保存更新
	if err = h.userRepo.Update(ctx, user); err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}
