package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type RoleMapper struct{}

func (m *RoleMapper) ToDomain(e *entity.Role, permissions []*model.Permissions) *model.Role {
	return &model.Role{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Localize:    e.Localize,
		Description: e.Description,
		Sequence:    e.Sequence,
		Status:      e.Status,
		Type:        model.RoleType(e.Type),
		Permissions: permissions,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (m *RoleMapper) ToEntity(d *model.Role) *entity.Role {
	return &entity.Role{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Localize:    d.Localize,
		Description: d.Description,
		Sequence:    d.Sequence,
		Type:        int8(d.Type),
		Status:      d.Status,
	}
}

func (m *RoleMapper) ToRolePermissions(roleID int64, permissions []*model.Permissions) []*entity.RolePermissions {
	result := make([]*entity.RolePermissions, len(permissions))
	for i, perm := range permissions {
		result[i] = &entity.RolePermissions{
			RoleID:       roleID,
			PermissionID: perm.ID,
		}
	}
	return result
}
func (m *RoleMapper) ToDomainList(e []*entity.Role) []*model.Role {
	if len(e) == 0 {
		return nil
	}
	list := make([]*model.Role, len(e))
	for i, user := range e {
		list[i] = m.ToDomain(user, make([]*model.Permissions, 0))
	}
	return list
}

func (m *RoleMapper) ToEntityList(d []*model.Role) []*entity.Role {
	if len(d) == 0 {
		return nil
	}
	list := make([]*entity.Role, len(d))
	for i, user := range d {
		list[i] = m.ToEntity(user)
	}
	return list
}
