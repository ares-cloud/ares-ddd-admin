package dto

import "github.com/ares-cloud/ares-ddd-admin/internal/domain/model"

type PermissionsTreeResult struct {
	Tree []*PermissionsTreeDto `json:"tree"`
	Ids  []int64               `json:"ids"`
}

// PermissionsTreeDto 简化的权限树数据传输对象
type PermissionsTreeDto struct {
	ID       int64                 `json:"id"`
	Code     string                `json:"code"`
	Name     string                `json:"name"`
	Localize string                `json:"localize"`
	Icon     string                `json:"icon"`
	ParentID int64                 `json:"parentId"`
	Children []*PermissionsTreeDto `json:"children,omitempty"`
}

// ToPermissionsTreeDto 转换为简化的权限树DTO
func ToPermissionsTreeDto(p *model.Permissions) *PermissionsTreeDto {
	if p == nil {
		return nil
	}

	dto := &PermissionsTreeDto{
		ID:       p.ID,
		Code:     p.Code,
		Name:     p.Name,
		Localize: p.Localize,
		Icon:     p.Icon,
		ParentID: p.ParentID,
	}

	if len(p.Children) > 0 {
		dto.Children = make([]*PermissionsTreeDto, 0, len(p.Children))
		for _, child := range p.Children {
			if childDto := ToPermissionsTreeDto(child); childDto != nil {
				dto.Children = append(dto.Children, childDto)
			}
		}
	}

	return dto
}

// ToPermissionsTreeDtoList 批量转换为简化的权限树DTO
func ToPermissionsTreeDtoList(permissions []*model.Permissions) []*PermissionsTreeDto {
	if len(permissions) == 0 {
		return make([]*PermissionsTreeDto, 0)
	}

	dtos := make([]*PermissionsTreeDto, 0, len(permissions))
	for _, p := range permissions {
		if dto := ToPermissionsTreeDto(p); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
