package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

type OperationLogMapper struct{}

// ToEntity 领域模型转换为实体
func (m *OperationLogMapper) ToEntity(domain *model.OperationLog) *entity.OperationLog {
	if domain == nil {
		return nil
	}
	return &entity.OperationLog{
		ID:        domain.ID,
		UserID:    domain.UserID,
		Username:  domain.Username,
		TenantID:  domain.TenantID,
		Method:    domain.Method,
		Path:      domain.Path,
		Query:     domain.Query,
		Body:      domain.Body,
		IP:        domain.IP,
		UserAgent: domain.UserAgent,
		Status:    domain.Status,
		Error:     domain.Error,
		Duration:  domain.Duration,
		Module:    domain.Module,
		Action:    domain.Action,
		BaseIntTime: database.BaseIntTime{
			CreatedAt: domain.CreatedAt,
		},
	}
}

// ToDomain 实体转换为领域模型
func (m *OperationLogMapper) ToDomain(entity *entity.OperationLog) *model.OperationLog {
	if entity == nil {
		return nil
	}
	return &model.OperationLog{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Username:  entity.Username,
		TenantID:  entity.TenantID,
		Method:    entity.Method,
		Path:      entity.Path,
		Query:     entity.Query,
		Body:      entity.Body,
		IP:        entity.IP,
		UserAgent: entity.UserAgent,
		Status:    entity.Status,
		Error:     entity.Error,
		Duration:  entity.Duration,
		Module:    entity.Module,
		Action:    entity.Action,
		CreatedAt: entity.CreatedAt,
	}
}

// ToEntityList 领域模型列表转换为实体列表
func (m *OperationLogMapper) ToEntityList(domains []*model.OperationLog) []*entity.OperationLog {
	if domains == nil {
		return nil
	}
	entities := make([]*entity.OperationLog, 0, len(domains))
	for _, domain := range domains {
		if entity := m.ToEntity(domain); entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
}

// ToDomainList 实体列表转换为领域模型列表
func (m *OperationLogMapper) ToDomainList(entities []*entity.OperationLog) []*model.OperationLog {
	if entities == nil {
		return nil
	}
	domains := make([]*model.OperationLog, 0, len(entities))
	for _, entity := range entities {
		if domain := m.ToDomain(entity); domain != nil {
			domains = append(domains, domain)
		}
	}
	return domains
}
