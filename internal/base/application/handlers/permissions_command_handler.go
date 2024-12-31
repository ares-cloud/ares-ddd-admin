package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type PermissionsCommandHandler struct {
	permService *service.PermissionService
	ef          *casbin.Enforcer
}

func NewPermissionsCommandHandler(
	permService *service.PermissionService,
	ef *casbin.Enforcer,
) *PermissionsCommandHandler {
	return &PermissionsCommandHandler{
		permService: permService,
		ef:          ef,
	}
}

func (h *PermissionsCommandHandler) HandleCreate(ctx context.Context, cmd commands.CreatePermissionsCommand) herrors.Herr {
	perm := model.NewPermissions(cmd.Code, cmd.Name, cmd.Type, cmd.Sequence)
	perm.Localize = cmd.Localize
	perm.Icon = cmd.Icon
	perm.Description = cmd.Description
	perm.Path = cmd.Path
	perm.Properties = cmd.Properties
	perm.ParentID = cmd.ParentID

	// 添加资源
	for _, resource := range cmd.Resources {
		if err := perm.AddResource(resource.Method, resource.Path); err != nil {
			hlog.CtxErrorf(ctx, "add resource failed: %s", err)
			return err
		}
	}

	if err := h.permService.CreatePermission(ctx, perm); err != nil {
		hlog.CtxErrorf(ctx, "permission create failed: %s", err)
		return err
	}
	return nil
}

func (h *PermissionsCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdatePermissionsCommand) herrors.Herr {
	// 1. 查询权限
	perm, err := h.permService.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission find failed: %s", err)
		return err
	}

	// 2. 更新基本信息
	if err := perm.UpdateBasicInfo(cmd.Name, cmd.Description, cmd.Sequence); err != nil {
		hlog.CtxErrorf(ctx, "update basic info failed: %s", err)
		return err
	}

	if cmd.Status != nil {
		if err := perm.UpdateStatus(*cmd.Status); err != nil {
			hlog.CtxErrorf(ctx, "update status failed: %s", err)
			return err
		}
	}

	perm.Icon = cmd.Icon
	perm.Path = cmd.Path
	perm.Properties = cmd.Properties
	if err := perm.ChangeType(cmd.Type); err != nil {
		hlog.CtxErrorf(ctx, "change type failed: %s", err)
		return err
	}
	if err := perm.ChangeParentID(cmd.ParentID); err != nil {
		hlog.CtxErrorf(ctx, "change parent id failed: %s", err)
		return err
	}
	perm.Localize = cmd.Localize

	// 3. 更新资源列表
	if len(cmd.Resources) > 0 {
		resources := make([]*model.PermissionsResource, len(cmd.Resources))
		for i, r := range cmd.Resources {
			resources[i] = &model.PermissionsResource{
				Method: r.Method,
				Path:   r.Path,
			}
		}
		if err := perm.UpdateResources(resources); err != nil {
			hlog.CtxErrorf(ctx, "update resources failed: %s", err)
			return err
		}
	}

	// 4. 更新权限
	if err := h.permService.UpdatePermission(ctx, perm); err != nil {
		hlog.CtxErrorf(ctx, "permission update failed %s", err)
		return err
	}

	// 5. 发布权限更新消息
	if err := h.ef.PublishUpdate(ctx); err != nil {
		hlog.CtxErrorf(ctx, "publish permission update error: %v", err)
	}

	return nil
}

func (h *PermissionsCommandHandler) HandleDelete(ctx context.Context, id int64) herrors.Herr {
	if err := h.permService.DeletePermission(ctx, id); err != nil {
		hlog.CtxErrorf(ctx, "permission delete failed: %s", err)
		return err
	}
	return nil
}
