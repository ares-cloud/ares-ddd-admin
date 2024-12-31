package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type RoleConverter struct{}

func NewRoleConverter() *RoleConverter {
	return &RoleConverter{}
}

// ToDTO 将实体转换为DTO
func (c *RoleConverter) ToDTO(role *entity.Role, permIds []int64) *dto.RoleDto {
	if role == nil {
		return nil
	}
	return &dto.RoleDto{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Type:        role.Type,
		Localize:    role.Localize,
		Description: role.Description,
		Sequence:    role.Sequence,
		Status:      role.Status,
		PermIds:     permIds,
		TenantID:    role.TenantID,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// ToDTOList 将实体列表转换为DTO列表
func (c *RoleConverter) ToDTOList(roles []*entity.Role) []*dto.RoleDto {
	dtos := make([]*dto.RoleDto, 0, len(roles))
	for _, role := range roles {
		if dto := c.ToDTO(role, nil); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
