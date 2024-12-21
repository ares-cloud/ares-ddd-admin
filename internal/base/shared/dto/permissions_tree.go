package dto

import "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"

type PermissionsTreeResult struct {
	Tree []*PermissionsTreeDto `json:"tree"`
	Ids  []int64               `json:"ids"`
}

// PermissionsTreeDto 简化的权限树数据传输对象
type PermissionsTreeDto struct {
	ID          int64                 `json:"id"`          // ID
	Code        string                `json:"code"`        // 编码
	Name        string                `json:"name"`        // 名称
	Localize    string                `json:"localize"`    // 本地化
	Icon        string                `json:"icon"`        // 图标
	Description string                `json:"description"` // 描述
	Sequence    int                   `json:"sequence"`    // 排序
	Type        int8                  `json:"type"`        // 类型
	Path        string                `json:"path"`        // 路径
	Properties  string                `json:"properties"`  // 属性
	Status      int8                  `json:"status"`      // 状态
	ParentID    int64                 `json:"parentId"`    // 父级ID
	ParentPath  string                `json:"parent_path"` // 父级路径
	Children    []*PermissionsTreeDto `json:"children,omitempty"`
}

// ToPermissionsTreeDto 转换为简化的权限树DTO
func ToPermissionsTreeDto(p *model.Permissions) *PermissionsTreeDto {
	if p == nil {
		return nil
	}

	dto := &PermissionsTreeDto{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Localize:    p.Localize,
		Icon:        p.Icon,
		Description: p.Description,
		Sequence:    p.Sequence,
		Type:        p.Type,
		Path:        p.Path,
		Properties:  p.Properties,
		Status:      p.Status,
		ParentID:    p.ParentID,
		ParentPath:  p.ParentPath,
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
