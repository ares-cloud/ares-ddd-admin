package service

import (
	"context"
	"strconv"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
)

type DataPermissionService struct {
	permRepo repository.IDataPermissionRepository
	deptRepo repository.IDepartmentRepository
	roleRepo repository.IRoleRepository
}

func NewDataPermissionService(
	permRepo repository.IDataPermissionRepository,
	deptRepo repository.IDepartmentRepository,
	roleRepo repository.IRoleRepository,
) *DataPermissionService {
	return &DataPermissionService{
		permRepo: permRepo,
		deptRepo: deptRepo,
		roleRepo: roleRepo,
	}
}

// GetDataRoles 获取数据权限角色列表
func (s *DataPermissionService) GetDataRoles(ctx context.Context) ([]*model.Role, error) {
	// 查询数据权限类型的角色
	return s.roleRepo.FindByType(ctx, int8(model.RoleTypeData))
}

// GetResourceRoles 获取资源角色列表
func (s *DataPermissionService) GetResourceRoles(ctx context.Context) ([]*model.Role, error) {
	// 查询资源类型的角色
	return s.roleRepo.FindByType(ctx, int8(model.RoleTypeResource))
}

// GetDataScope 获取用户的数据权限范围
func (s *DataPermissionService) GetDataScope(ctx context.Context, userID string) (*model.DataPermission, error) {
	// 1. 获取用户角色
	roleIDs := actx.GetRoles(ctx)
	if len(roleIDs) == 0 {
		// 没有角色时返回默认的部门及下级数据权限
		return &model.DataPermission{
			Scope: model.DataScopeDeptTree,
		}, nil
	}

	// 2. 只获取数据权限类型角色的权限
	dataRoleIDs := make([]int64, 0)
	for _, roleIDStr := range roleIDs {
		roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
		if err != nil {
			continue
		}
		if roleModel, err := s.roleRepo.FindByID(ctx, roleID); err == nil && roleModel.IsDataRole() {
			dataRoleIDs = append(dataRoleIDs, roleID)
		}
	}

	// 3. 获取角色的数据权限
	perms, err := s.permRepo.GetByRoleIDs(ctx, dataRoleIDs)
	if err != nil {
		return nil, err
	}

	// 4. 如果没有配置数据权限,使用默认的部门及下级数据权限
	if len(perms) == 0 {
		return &model.DataPermission{
			Scope: model.DataScopeDeptTree,
		}, nil
	}

	// 5. 取最大权限范围
	var maxScope model.DataScope
	var customDeptIDs []string
	for _, perm := range perms {
		if perm.Scope > maxScope {
			maxScope = perm.Scope
			customDeptIDs = perm.DeptIDs
		}
	}

	return &model.DataPermission{
		Scope:   maxScope,
		DeptIDs: customDeptIDs,
	}, nil
}

// GetAccessibleDeptIDs 获取可访问的部门ID列表
func (s *DataPermissionService) GetAccessibleDeptIDs(ctx context.Context, userID string) ([]string, error) {
	// 1. 获取数据权限范围
	scope, err := s.GetDataScope(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 根据权限范围获取部门ID列表
	switch scope.Scope {
	case model.DataScopeAll:
		return nil, nil // 全部数据返回nil,表示不限制
	case model.DataScopeSelf:
		return nil, nil //仅查询自己的数据可以不要部门的限制
	case model.DataScopeDept:
		return []string{actx.GetDeptId(ctx)}, nil
	case model.DataScopeDeptTree:
		return s.getDeptTreeIDs(ctx, actx.GetDeptId(ctx))
	case model.DataScopeCustom:
		return scope.DeptIDs, nil
	default:
		return []string{}, nil
	}
}

// getDeptTreeIDs 获取部门树下的所有部门ID
func (s *DataPermissionService) getDeptTreeIDs(ctx context.Context, deptID string) ([]string, error) {
	ids, err := s.GetChildrenRecursively(ctx, deptID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetChildrenRecursively 递归获取所有子部门
func (s *DataPermissionService) GetChildrenRecursively(ctx context.Context, deptID string) ([]string, error) {
	// 1. 获取直接子部门
	children, err := s.deptRepo.GetByParentID(ctx, deptID)
	if err != nil {
		return nil, err
	}

	// 2. 初始化结果切片,包含当前部门ID
	result := make([]string, 0)
	result = append(result, deptID)

	// 3. 递归获取每个子部门的子部门
	for _, child := range children {
		// 获取子部门的所有下级部门ID
		childrenIDs, err := s.GetChildrenRecursively(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		// 将子部门ID添加到结果中
		result = append(result, childrenIDs...)
	}

	return result, nil
}
