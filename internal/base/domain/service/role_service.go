package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
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

// GetAllDataPermissionRoles 获取所有数据权限角色
func (s *RoleService) GetAllDataPermissionRoles(ctx context.Context) ([]*model.Role, error) {
	// 构建查询条件
	qb := query.NewQueryBuilder().
		Where("type", query.Eq, int8(entity.RoleTypeData)).
		Where("status", query.Eq, 1) // 只查询启用状态的角色

	// 查询数据权限角色
	roles, err := s.roleRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
