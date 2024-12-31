package impl

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type DepartmentQueryService struct {
	deptRepo      repository.ISysDepartmentRepo
	userRepo      repository.ISysUserRepo
	deptConverter *converter.DepartmentConverter
	userConverter *converter.UserConverter
}

func NewDepartmentQueryService(
	deptRepo repository.ISysDepartmentRepo,
	userRepo repository.ISysUserRepo,
	deptConverter *converter.DepartmentConverter,
	userConverter *converter.UserConverter,
) *DepartmentQueryService {
	return &DepartmentQueryService{
		deptRepo:      deptRepo,
		userRepo:      userRepo,
		deptConverter: deptConverter,
		userConverter: userConverter,
	}
}

// GetDepartment 获取部门详情
func (d *DepartmentQueryService) GetDepartment(ctx context.Context, id string) (*dto.DepartmentDto, error) {
	dept, err := d.deptRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return d.deptConverter.ToDTO(dept), nil
}

// FindDepartments 查询部门列表
func (d *DepartmentQueryService) FindDepartments(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.DepartmentDto, error) {
	depts, err := d.deptRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}
	return d.deptConverter.ToDTOList(depts), nil
}

// CountDepartments 统计部门数量
func (d *DepartmentQueryService) CountDepartments(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return d.deptRepo.Count(ctx, qb)
}

// GetDepartmentTree 获取部门树
func (d *DepartmentQueryService) GetDepartmentTree(ctx context.Context, parentID string) ([]*dto.DepartmentTreeDto, error) {
	// 1. 获取直接子部门
	children, err := d.deptRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 2. 构建部门树
	tree := make([]*dto.DepartmentTreeDto, len(children))
	// 复制一份数据，避免修改原始实体
	for i, child := range children {
		tree[i] = d.deptConverter.ToTreeDTO(child)
	}

	for _, child := range tree {
		// 递归获取子部门
		subChildren, err := d.getChildrenRecursively(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		// 添加子部门
		child.Children = subChildren
	}

	return tree, nil
}

// getChildrenRecursively 递归获取子部门
func (d *DepartmentQueryService) getChildrenRecursively(ctx context.Context, parentID string) ([]*dto.DepartmentTreeDto, error) {
	children, err := d.deptRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 复制一份数据，避免修改原始实体
	result := make([]*dto.DepartmentTreeDto, len(children))
	for i, child := range children {
		result[i] = d.deptConverter.ToTreeDTO(child)
	}

	for _, child := range result {
		subChildren, err := d.getChildrenRecursively(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		child.Children = subChildren
	}

	return result, nil
}

// GetUserDepartments 获取用户部门
func (d *DepartmentQueryService) GetUserDepartments(ctx context.Context, userID string) ([]*dto.DepartmentDto, error) {
	// 1. 获取用户部门关联
	userDepts, err := d.deptRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(userDepts) == 0 {
		return []*dto.DepartmentDto{}, nil
	}

	// 2. 获取部门ID列表
	deptIDs := make([]string, 0, len(userDepts))
	for _, ud := range userDepts {
		deptIDs = append(deptIDs, ud.DeptID)
	}

	// 3. 查询部门信息
	depts, err := d.deptRepo.FindByIds(ctx, deptIDs)
	if err != nil {
		return nil, err
	}

	return d.deptConverter.ToDTOList(depts), nil
}

// GetDepartmentUsers 获取部门用户
func (d *DepartmentQueryService) GetDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	users, err := d.userRepo.FindByDepartment(ctx, deptID, excludeAdminID, qb)
	if err != nil {
		return nil, err
	}
	return d.userConverter.ToDTOList(users), nil
}

// CountDepartmentUsers 统计部门用户数量
func (d *DepartmentQueryService) CountDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	return d.userRepo.CountByDepartment(ctx, deptID, excludeAdminID, qb)
}

// GetUnassignedUsers 获取未分配部门的用户
func (d *DepartmentQueryService) GetUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	users, err := d.userRepo.FindUnassignedUsers(ctx, qb)
	if err != nil {
		return nil, err
	}
	return d.userConverter.ToDTOList(users), nil
}

// CountUnassignedUsers 统计未分配部门的用户数量
func (d *DepartmentQueryService) CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return d.userRepo.CountUnassignedUsers(ctx, qb)
}
