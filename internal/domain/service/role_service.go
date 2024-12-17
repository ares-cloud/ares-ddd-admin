package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type RoleService struct {
	roleRepo repository.IRoleRepository
}

func NewRoleService(roleRepo repository.IRoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

// CheckRolePermissions 检查角色是否具有特定权限
func (s *RoleService) CheckRolePermissions(ctx context.Context, roleID int64, permissionCode string) (bool, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return false, err
	}

	for _, perm := range role.Permissions {
		if perm.Code == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

// GetRolePermissionCodes 获取角色的所有权限代码
func (s *RoleService) GetRolePermissionCodes(ctx context.Context, roleID int64) ([]string, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(role.Permissions))
	for i, perm := range role.Permissions {
		codes[i] = perm.Code
	}

	return codes, nil
}
