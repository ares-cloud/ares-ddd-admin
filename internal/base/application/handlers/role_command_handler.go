package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type RoleCommandHandler struct {
	roleService *service.RoleCommandService
}

func NewRoleCommandHandler(roleService *service.RoleCommandService) *RoleCommandHandler {
	return &RoleCommandHandler{
		roleService: roleService,
	}
}

// HandleCreate 处理创建角色命令
func (h *RoleCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateRoleCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 创建角色领域模型
	role := model.NewRole(actx.GetTenantId(ctx), cmd.Code, cmd.Name)
	role.Description = cmd.Description
	role.Localize = cmd.Localize
	role.Sequence = cmd.Sequence
	role.Type = cmd.Type

	// 创建角色
	return h.roleService.CreateRole(ctx, role)
}

// HandleUpdate 处理更新角色命令
func (h *RoleCommandHandler) HandleUpdate(ctx context.Context, cmd *commands.UpdateRoleCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 获取现有角色
	role, hr := h.roleService.GetRole(ctx, cmd.ID)
	if herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to get role: %s", hr)
		return hr
	}

	// 更新基本信息
	role.UpdateBasicInfo(cmd.Name, cmd.Localize, cmd.Description, cmd.Sequence)
	if cmd.Status != 0 {
		if hr := role.UpdateStatus(cmd.Status); herrors.HaveError(hr) {
			return hr
		}
	}

	// 保存更新
	return h.roleService.UpdateRole(ctx, role)
}

// HandleDelete 处理删除角色命令
func (h *RoleCommandHandler) HandleDelete(ctx context.Context, cmd *commands.DeleteRoleCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 删除角色
	return h.roleService.DeleteRole(ctx, cmd.ID)
}

// HandleAssignPermissions 处理分配权限命令
func (h *RoleCommandHandler) HandleAssignPermissions(ctx context.Context, cmd *commands.AssignRolePermissionsCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 分配权限
	return h.roleService.AssignPermissions(ctx, cmd.RoleID, cmd.PermissionIDs)
}
