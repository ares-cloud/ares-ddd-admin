package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type TenantMapper struct {
	userMapper *UserMapper
}

func NewTenantMapper(userMapper *UserMapper) *TenantMapper {
	return &TenantMapper{
		userMapper: userMapper,
	}
}

// ToEntity 领域模型转换为实体
func (m *TenantMapper) ToEntity(domain *model.Tenant) *entity.Tenant {
	if domain == nil {
		return nil
	}

	return &entity.Tenant{
		ID:          domain.ID,
		Code:        domain.Code,
		Name:        domain.Name,
		Domain:      domain.Domain,
		Status:      domain.Status,
		IsDefault:   domain.IsDefault,
		ExpireTime:  domain.ExpireTime,
		Description: domain.Description,
	}
}

// ToDomain 实体转换为领域模型
func (m *TenantMapper) ToDomain(entity *entity.Tenant, adminUser *entity.SysUser) *model.Tenant {
	if entity == nil {
		return nil
	}

	tenant := &model.Tenant{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Domain:      entity.Domain,
		Status:      entity.Status,
		IsDefault:   entity.IsDefault,
		ExpireTime:  entity.ExpireTime,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	// 转换管理员用户
	if adminUser != nil {
		tenant.AdminUser = m.userMapper.ToDomain(adminUser, nil)
	}

	return tenant
}

// ToDomainList 实体列表转换为领域模型列表
func (m *TenantMapper) ToDomainList(entities []*entity.Tenant) []*model.Tenant {
	if len(entities) == 0 {
		return make([]*model.Tenant, 0)
	}

	domains := make([]*model.Tenant, len(entities))
	for i, e := range entities {
		domains[i] = m.ToDomain(e, nil)
	}
	return domains
}

// ToEntityList 领域模型列表转换为实体列表
func (m *TenantMapper) ToEntityList(domains []*model.Tenant) []*entity.Tenant {
	if len(domains) == 0 {
		return make([]*entity.Tenant, 0)
	}

	entities := make([]*entity.Tenant, len(domains))
	for i, d := range domains {
		entities[i] = m.ToEntity(d)
	}
	return entities
}
