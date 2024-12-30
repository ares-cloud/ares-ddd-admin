package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type TenantCommandHandler struct {
	tenantService *service.TenantCommandService
}

func NewTenantCommandHandler(tenantService *service.TenantCommandService) *TenantCommandHandler {
	return &TenantCommandHandler{
		tenantService: tenantService,
	}
}

func (h *TenantCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateTenantCommand) herrors.Herr {
	// 检查租户编码是否已存在
	exists, err := h.tenantService.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		return err
	}
	if exists {
		return herrors.ErrRecordNotFount
	}

	// 创建管理员用户
	adminUser := model.NewUser("", cmd.AdminUser.Nickname, cmd.AdminUser.Password)
	adminUser.Phone = cmd.AdminUser.Phone
	adminUser.Email = cmd.AdminUser.Email
	adminUser.Username = cmd.AdminUser.Username
	if err := adminUser.HashPassword(); err != nil {
		hlog.CtxErrorf(ctx, "hash password: %v", err)
		return herrors.CreateFail(err)
	}

	// 创建租户
	tenant := model.NewTenant(cmd.Code, cmd.Name, adminUser)
	tenant.Description = cmd.Description
	tenant.IsDefault = cmd.IsDefault
	if cmd.ExpireTime > 0 {
		tenant.ExpireTime = cmd.ExpireTime
	}

	if err := h.tenantService.CreateTenant(ctx, tenant); err != nil {
		hlog.CtxErrorf(ctx, "create tenant err: %v", err)
		return err
	}
	return nil
}

func (h *TenantCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateTenantCommand) herrors.Herr {
	// 查找现有租户
	tenant, err := h.tenantService.GetTenant(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// 更新基本信息
	tenant.UpdateBasicInfo(cmd.Name, cmd.Description)

	// 更新过期时间
	if cmd.ExpireTime > 0 {
		tenant.UpdateExpireTime(cmd.ExpireTime)
	}

	// 更新默认状态
	if cmd.IsDefault != 0 {
		if err := tenant.UpdateIsDefault(cmd.IsDefault); err != nil {
			return herrors.UpdateFail(err)
		}
	}

	// 保存更新
	if err := h.tenantService.UpdateTenant(ctx, tenant); err != nil {
		return err
	}

	return nil
}

func (h *TenantCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteTenantCommand) herrors.Herr {
	// 删除租户
	if err := h.tenantService.DeleteTenant(ctx, cmd.ID); err != nil {
		return err
	}
	return nil
}

func (h *TenantCommandHandler) HandleAssignPermissions(ctx context.Context, cmd commands.AssignTenantPermissionsCommand) herrors.Herr {
	// 分配权限
	if err := h.tenantService.AssignPermissions(ctx, cmd.TenantID, cmd.PermissionIDs); err != nil {
		return err
	}
	return nil
}
