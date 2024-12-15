package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
)

type PermissionsMapper struct{}

// ToDomain 实体转换为领域模型
func (m *PermissionsMapper) ToDomain(e *entity.Permissions, resources []*entity.PermissionsResource) *model.Permissions {
	if e == nil {
		return nil
	}

	perm := &model.Permissions{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Localize:    e.Localize,
		Icon:        e.Icon,
		Description: e.Description,
		Sequence:    e.Sequence,
		Type:        e.Type,
		Path:        e.Path,
		Properties:  e.Properties,
		Status:      e.Status,
		ParentID:    e.ParentID,
		ParentPath:  e.ParentPath,
		Resources:   make([]*model.PermissionsResource, 0),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	// 转换资源列表
	if len(resources) > 0 {
		for _, r := range resources {
			perm.Resources = append(perm.Resources, &model.PermissionsResource{
				ID:            r.ID,
				PermissionsID: r.PermissionsID,
				Method:        r.Method,
				Path:          r.Path,
			})
		}
	}

	return perm
}

// ToEntity 领域模型转换为实体
func (m *PermissionsMapper) ToEntity(d *model.Permissions) (*entity.Permissions, []*entity.PermissionsResource) {
	if d == nil {
		return nil, nil
	}

	permEntity := &entity.Permissions{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Localize:    d.Localize,
		Icon:        d.Icon,
		Description: d.Description,
		Sequence:    d.Sequence,
		Type:        d.Type,
		Path:        d.Path,
		Properties:  d.Properties,
		Status:      d.Status,
		ParentID:    d.ParentID,
		ParentPath:  d.ParentPath,
	}

	// 转换资源列表
	var resources []*entity.PermissionsResource
	if len(d.Resources) > 0 {
		resources = make([]*entity.PermissionsResource, len(d.Resources))
		for i, r := range d.Resources {
			resources[i] = &entity.PermissionsResource{
				PermissionsID: d.ID,
				Method:        r.Method,
				Path:          r.Path,
			}
		}
	}

	return permEntity, resources
}

// ToDomainList 实体列表转换为领域模型列表
func (m *PermissionsMapper) ToDomainList(e []*entity.Permissions, r []*entity.PermissionsResource) []*model.Permissions {
	if len(e) == 0 {
		return make([]*model.Permissions, 0)
	}
	list := make([]*model.Permissions, len(e))
	resourceMap := m.GroupResourcesByPermissionID(r)
	for i, v := range e {
		resources := resourceMap[v.ID]
		list[i] = m.ToDomain(v, resources)
	}
	return list
}

// GroupResourcesByPermissionID 将权限资源按权限ID分组
func (m *PermissionsMapper) GroupResourcesByPermissionID(resources []*entity.PermissionsResource) map[int64][]*entity.PermissionsResource {
	resourceMap := make(map[int64][]*entity.PermissionsResource)
	for _, res := range resources {
		resourceMap[res.PermissionsID] = append(resourceMap[res.PermissionsID], res)
	}
	return resourceMap
}
