package handlers

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type DepartmentQueryHandler struct {
	deptRepo    repository.IDepartmentRepository
	deptService *service.DepartmentService
}

func NewDepartmentQueryHandler(deptRepo repository.IDepartmentRepository, deptService *service.DepartmentService) *DepartmentQueryHandler {
	return &DepartmentQueryHandler{
		deptRepo:    deptRepo,
		deptService: deptService,
	}
}

// HandleList 处理部门列表查询
func (h *DepartmentQueryHandler) HandleList(ctx context.Context, query *queries.ListDepartmentsQuery) ([]*model.Department, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 构建查询条件
	listQuery := &repository.ListDepartmentQuery{
		Name:   query.Name,
		Code:   query.Code,
		Status: query.Status,
	}

	// 查询部门列表
	depts, err := h.deptRepo.List(ctx, listQuery)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to list departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return depts, nil
}

// HandleGet 处理获取部门查询
func (h *DepartmentQueryHandler) HandleGet(ctx context.Context, query *queries.GetDepartmentQuery) (*model.Department, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 查询部门
	dept, err := h.deptRepo.GetByID(ctx, query.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get department: %s", err)
		return nil, herrors.QueryFail(err)
	}
	if dept == nil {
		return nil, herrors.QueryFail(fmt.Errorf("department not found: %s", query.ID))
	}

	return dept, nil
}

// HandleGetTree 处理获取部门树查询
func (h *DepartmentQueryHandler) HandleGetTree(ctx context.Context, query *queries.GetDepartmentTreeQuery) ([]*model.Department, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 获取部门树
	tree, err := h.deptService.GetDepartmentTree(ctx, query.ParentID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get department tree: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return tree, nil
}

// HandleGetUserDepartments 处理获取用户部门查询
func (h *DepartmentQueryHandler) HandleGetUserDepartments(ctx context.Context, query *queries.GetUserDepartmentsQuery) ([]*model.Department, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 获取用户部门
	depts, err := h.deptRepo.GetUserDepartments(ctx, query.UserID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get user departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return depts, nil
}
