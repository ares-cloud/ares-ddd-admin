package converter

import (
	"sort"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type PermissionsConverter struct{}

func NewPermissionsConverter() *PermissionsConverter {
	return &PermissionsConverter{}
}

// ToDTO 将实体转换为DTO
func (c *PermissionsConverter) ToDTO(p *entity.Permissions, es []*entity.PermissionsResource) *dto.PermissionsDto {
	if p == nil {
		return nil
	}

	resources := make([]*dto.PermissionsResourceDto, 0)
	if len(es) > 0 {
		for _, r := range es {
			resources = append(resources, &dto.PermissionsResourceDto{
				Method: r.Method,
				Path:   r.Path,
			})
		}
	}

	return &dto.PermissionsDto{
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

// ToDTOList 将实体列表转换为DTO列表
func (c *PermissionsConverter) ToDTOList(permissions []*entity.Permissions) []*dto.PermissionsDto {
	dtos := make([]*dto.PermissionsDto, 0, len(permissions))
	for _, p := range permissions {
		if dto := c.ToDTO(p, nil); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToTreeDTOList 将权限列表转换为树形结构DTO
func (c *PermissionsConverter) ToTreeDTOList(permissions []*entity.Permissions) []*dto.PermissionsTreeDto {
	// 1. 转换为DTO
	dtos := make([]*dto.PermissionsTreeDto, 0, len(permissions))
	for _, p := range permissions {
		if dto := c.toTreeDTO(p); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	// 2. 构建树形结构
	return buildPermissionTree(dtos)
}

// toTreeDTO 将单个权限实体转换为树形DTO
func (c *PermissionsConverter) toTreeDTO(p *entity.Permissions) *dto.PermissionsTreeDto {
	if p == nil {
		return nil
	}

	//resourceDtos := make([]*dto.PermissionsResourceDto, 0)
	//if len(resources) > 0 {
	//	for _, r := range resources {
	//		resourceDtos = append(resourceDtos, &dto.PermissionsResourceDto{
	//			Method: r.Method,
	//			Path:   r.Path,
	//		})
	//	}
	//}

	return &dto.PermissionsTreeDto{
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
		Children:    make([]*dto.PermissionsTreeDto, 0),
	}
}

// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []*dto.PermissionsTreeDto) []*dto.PermissionsTreeDto {
	// 1. 构建ID到权限的映射
	permMap := make(map[int64]*dto.PermissionsTreeDto)
	for _, p := range permissions {
		permMap[p.ID] = p
	}

	// 2. 构建树形结构
	var roots []*dto.PermissionsTreeDto
	for _, p := range permissions {
		if p.ParentID == 0 {
			roots = append(roots, p)
		} else {
			if parent, ok := permMap[p.ParentID]; ok {
				parent.Children = append(parent.Children, p)
			}
		}
	}

	// 3. 对每个节点的子节点进行排序
	for _, p := range permissions {
		if len(p.Children) > 0 {
			sortPermissions(p.Children)
		}
	}

	// 4. 对根节点进行排序
	sortPermissions(roots)

	return roots
}

// sortPermissions 根据序号对权限列表进行排序
func sortPermissions(permissions []*dto.PermissionsTreeDto) {
	sort.Slice(permissions, func(i, j int) bool {
		if permissions[i].Sequence == permissions[j].Sequence {
			return permissions[i].ID < permissions[j].ID
		}
		return permissions[i].Sequence < permissions[j].Sequence
	})
}

// getResourcesByPermissionID 获取指定权限ID的资源列表
func getResourcesByPermissionID(resources []*entity.PermissionsResource, permissionID int64) []*entity.PermissionsResource {
	result := make([]*entity.PermissionsResource, 0)
	for _, r := range resources {
		if r.PermissionsID == permissionID {
			result = append(result, r)
		}
	}
	return result
}
