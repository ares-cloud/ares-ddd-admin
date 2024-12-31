package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type TenantConverter struct {
	userConverter *UserConverter
}

func NewTenantConverter(userConverter *UserConverter) *TenantConverter {
	return &TenantConverter{
		userConverter: userConverter,
	}
}

// ToDTO 将实体转换为DTO
func (c *TenantConverter) ToDTO(t *entity.Tenant, adminUser *entity.SysUser) *dto.TenantDto {
	if t == nil {
		return nil
	}

	dto := &dto.TenantDto{
		ID:          t.ID,
		Code:        t.Code,
		Name:        t.Name,
		Domain:      t.Domain,
		Description: t.Description,
		IsDefault:   t.IsDefault,
		Status:      t.Status,
		ExpireTime:  t.ExpireTime,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	// 转换管理员用户
	if adminUser != nil {
		dto.AdminUser = c.userConverter.ToDTO(adminUser, nil)
	}

	return dto
}

// ToDTOList 将实体列表转换为DTO列表
func (c *TenantConverter) ToDTOList(tenants []*entity.Tenant) []*dto.TenantDto {
	dtos := make([]*dto.TenantDto, 0, len(tenants))
	for _, t := range tenants {
		if dto := c.ToDTO(t, nil); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
