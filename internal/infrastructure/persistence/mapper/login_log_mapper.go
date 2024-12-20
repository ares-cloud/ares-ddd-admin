package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

type LoginLogMapper struct{}

// ToEntity 领域模型转换为实体
func (m *LoginLogMapper) ToEntity(domain *model.LoginLog) *entity.LoginLog {
	if domain == nil {
		return nil
	}
	return &entity.LoginLog{
		ID:        domain.ID,
		UserID:    domain.UserID,
		Username:  domain.Username,
		TenantID:  domain.TenantID,
		IP:        domain.IP,
		Location:  domain.Location,
		Device:    domain.Device,
		OS:        domain.OS,
		Browser:   domain.Browser,
		Status:    domain.Status,
		Message:   domain.Message,
		LoginTime: domain.LoginTime,
		BaseIntTime: database.BaseIntTime{
			CreatedAt: domain.CreatedAt,
			UpdatedAt: domain.UpdatedAt,
		},
	}
}

// ToDomain 实体转换为领域模型
func (m *LoginLogMapper) ToDomain(entity *entity.LoginLog) *model.LoginLog {
	if entity == nil {
		return nil
	}
	return &model.LoginLog{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Username:  entity.Username,
		TenantID:  entity.TenantID,
		IP:        entity.IP,
		Location:  entity.Location,
		Device:    entity.Device,
		OS:        entity.OS,
		Browser:   entity.Browser,
		Status:    entity.Status,
		Message:   entity.Message,
		LoginTime: entity.LoginTime,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToEntityList 领域模型列表转换为实体列表
func (m *LoginLogMapper) ToEntityList(domains []*model.LoginLog) []*entity.LoginLog {
	if domains == nil {
		return nil
	}
	entities := make([]*entity.LoginLog, 0, len(domains))
	for _, domain := range domains {
		if entity := m.ToEntity(domain); entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
}

// ToDomainList 实体列表转换为领域模型列表
func (m *LoginLogMapper) ToDomainList(entities []*entity.LoginLog) []*model.LoginLog {
	if entities == nil {
		return nil
	}
	domains := make([]*model.LoginLog, 0, len(entities))
	for _, entity := range entities {
		if domain := m.ToDomain(entity); domain != nil {
			domains = append(domains, domain)
		}
	}
	return domains
}
