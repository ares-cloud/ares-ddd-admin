package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type DataPermissionConverter struct{}

func NewDataPermissionConverter() *DataPermissionConverter {
	return &DataPermissionConverter{}
}

// ToDTO 将实体转换为DTO
func (c *DataPermissionConverter) ToDTO(dp *entity.DataPermission, deptIds []string) *dto.DataPermissionDto {
	if dp == nil {
		return nil
	}
	return &dto.DataPermissionDto{
		ID:       dp.ID,
		RoleID:   dp.RoleID,
		Scope:    dp.Scope,
		DeptIDs:  deptIds,
		TenantID: dp.TenantID,
	}
}
