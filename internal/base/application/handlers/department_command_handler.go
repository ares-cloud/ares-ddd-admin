package handlers

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type DepartmentCommandHandler struct {
	deptRepo    repository.IDepartmentRepository
	deptService *service.DepartmentService
}

func NewDepartmentCommandHandler(deptRepo repository.IDepartmentRepository, deptService *service.DepartmentService) *DepartmentCommandHandler {
	return &DepartmentCommandHandler{
		deptRepo:    deptRepo,
		deptService: deptService,
	}
}

// HandleCreate 处理创建部门命令
func (h *DepartmentCommandHandler) HandleCreate(ctx context.Context, cmd *commands.CreateDepartmentCommand) herrors.Herr {
	if validate := cmd.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Command validation error: %s", validate)
		return validate
	}
	// 1. 检查部门编码是否已存在
	exists, err := h.deptRepo.GetByCode(ctx, cmd.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "check department code exists failed: %s", err)
		return herrors.CreateFail(err)
	}
	if exists != nil {
		return herrors.CreateFail(fmt.Errorf("department code %s already exists", cmd.Code))
	}

	// 2. 创建部门实体
	dept := model.NewDepartment(cmd.Code, cmd.Name, cmd.Sort)
	dept.ParentID = cmd.ParentID
	dept.Leader = cmd.Leader
	dept.Phone = cmd.Phone
	dept.Email = cmd.Email
	dept.Status = cmd.Status
	dept.Description = cmd.Description
	dept.TenantID = cmd.TenantID

	// 3. 保存部门
	if err := h.deptRepo.Create(ctx, dept); err != nil {
		hlog.CtxErrorf(ctx, "failed to create department: %s", err)
		return herrors.CreateFail(err)
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
	dept, err := h.deptRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to find department: %s", err)
		return herrors.UpdateFail(err)
	}
	if dept == nil {
		return herrors.UpdateFail(fmt.Errorf("department not found: %s", cmd.ID))
	}

	// 2. 检查编码是否重复
	if dept.Code != cmd.Code {
		exists, err := h.deptRepo.GetByCode(ctx, cmd.Code)
		if err != nil {
			hlog.CtxErrorf(ctx, "check department code exists failed: %s", err)
			return herrors.UpdateFail(err)
		}
		if exists != nil {
			return herrors.UpdateFail(fmt.Errorf("department code %s already exists", cmd.Code))
		}
	}

	// 3. 更新部门信息
	dept.UpdateBasicInfo(cmd.Name, cmd.Code, cmd.Sort)
	dept.UpdateContactInfo(cmd.Leader, cmd.Phone, cmd.Email)
	dept.UpdateStatus(cmd.Status)
	dept.UpdateParent(cmd.ParentID)
	dept.Description = cmd.Description

	// 4. 保存更新
	if err := h.deptRepo.Update(ctx, dept); err != nil {
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
	// 1. 检查部门是否存在
	dept, err := h.deptRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to find department: %s", err)
		return herrors.DeleteFail(err)
	}
	if dept == nil {
		return herrors.DeleteFail(fmt.Errorf("department not found: %s", cmd.ID))
	}

	// 2. 检查是否有子部门
	children, err := h.deptRepo.GetByParentID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get sub departments: %s", err)
		return herrors.DeleteFail(err)
	}
	if len(children) > 0 {
		return herrors.DeleteFail(fmt.Errorf("department has sub departments"))
	}

	// 3. 删除部门
	if err := h.deptRepo.Delete(ctx, cmd.ID); err != nil {
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
