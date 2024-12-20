package handlers

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type RoleCommandHandler struct {
	roleRepo repository.IRoleRepository
	permRepo repository.IPermissionsRepository
	ef       *casbin.Enforcer
}

func NewRoleCommandHandler(roleRepo repository.IRoleRepository, permRepo repository.IPermissionsRepository, ef *casbin.Enforcer) *RoleCommandHandler {
	return &RoleCommandHandler{
		roleRepo: roleRepo,
		permRepo: permRepo,
		ef:       ef,
	}
}

func (h *RoleCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateRoleCommand) herrors.Herr {
	// 检查角色编码是否已存在
	exists, err := h.roleRepo.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "check role code exists failed: %s", err)
		return herrors.CreateFail(err)
	}
	if exists {
		return herrors.CreateFail(fmt.Errorf("role code %s already exists", cmd.Code))
	}

	role := model.NewRole(cmd.Code, cmd.Name, cmd.Sequence)
	role.Description = cmd.Description
	role.Localize = cmd.Localize

	if len(cmd.PermIDs) > 0 {
		perms := make([]*model.Permissions, 0, len(cmd.PermIDs))
		for _, permID := range cmd.PermIDs {
			perm, err := h.permRepo.FindByID(ctx, permID)
			if err != nil {
				hlog.CtxErrorf(ctx, "failed to find perm with id %s", permID)
				return herrors.CreateFail(err)
			}
			perms = append(perms, perm)
		}
		role.AssignPermissions(perms)
	}

	err = h.roleRepo.Create(ctx, role)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to create role: %s", err)
		return herrors.CreateFail(err)
	}

	// 发布权限更新消息
	if err := h.ef.PublishUpdate(ctx); err != nil {
		hlog.CtxErrorf(ctx, "publish permission update error: %v", err)
	}

	return nil
}

func (h *RoleCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateRoleCommand) herrors.Herr {
	role, err := h.roleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to find role: %s", err)
		return herrors.CreateFail(err)
	}

	role.UpdateBasicInfo(cmd.Name, cmd.Description, cmd.Sequence)
	if cmd.Status != nil {
		role.UpdateStatus(*cmd.Status)
	}

	if len(cmd.PermIDs) > 0 {
		perms := make([]*model.Permissions, 0, len(cmd.PermIDs))
		for _, permID := range cmd.PermIDs {
			perm, err := h.permRepo.FindByID(ctx, permID)
			if err != nil {
				hlog.CtxErrorf(ctx, "failed to find perm with id %s", permID)
				return herrors.CreateFail(err)
			}
			perms = append(perms, perm)
		}
		role.AssignPermissions(perms)
	} else {
		role.AssignPermissions(nil)
	}
	role.UpdateLocalize(cmd.Localize)
	err = h.roleRepo.Update(ctx, role)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to update role: %s", err)
		return herrors.CreateFail(err)
	}
	return nil
}

func (h *RoleCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteRoleCommand) herrors.Herr {
	// 检查角色是否存在
	role, err := h.roleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to find role: %s", err)
		return herrors.CreateFail(err)
	}

	err = h.roleRepo.Delete(ctx, role.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to delete role: %s", err)
		return herrors.CreateFail(err)
	}
	return nil
}
