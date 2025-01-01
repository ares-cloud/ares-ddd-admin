package impl

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
)

type DataPermissionQueryService struct {
	permissionConverter *converter.DataPermissionConverter
	repo                repository.IDataPermissionRepo
}

func NewDataPermissionQueryService(
	repo repository.IDataPermissionRepo,
	permissionConverter *converter.DataPermissionConverter,
) *DataPermissionQueryService {
	return &DataPermissionQueryService{
		repo:                repo,
		permissionConverter: permissionConverter,
	}
}

// GetByRoleID 获取角色的数据权限
func (s *DataPermissionQueryService) GetByRoleID(ctx context.Context, roleID int64) (*dto.DataPermissionDto, error) {
	perm, err := s.repo.FindByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, nil
	}
	return s.permissionConverter.ToDTO(perm), nil
}
