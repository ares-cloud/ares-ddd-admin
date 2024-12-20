package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type PermissionsCommandHandler struct {
	permRepo repository.IPermissionsRepository
	ef       *casbin.Enforcer
}

func NewPermissionsCommandHandler(permRepo repository.IPermissionsRepository, ef *casbin.Enforcer) *PermissionsCommandHandler {
	return &PermissionsCommandHandler{
		permRepo: permRepo,
		ef:       ef,
	}
}

func (h *PermissionsCommandHandler) HandleCreate(ctx context.Context, cmd commands.CreatePermissionsCommand) herrors.Herr {
	// 检查权限编码是否已存在
	exists, err := h.permRepo.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "check permission code exists failed: %s", err)
		return herrors.CreateFail(err)
	}
	if exists {
		return herrors.CreateFail(fmt.Errorf("permission code %s already exists", cmd.Code))
	}

	perm := model.NewPermissions(cmd.Code, cmd.Name, cmd.Type, cmd.Sequence)
	perm.Localize = cmd.Localize
	perm.Icon = cmd.Icon
	perm.Description = cmd.Description
	perm.Path = cmd.Path
	perm.Properties = cmd.Properties
	perm.ParentID = cmd.ParentID

	// 添加资源
	for _, resource := range cmd.Resources {
		perm.AddResource(resource.Method, resource.Path)
	}

	err = h.permRepo.Create(ctx, perm)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission create failed: %s", err)
		return herrors.CreateFail(err)
	}
	return nil
}

func (h *PermissionsCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdatePermissionsCommand) herrors.Herr {
	perm, err := h.permRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission find failed: %s", err)
		return herrors.UpdateFail(err)
	}

	// 更新基本信息
	perm.UpdateBasicInfo(cmd.Name, cmd.Description, cmd.Sequence)
	if cmd.Status != nil {
		perm.UpdateStatus(*cmd.Status)
	}

	perm.Icon = cmd.Icon
	perm.Path = cmd.Path
	perm.Properties = cmd.Properties
	perm.ChangeType(cmd.Type)
	perm.ChangeParentID(cmd.ParentID)
	perm.Localize = cmd.Localize
	// 更新资源列表
	if len(cmd.Resources) > 0 {
		resources := make([]*model.PermissionsResource, len(cmd.Resources))
		for i, r := range cmd.Resources {
			resources[i] = &model.PermissionsResource{
				Method: r.Method,
				Path:   r.Path,
			}
		}
		perm.UpdateResources(resources)
	}

	err = h.permRepo.Update(ctx, perm)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission update failed %s", err)
		return herrors.UpdateFail(err)
	}

	// 发布权限更新消息
	if err := h.ef.PublishUpdate(ctx); err != nil {
		hlog.CtxErrorf(ctx, "publish permission update error: %v", err)
	}

	return nil
}

// HandleDelete 处理删除权限命令
func (h *PermissionsCommandHandler) HandleDelete(ctx context.Context, id int64) herrors.Herr {
	// 查找权限是否存在
	perm, err := h.permRepo.FindByID(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission find failed: %s", err)
		return herrors.DeleteFail(err)
	}

	// 检查是否有子权限
	if len(perm.Children) > 0 {
		return herrors.DeleteFail(errors.New("cannot delete permission with children"))
	}

	// 执行删除
	err = h.permRepo.Delete(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "permission delete failed: %s", err)
		return herrors.DeleteFail(err)
	}

	return nil
}
