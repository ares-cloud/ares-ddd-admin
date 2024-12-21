package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// PermissionsDto 权限数据传输对象
type PermissionsDto struct {
	ID          int64                     `json:"id"`          // ID
	Code        string                    `json:"code"`        // 编码
	Name        string                    `json:"name"`        // 名称
	Localize    string                    `json:"localize"`    // 本地化
	Icon        string                    `json:"icon"`        // 图标
	Description string                    `json:"description"` // 描述
	Sequence    int                       `json:"sequence"`    // 排序
	Type        int8                      `json:"type"`        // 类型
	Path        string                    `json:"path"`        // 路径
	Properties  string                    `json:"properties"`  // 属性
	Status      int8                      `json:"status"`      // 状态
	ParentID    int64                     `json:"parentId"`    // 父级ID
	ParentPath  string                    `json:"parent_path"` // 父级路径
	Resources   []*PermissionsResourceDto `json:"resources"`   // 资源列表
	CreatedAt   int64                     `json:"createdAt"`   // 创建时间
	UpdatedAt   int64                     `json:"updatedAt"`   // 更新时间
}

// PermissionsResourceDto 权限资源数据传输对象
type PermissionsResourceDto struct {
	Method string `json:"method"` // HTTP方法
	Path   string `json:"path"`   // 资源路径
}

// ToDto 领域模型转换为DTO
func ToPermissionsDto(p *model.Permissions) *PermissionsDto {
	if p == nil {
		return nil
	}

	resources := make([]*PermissionsResourceDto, 0)
	if len(p.Resources) > 0 {
		for _, r := range p.Resources {
			resources = append(resources, &PermissionsResourceDto{
				Method: r.Method,
				Path:   r.Path,
			})
		}
	}
	return &PermissionsDto{
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
		Resources:   resources,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToDtoList 领域模型列表转换为DTO列表
func ToPermissionsDtoList(permissions []*model.Permissions) []*PermissionsDto {
	if permissions == nil {
		return nil
	}

	dtos := make([]*PermissionsDto, 0, len(permissions))
	for _, p := range permissions {
		if dto := ToPermissionsDto(p); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
