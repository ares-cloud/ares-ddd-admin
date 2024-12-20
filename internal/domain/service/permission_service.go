package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
)

type PermissionService struct {
	permissionRepo repository.IPermissionsRepository
	tenantRepo     repository.ITenantRepository
}

func NewPermissionService(permissionRepo repository.IPermissionsRepository, tenantRepo repository.ITenantRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		tenantRepo:     tenantRepo,
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

func (s *PermissionService) FindAllEnabled(ctx context.Context) ([]*model.Permissions, error) {
	tenantId := actx.GetTenantId(ctx)
	if tenantId != "" {
		return s.tenantRepo.GetPermissions(ctx, tenantId)
	}
	return s.permissionRepo.FindAllEnabled(ctx)
}
