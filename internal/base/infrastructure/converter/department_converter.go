package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type DepartmentConverter struct{}

func NewDepartmentConverter() *DepartmentConverter {
	return &DepartmentConverter{}
}

// ToDTO 将实体转换为DTO
func (c *DepartmentConverter) ToDTO(dept *entity.Department) *dto.DepartmentDto {
	if dept == nil {
		return nil
	}
	return &dto.DepartmentDto{
		ID:          dept.ID,
		ParentID:    dept.ParentID,
		Name:        dept.Name,
		Code:        dept.Code,
		Sequence:    dept.Sequence,
		AdminID:     dept.AdminID,
		Leader:      dept.Leader,
		Phone:       dept.Phone,
		Email:       dept.Email,
		Status:      dept.Status,
		Description: dept.Description,
	}
}

// ToDTOList 将实体列表转换为DTO列表
func (c *DepartmentConverter) ToDTOList(depts []*entity.Department) []*dto.DepartmentDto {
	dtos := make([]*dto.DepartmentDto, 0, len(depts))
	for _, dept := range depts {
		if dto := c.ToDTO(dept); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToTreeDTO 将实体转换为树形DTO
func (c *DepartmentConverter) ToTreeDTO(dept *entity.Department) *dto.DepartmentTreeDto {
	if dept == nil {
		return nil
	}
	return &dto.DepartmentTreeDto{
		ID:       dept.ID,
		ParentID: dept.ParentID,
		Name:     dept.Name,
	}
}
