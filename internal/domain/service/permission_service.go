package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type PermissionService struct {
	permissionRepo repository.IPermissionsRepository
}

func NewPermissionService(permissionRepo repository.IPermissionsRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

// ValidatePermissionCode 验证权限代码是否有效
func (s *PermissionService) ValidatePermissionCode(ctx context.Context, code string) (bool, error) {
	return s.permissionRepo.ExistsByCode(ctx, code)
}

// GetPermissionsByType 根据类型获取权限列表
func (s *PermissionService) GetPermissionsByType(ctx context.Context, permType int8) ([]*model.Permissions, error) {
	return s.permissionRepo.FindByType(ctx, permType)
}
