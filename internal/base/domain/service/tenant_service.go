package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	pkgEvents "github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type TenantCommandService struct {
	tenantRepo repository.ITenantRepository
	publisher  pkgEvents.IEventBus
}

func NewTenantCommandService(
	tenantRepo repository.ITenantRepository,
	publisher pkgEvents.IEventBus,
) *TenantCommandService {
	return &TenantCommandService{
		tenantRepo: tenantRepo,
		publisher:  publisher,
	}
}

// CreateTenant 创建租户
func (s *TenantCommandService) CreateTenant(ctx context.Context, tenant *model.Tenant) herrors.Herr {
	// 验证租户模型
	if err := tenant.Validate(); err != nil {
		return err
	}

	// 检查租户编码是否已存在
	exists, err := s.tenantRepo.ExistsByCode(ctx, tenant.Code)
	if err != nil {
		return herrors.NewErr(err)
	}
	if exists {
		return errors.TenantCodeExists(tenant.Code)
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return herrors.NewErr(err)
	}

	// 发布租户创建事件
	event := events.NewTenantEvent(tenant.ID, events.TenantCreated)
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// UpdateTenant 更新租户
func (s *TenantCommandService) UpdateTenant(ctx context.Context, tenant *model.Tenant) herrors.Herr {
	// 检查租户是否存在
	old, err := s.tenantRepo.FindByID(ctx, tenant.ID)
	if err != nil {
		return herrors.NewErr(err)
	}
	if old == nil {
		return errors.TenantNotFound(tenant.ID)
	}

	// 检查租户状态
	if locked, reason := old.IsLocked(); locked {
		return errors.TenantDisabled(reason)
	}

	// 验证租户模型
	if err := tenant.Validate(); err != nil {
		return err
	}

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return herrors.NewErr(err)
	}

	// 发布租户更新事件
	event := events.NewTenantEvent(tenant.ID, events.TenantUpdated)
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// DeleteTenant 删除租户
func (s *TenantCommandService) DeleteTenant(ctx context.Context, id string) herrors.Herr {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return herrors.NewErr(err)
	}
	if tenant == nil {
		return errors.TenantNotFound(id)
	}

	// 检查租户状态
	if locked, reason := tenant.IsLocked(); locked {
		return errors.TenantDisabled(reason)
	}

	// 检查是否为默认租户
	if tenant.IsDefaultTenant() {
		return errors.TenantIsDefault()
	}

	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return herrors.NewErr(err)
	}

	// 发布租户删除事件
	event := events.NewTenantEvent(tenant.ID, events.TenantDeleted)
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// AssignPermissions 分配权限
func (s *TenantCommandService) AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) herrors.Herr {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return herrors.NewErr(err)
	}
	if tenant == nil {
		return errors.TenantNotFound(tenantID)
	}

	// 检查租户状态
	if locked, reason := tenant.IsLocked(); locked {
		return errors.TenantDisabled(reason)
	}

	if err := s.tenantRepo.AssignPermissions(ctx, tenantID, permissionIDs); err != nil {
		return herrors.NewErr(err)
	}

	// 发布权限变更事件
	err = s.publisher.Publish(ctx, events.NewTenantPermissionEvent(tenantID, permissionIDs))
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// GetTenant 获取租户信息
func (s *TenantCommandService) GetTenant(ctx context.Context, id string) (*model.Tenant, herrors.Herr) {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, herrors.NewErr(err)
	}
	if tenant == nil {
		return nil, errors.TenantNotFound(id)
	}
	return tenant, nil
}

// ExistsByCode 检查租户编码是否存在
func (s *TenantCommandService) ExistsByCode(ctx context.Context, code string) (bool, herrors.Herr) {
	exists, err := s.tenantRepo.ExistsByCode(ctx, code)
	if err != nil {
		return false, herrors.NewErr(err)
	}
	return exists, nil
}

// LockTenant 锁定租户
func (s *TenantCommandService) LockTenant(ctx context.Context, id string, reason string) herrors.Herr {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return herrors.NewErr(err)
	}
	if tenant == nil {
		return errors.TenantNotFound(id)
	}

	// 锁定租户
	if err := tenant.Lock(reason); err != nil {
		return err
	}

	// 保存更新
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return herrors.NewErr(err)
	}

	// 发布租户锁定事件
	event := events.NewTenantEvent(tenant.ID, events.TenantLocked)
	if err := s.publisher.Publish(ctx, event); err != nil {
		return herrors.NewErr(err)
	}

	return nil
}

// UnlockTenant 解锁租户
func (s *TenantCommandService) UnlockTenant(ctx context.Context, id string) herrors.Herr {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return herrors.NewErr(err)
	}
	if tenant == nil {
		return errors.TenantNotFound(id)
	}

	// 解锁租户
	if err := tenant.Unlock(); err != nil {
		return err
	}

	// 保存更新
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return herrors.NewErr(err)
	}

	// 发布租户解锁事件
	event := events.NewTenantEvent(tenant.ID, events.TenantUnlocked)
	if err := s.publisher.Publish(ctx, event); err != nil {
		return herrors.NewErr(err)
	}

	return nil
}
