package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"strings"
)

type DataPermissionConverter struct{}

func NewDataPermissionConverter() *DataPermissionConverter {
	return &DataPermissionConverter{}
}

func (c *DataPermissionConverter) ToDTO(perm *entity.DataPermission) *dto.DataPermissionDto {
	if perm == nil {
		return nil
	}
	deptIds := make([]string, 0)
	if perm.DeptIDs != "" {
		deptIds = strings.Split(perm.DeptIDs, ",")
	}
	return &dto.DataPermissionDto{
		ID:      perm.ID,
		RoleID:  perm.RoleID,
		Scope:   perm.Scope,
		DeptIDs: deptIds,
	}
}
