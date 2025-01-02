package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type DepartmentCommandHandler struct {
	deptService *service.DepartmentService
}

func NewDepartmentCommandHandler(deptService *service.DepartmentService) *DepartmentCommandHandler {
	return &DepartmentCommandHandler{
		deptService: deptService,
	}
}

// HandleCreate 处理创建部门命令
func (h *DepartmentCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateDepartmentCommand) herrors.Herr {
	if validate := cmd.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", validate)
		return validate
	}

	// 2. 创建部门实体
	dept := model.NewDepartment(cmd.Code, cmd.Name, cmd.Sort)
	dept.ParentID = cmd.ParentID
	dept.Leader = cmd.Leader
	dept.Phone = cmd.Phone
	dept.Email = cmd.Email
	dept.Status = cmd.Status
	dept.Description = cmd.Description
	dept.TenantID = actx.GetTenantId(ctx)

	// 调用领域服务创建部门
	if err := h.deptService.CreateDepartment(ctx, dept); err != nil {
		hlog.CtxErrorf(ctx, "failed to create department: %s", err)
		return err
	}

	return nil
}

// HandleUpdate 处理更新部门命令
func (h *DepartmentCommandHandler) HandleUpdate(ctx context.Context, cmd *commands.UpdateDepartmentCommand) herrors.Herr {
	if validate := cmd.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", validate)
		return validate
	}

	// 1. 获取部门
	dept, err := h.deptService.GetByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to find department: %s", err)
		return herrors.UpdateFail(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(cmd.ID)
	}

	// 3. 更新部门信息
	dept.UpdateBasicInfo(cmd.Name, cmd.Code, cmd.Sort)
	dept.UpdateContactInfo(cmd.Leader, cmd.Phone, cmd.Email)
	dept.UpdateStatus(cmd.Status)
	dept.UpdateParent(cmd.ParentID)
	dept.Description = cmd.Description

	// 调用领域服务更新部门
	if err := h.deptService.UpdateDepartment(ctx, dept); err != nil {
		hlog.CtxErrorf(ctx, "failed to update department: %s", err)
		return herrors.UpdateFail(err)
	}

	return nil
}

// HandleDelete 处理删除部门命令
func (h *DepartmentCommandHandler) HandleDelete(ctx context.Context, cmd *commands.DeleteDepartmentCommand) herrors.Herr {
	if validate := cmd.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", validate)
		return validate
	}

	// 调用领域服务删除部门
	if err := h.deptService.DeleteDepartment(ctx, cmd.ID); err != nil {
		hlog.CtxErrorf(ctx, "failed to delete department: %s", err)
		return herrors.DeleteFail(err)
	}

	return nil
}

// HandleMove 处理移动部门命令
func (h *DepartmentCommandHandler) HandleMove(ctx context.Context, cmd *commands.MoveDepartmentCommand) herrors.Herr {
	if validate := cmd.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", validate)
		return validate
	}
	if err := h.deptService.MoveDepartment(ctx, cmd.ID, cmd.TargetParent); err != nil {
		hlog.CtxErrorf(ctx, "failed to move department: %s", err)
		return herrors.UpdateFail(err)
	}
	return nil
}

// HandleSetAdmin 处理设置部门管理员
func (h *DepartmentCommandHandler) HandleSetAdmin(ctx context.Context, cmd *commands.SetDepartmentAdminCommand) herrors.Herr {
	if hr := h.deptService.SetDepartmentAdmin(ctx, cmd.DeptID, cmd.AdminID); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to set department admin: %s", hr)
		return hr
	}

	return nil
}

// HandleAssignUsers 处理分配用户到部门
func (h *DepartmentCommandHandler) HandleAssignUsers(ctx context.Context, cmd *commands.AssignUsersToDepartmentCommand) herrors.Herr {
	// 调用领域服务分配用户
	if err := h.deptService.AssignUsers(ctx, cmd.DeptID, cmd.UserIDs); err != nil {
		hlog.CtxErrorf(ctx, "failed to assign users to department: %s", err)
		return herrors.UpdateFail(err)
	}

	return nil
}

// HandleRemoveUsers 处理从部门移除用户
func (h *DepartmentCommandHandler) HandleRemoveUsers(ctx context.Context, cmd *commands.RemoveUsersFromDepartmentCommand) herrors.Herr {
	// 调用领域服务移除用户
	if err := h.deptService.RemoveUsers(ctx, cmd.DeptID, cmd.UserIDs); err != nil {
		hlog.CtxErrorf(ctx, "failed to remove users from department: %s", err)
		return herrors.UpdateFail(err)
	}

	return nil
}

// HandleTransferUser 处理人员部门调动
func (h *DepartmentCommandHandler) HandleTransferUser(ctx context.Context, cmd *commands.TransferUserCommand) herrors.Herr {
	if hr := cmd.Validate(); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", hr)
		return hr
	}

	// 调用用户服务执行部门调动
	if hr := h.deptService.TransferUser(ctx, cmd.UserID, cmd.FromDeptID, cmd.ToDeptID); herrors.HaveError(hr) {
		hlog.CtxErrorf(ctx, "failed to transfer user: %s", hr)
		return hr
	}

	return nil
}
