package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type DataPermissionCommandHandler struct {
	permRepo repository.IDataPermissionRepository
}

func NewDataPermissionCommandHandler(permRepo repository.IDataPermissionRepository) *DataPermissionCommandHandler {
	return &DataPermissionCommandHandler{
		permRepo: permRepo,
	}
}

// HandleAssign 处理分配数据权限
func (h *DataPermissionCommandHandler) HandleAssign(ctx context.Context, cmd *commands.AssignDataPermissionCommand) herrors.Herr {
	// 1. 检查角色是否已有数据权限
	perm, err := h.permRepo.GetByRoleID(ctx, cmd.RoleID)
	if err != nil && !database.IfErrorNotFound(err) {
		return herrors.NewServerHError(err)
	}

	// 如果没有指定数据范围,使用默认的部门及下级数据权限
	if cmd.Scope == 0 {
		cmd.Scope = int8(model.DataScopeDeptTree)
	}

	// 2. 创建或更新数据权限
	if perm == nil {
		// 创建新的数据权限
		perm = &model.DataPermission{
			RoleID:   cmd.RoleID,
			Scope:    model.DataScope(cmd.Scope),
			DeptIDs:  cmd.DeptIDs,
			TenantID: actx.GetTenantId(ctx),
		}
		if err := h.permRepo.Create(ctx, perm); err != nil {
			return herrors.NewServerHError(err)
		}
	} else {
		// 更新现有数据权限
		perm.Scope = model.DataScope(cmd.Scope)
		perm.DeptIDs = cmd.DeptIDs
		if err := h.permRepo.Update(ctx, perm); err != nil {
			return herrors.NewServerHError(err)
		}
	}

	return nil
}

// HandleRemove 处理移除数据权限
func (h *DataPermissionCommandHandler) HandleRemove(ctx context.Context, cmd *commands.RemoveDataPermissionCommand) herrors.Herr {
	// 1. 获取角色的数据权限
	perm, err := h.permRepo.GetByRoleID(ctx, cmd.RoleID)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil
		}
		return herrors.NewServerHError(err)
	}

	// 2. 删除数据权限
	if err := h.permRepo.Delete(ctx, perm.ID); err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}
