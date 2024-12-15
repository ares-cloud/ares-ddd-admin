package handlers

import (
	"context"
	"errors"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type TenantCommandHandler struct {
	tenantRepo repository.ITenantRepository
}

func NewTenantCommandHandler(tenantRepo repository.ITenantRepository) *TenantCommandHandler {
	return &TenantCommandHandler{
		tenantRepo: tenantRepo,
	}
}

func (h *TenantCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateTenantCommand) herrors.Herr {
	// 检查租户编码是否已存在
	exists, err := h.tenantRepo.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		return herrors.CreateFail(err)
	}
	if exists {
		return herrors.CreateFail(errors.New("tenant code already exists"))
	}

	// 创建管理员用户
	adminUser := model.NewUser(cmd.AdminUser.Username, cmd.AdminUser.Name, cmd.AdminUser.Password)
	adminUser.Phone = cmd.AdminUser.Phone
	adminUser.Email = cmd.AdminUser.Email

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

	err = h.tenantRepo.Create(ctx, tenant)
	if err != nil {
		hlog.CtxErrorf(ctx, "create tenant err: %v", err)
		return herrors.CreateFail(err)
	}
	return nil
}

func (h *TenantCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateTenantCommand) herrors.Herr {
	// 查找现有租户
	tenant, err := h.tenantRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return herrors.UpdateFail(err)
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
	if err := h.tenantRepo.Update(ctx, tenant); err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}

func (h *TenantCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteTenantCommand) herrors.Herr {
	// 查找租户
	tenant, err := h.tenantRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return herrors.DeleteFail(err)
	}

	// 检查是否为默认租户
	if tenant.IsDefaultTenant() {
		return herrors.DeleteFail(errors.New("cannot delete default tenant"))
	}

	// 删除租户
	if err := h.tenantRepo.Delete(ctx, cmd.ID); err != nil {
		return herrors.DeleteFail(err)
	}

	return nil
}

func (h *TenantCommandHandler) HandleAssignPermissions(ctx context.Context, cmd commands.AssignTenantPermissionsCommand) herrors.Herr {
	// 查找租户
	tenant, err := h.tenantRepo.FindByID(ctx, cmd.TenantID)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	// 分配权限
	err = h.tenantRepo.AssignPermissions(ctx, tenant.ID, cmd.PermissionIDs)
	if err != nil {
		return herrors.UpdateFail(err)
	}

	return nil
}
