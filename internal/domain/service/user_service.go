package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type UserService struct {
	userRepo repository.IUserRepository
	permRepo repository.IPermissionsRepository
	roleRepo repository.IRoleRepository
}

func NewUserService(userRepo repository.IUserRepository, permRepo repository.IPermissionsRepository, roleRepo repository.IRoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		permRepo: permRepo,
	}
}

// CheckUserPermission 检查用户是否具有特定权限
func (s *UserService) CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range user.Roles {
		role, err := s.roleRepo.FindByID(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, perm := range role.Permissions {
			if perm.Code == permissionCode {
				return true, nil
			}
		}
	}
	return false, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *UserService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionMap := make(map[string]struct{})
	for _, role := range user.Roles {
		role, err := s.roleRepo.FindByID(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, perm := range role.Permissions {
			permissionMap[perm.Code] = struct{}{}
		}
	}

	permissions := make([]string, 0, len(permissionMap))
	for code := range permissionMap {
		permissions = append(permissions, code)
	}
	return permissions, nil
}

// GetUserMenus 获取用户菜单
// 根据用户绑定的角色获取类型为1(页面)的权限
func (s *UserService) GetUserMenus(ctx context.Context, userID string) ([]*model.Permissions, error) {
	// 获取用户信息及关联的角色
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果用户没有角色，返回空菜单
	if len(user.Roles) == 0 {
		return []*model.Permissions{}, nil
	}

	// 获取用户所有角色的权限
	var allPermissions []*model.Permissions
	for _, role := range user.Roles {
		permissions, err := s.permRepo.FindByRoleID(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		allPermissions = append(allPermissions, permissions...)
	}

	// 过滤出类型为1(页面)的权限
	menuPermissions := make([]*model.Permissions, 0)
	permMap := make(map[int64]bool) // 用于去重

	for _, perm := range allPermissions {
		// 只获取类型为1(页面)且状态为启用的权限
		if perm.Type == 1 && perm.Status == 1 {
			if !permMap[perm.ID] {
				menuPermissions = append(menuPermissions, perm)
				permMap[perm.ID] = true
			}
		}
	}

	// 构建菜单树
	return buildPermissionTree(menuPermissions), nil
}

// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []*model.Permissions) []*model.Permissions {
	// 创建ID到权限的映射
	permMap := make(map[int64]*model.Permissions)
	for _, p := range permissions {
		permMap[p.ID] = p
		p.Children = make([]*model.Permissions, 0) // 初始化子节点列表
	}

	// 构建树结构
	var roots []*model.Permissions
	for _, p := range permissions {
		if p.ParentID == 0 {
			roots = append(roots, p)
		} else {
			if parent, ok := permMap[p.ParentID]; ok {
				parent.Children = append(parent.Children, p)
			}
		}
	}

	return roots
}
