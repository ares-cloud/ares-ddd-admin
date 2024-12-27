package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// RoleDto 角色DTO
type RoleDto struct {
	ID          int64   `json:"id"`          // 角色ID
	Code        string  `json:"code"`        // 角色代码
	Name        string  `json:"name"`        // 角色名称
	Type        int8    `json:"type"`        // 角色类型(1:资源角色 2:数据权限角色)
	Localize    string  `json:"localize"`    // 国际化key
	Description string  `json:"description"` // 描述
	Sequence    int     `json:"sequence"`    // 排序
	Status      int8    `json:"status"`      // 状态
	PermIds     []int64 `json:"permIds"`     //权限id
	TenantID    string  `json:"tenantId"`    // 租户ID
	CreatedAt   int64   `json:"createdAt"`   // 创建时间
	UpdatedAt   int64   `json:"updatedAt"`   // 更新时间
}

// ToRoleDto 领域模型转换为DTO
func ToRoleDto(r *model.Role) *RoleDto {
	if r == nil {
		return nil
	}

	return &RoleDto{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Type:        int8(r.Type),
		Localize:    r.Localize,
		Description: r.Description,
		Sequence:    r.Sequence,
		Status:      r.Status,
		TenantID:    r.TenantID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		PermIds:     make([]int64, 0), // 初始化为空切片，避免返回null
	}
}

// ToRoleDtoList 领域模型列表转换为DTO列表
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
