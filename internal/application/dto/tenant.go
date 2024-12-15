package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
)

// TenantDto 租户数据传输对象
type TenantDto struct {
	ID          string   `json:"id"`          // ID
	Code        string   `json:"code"`        // 租户编码
	Name        string   `json:"name"`        // 租户名称
	Description string   `json:"description"` // 描述
	IsDefault   int8     `json:"isDefault"`   // 是否默认租户
	Status      int8     `json:"status"`      // 状态
	AdminUser   *UserDto `json:"adminUser"`   // 管理员用户
	ExpireTime  int64    `json:"expireTime"`
	CreatedAt   int64    `json:"createdAt"` // 创建时间
	UpdatedAt   int64    `json:"updatedAt"` // 更新时间
}

// ToTenantDto 领域模型转换为DTO
func ToTenantDto(t *model.Tenant) *TenantDto {
	if t == nil {
		return nil
	}

	dto := &TenantDto{
		ID:          t.ID,
		Code:        t.Code,
		Name:        t.Name,
		Description: t.Description,
		IsDefault:   t.IsDefault,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		ExpireTime:  t.ExpireTime,
		UpdatedAt:   t.UpdatedAt,
	}

	// 转换管理员用户
	if t.AdminUser != nil {
		dto.AdminUser = ToUserDto(t.AdminUser)
	}

	return dto
}

// ToTenantDtoList 领域模型列表转换为DTO列表
func ToTenantDtoList(tenants []*model.Tenant) []*TenantDto {
	if tenants == nil {
		return nil
	}

	dtos := make([]*TenantDto, 0, len(tenants))
	for _, t := range tenants {
		if dto := ToTenantDto(t); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
