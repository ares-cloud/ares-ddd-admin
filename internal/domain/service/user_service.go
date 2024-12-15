package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type UserService struct {
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

func NewUserService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
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
