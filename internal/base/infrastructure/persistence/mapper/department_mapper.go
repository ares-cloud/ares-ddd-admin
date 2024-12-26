package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type DepartmentMapper struct{}

func (m *DepartmentMapper) ToDomain(e *entity.Department) *model.Department {
	if e == nil {
		return nil
	}
	return &model.Department{
		ID:          e.ID,
		TenantID:    e.TenantID,
		ParentID:    e.ParentID,
		Code:        e.Code,
		Name:        e.Name,
		Sort:        e.Sort,
		Leader:      e.Leader,
		Phone:       e.Phone,
		Email:       e.Email,
		Status:      e.Status,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (m *DepartmentMapper) ToDomainList(entities []*entity.Department) []*model.Department {
	if len(entities) == 0 {
		return make([]*model.Department, 0)
	}
	list := make([]*model.Department, len(entities))
	for i, e := range entities {
		list[i] = m.ToDomain(e)
	}
	return list
}

func (m *DepartmentMapper) ToEntity(d *model.Department) *entity.Department {
	if d == nil {
		return nil
	}
	return &entity.Department{
		ID:          d.ID,
		TenantID:    d.TenantID,
		ParentID:    d.ParentID,
		Code:        d.Code,
		Name:        d.Name,
		Sort:        d.Sort,
		Leader:      d.Leader,
		Phone:       d.Phone,
		Email:       d.Email,
		Status:      d.Status,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
