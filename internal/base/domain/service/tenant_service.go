package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type TenantService struct {
	tenantRepo repository.ITenantRepository
	eventBus   *pkgEvent.EventBus
}

func NewTenantService(
	tenantRepo repository.ITenantRepository,
	eventBus *pkgEvent.EventBus,
) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		eventBus:   eventBus,
	}
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	// 检查租户编码是否存在
	exists, err := s.tenantRepo.ExistsByCode(ctx, tenant.Code)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrTenantCodeExists
	}

	// 检查域名是否存在
	if tenant.Domain != "" {
		//exists, err = s.tenantRepo.ExistsByDomain(ctx, tenant.Domain)
		//if err != nil {
		//	return err
		//}
		//if exists {
		//	return errors.ErrTenantDomainExists
		//}
	}

	// 创建租户
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return err
	}

	//// 发布租户创建事件
	//event := &events.TenantEvent{
	//	BaseEvent: events.BaseEvent{TenantID: tenant.ID},
	//	Action:    events.TenantCreated,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	// 检查租户是否存在
	old, err := s.tenantRepo.FindByID(ctx, tenant.ID)
	if err != nil {
		return err
	}
	if old == nil {
		return errors.ErrTenantNotFound
	}

	// 检查域名是否被其他租户使用
	if tenant.Domain != old.Domain {
		//exists, err := s.tenantRepo.ExistsByDomain(ctx, tenant.Domain)
		//if err != nil {
		//	return err
		//}
		//if exists {
		//	return errors.ErrTenantDomainExists
		//}
	}

	// 更新租户
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return err
	}

	// 发布租户更新事件
	//event := &events.TenantEvent{
	//	BaseEvent: events.BaseEvent{TenantID: tenant.ID},
	//	Action:    events.TenantUpdated,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(ctx context.Context, id string) error {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if tenant == nil {
		return errors.ErrTenantNotFound
	}

	// 删除租户
	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 发布租户删除事件
	//event := &events.TenantEvent{
	//	BaseEvent: events.BaseEvent{TenantID: id},
	//	Action:    events.TenantDeleted,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}

// UpdateTenantStatus 更新租户状态
func (s *TenantService) UpdateTenantStatus(ctx context.Context, id string, status int8) error {
	//if err := s.tenantRepo.UpdateStatus(ctx, id, status); err != nil {
	//	return err
	//}

	//// 发布租户状态更新事件
	//event := &events.TenantEvent{
	//	BaseEvent: events.BaseEvent{TenantID: id},
	//	Action:    events.TenantStatusChanged,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}

// GetTenant 获取租户信息
func (s *TenantService) GetTenant(ctx context.Context, id string) (*model.Tenant, error) {
	return s.tenantRepo.FindByID(ctx, id)
}

// GetAllEnabledTenants 获取所有启用的租户
func (s *TenantService) GetAllEnabledTenants(ctx context.Context) ([]*model.Tenant, error) {
	//return s.tenantRepo.GetAllEnabled(ctx)
	return nil, nil
}

// GetDefaultTenant 获取默认租户
func (s *TenantService) GetDefaultTenant(ctx context.Context) (*model.Tenant, error) {
	//return s.tenantRepo.GetDefaultTenant(ctx)
	return nil, nil
}

// FindTenants 查询租户列表
func (s *TenantService) FindTenants(ctx context.Context, qb *query.QueryBuilder) ([]*model.Tenant, error) {
	return s.tenantRepo.Find(ctx, qb)
}

// CountTenants 统计租户数量
func (s *TenantService) CountTenants(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	return s.tenantRepo.Count(ctx, qb)
}
