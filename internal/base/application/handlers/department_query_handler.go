package handlers

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type DepartmentQueryHandler struct {
	deptRepo    repository.IDepartmentRepository
	userRepo    repository.IUserRepository
	deptService *service.DepartmentService
}

func NewDepartmentQueryHandler(
	deptRepo repository.IDepartmentRepository,
	userRepo repository.IUserRepository,
	deptService *service.DepartmentService,
) *DepartmentQueryHandler {
	return &DepartmentQueryHandler{
		deptRepo:    deptRepo,
		userRepo:    userRepo,
		deptService: deptService,
	}
}

// HandleList 处理部门列表查询
func (h *DepartmentQueryHandler) HandleList(ctx context.Context, req *queries.ListDepartmentsQuery) (*models.PageRes[dto.DepartmentDto], herrors.Herr) {
	if validate := req.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 构建查询条件
	qb := query.NewQueryBuilder()
	if req.Name != "" {
		qb.Where("name", query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		qb.Where("code", query.Like, "%"+req.Code+"%")
	}
	if req.Status != nil {
		qb.Where("status", query.Eq, *req.Status)
	}
	if req.ParentID != "" {
		qb.Where("parent_id", query.Eq, req.ParentID)
	}
	qb.OrderBy("sort", false)
	qb.WithPage(&req.Page)

	// 查询总数
	total, err := h.deptRepo.Count(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to count departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 查询列表数据
	depts, err := h.deptRepo.Find(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to list departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO并返回分页结果
	return &models.PageRes[dto.DepartmentDto]{
		Total: total,
		List:  dto.ToDepartmentDtoList(depts),
	}, nil
}

// HandleGet 处理获取部门查询
func (h *DepartmentQueryHandler) HandleGet(ctx context.Context, query *queries.GetDepartmentQuery) (*dto.DepartmentDto, herrors.Herr) {
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

	// 转换为DTO
	return dto.ToDepartmentDto(dept), nil
}

// HandleGetTree 处理获取部门树查询
func (h *DepartmentQueryHandler) HandleGetTree(ctx context.Context, query *queries.GetDepartmentTreeQuery) ([]*dto.DepartmentTreeDto, herrors.Herr) {
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

	// 转换为DTO
	return dto.ToDepartmentTreeDtoList(tree), nil
}

// HandleGetUserDepartments 处理获取用户部门查询
func (h *DepartmentQueryHandler) HandleGetUserDepartments(ctx context.Context, query *queries.GetUserDepartmentsQuery) ([]*dto.DepartmentDto, herrors.Herr) {
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

	// 转换为DTO
	return dto.ToDepartmentDtoList(depts), nil
}

// HandleGetUsers 处理获取部门用户
func (h *DepartmentQueryHandler) HandleGetUsers(ctx context.Context, req *queries.GetDepartmentUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 1. 查询部门信息
	_, err := h.deptRepo.GetByID(ctx, req.DeptID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	// 2. 构建查询条件
	qb := query.NewQueryBuilder()
	if req.Username != "" {
		qb.Where("username", query.Like, "%"+req.Username+"%")
	}
	if req.Name != "" {
		qb.Where("name", query.Like, "%"+req.Name+"%")
	}
	qb.WithPage(&req.Page)
	qb.OrderBy("created_at", true)

	//// 3. 查询总数
	//total, err := h.userRepo.CountByDepartment(ctx, req.DeptID, dept.AdminID, qb)
	//if err != nil {
	//	return nil, herrors.NewServerHError(err)
	//}
	//
	//// 4. 查询用户列表
	//users, err := h.userRepo.FindByDepartment(ctx, req.DeptID, dept.AdminID, qb)
	//if err != nil {
	//	return nil, herrors.NewServerHError(err)
	//}

	//// 5. 转换为DTO并返回分页结果
	//return &models.PageRes[dto.UserDto]{
	//	List:  dto.ToUserDtoList(users),
	//	Total: total,
	//}, nil
	return nil, nil
}

// HandleGetUnassignedUsers 处理获取未分配部门的用户查询
func (h *DepartmentQueryHandler) HandleGetUnassignedUsers(ctx context.Context, req *queries.GetUnassignedUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if req.Username != "" {
		qb.Where("username", query.Like, "%"+req.Username+"%")
	}
	if req.Name != "" {
		qb.Where("name", query.Like, "%"+req.Name+"%")
	}
	// 只查询启用状态的用户
	qb.Where("status", query.Eq, 1)
	// 添加分页
	qb.WithPage(&req.Page)
	qb.OrderBy("created_at", true)

	//// 查询总数
	//total, err := h.userRepo.CountUnassignedUsers(ctx, qb)
	//if err != nil {
	//	return nil, herrors.NewServerHError(err)
	//}
	//
	//// 查询用户列表
	//users, err := h.userRepo.FindUnassignedUsers(ctx, qb)
	//if err != nil {
	//	return nil, herrors.NewServerHError(err)
	//}

	// 转换为DTO并返回分页结果
	//return &models.PageRes[dto.UserDto]{
	//	List:  dto.ToUserDtoList(users),
	//	Total: total,
	//}, nil
	return nil, nil
}
