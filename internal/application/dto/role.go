package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
)

// RoleDto 角色数据传输对象
type RoleDto struct {
	ID          int64   `json:"id"`          // ID
	Code        string  `json:"code"`        // 编码
	Name        string  `json:"name"`        // 名称
	Description string  `json:"description"` // 描述
	Status      int8    `json:"status"`      // 状态
	PermIds     []int64 `json:"permIds"`     // 权限ID列表
	CreatedAt   int64   `json:"createdAt"`   // 创建时间
	UpdatedAt   int64   `json:"updatedAt"`   // 更新时间
}

// ToRoleDto 领域模型转换为DTO
func ToRoleDto(r *model.Role) *RoleDto {
	if r == nil {
		return nil
	}

	permIds := make([]int64, 0, len(r.Permissions))
	for _, perm := range r.Permissions {
		permIds = append(permIds, perm.ID)
	}

	return &RoleDto{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		PermIds:     permIds,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// ToRoleDtoList 领域���型列表转换为DTO列表
func ToRoleDtoList(roles []*model.Role) []*RoleDto {
	if roles == nil {
		return nil
	}

	dtos := make([]*RoleDto, 0, len(roles))
	for _, r := range roles {
		if dto := ToRoleDto(r); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
