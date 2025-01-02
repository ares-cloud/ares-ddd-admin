package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type PermissionService struct {
	permRepo repository.IPermissionsRepository
	eventBus pkgEvent.IEventBus
}

func NewPermissionService(
	permRepo repository.IPermissionsRepository,
	eventBus pkgEvent.IEventBus,
) *PermissionService {
	return &PermissionService{
		permRepo: permRepo,
		eventBus: eventBus,
	}
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(ctx context.Context, perm *model.Permissions) herrors.Herr {
	// 1. 检查权限编码是否已存在
	exists, err := s.permRepo.ExistsByCode(ctx, perm.Code)
	if err != nil {
		return errors.PermissionQueryFailed(err)
	}
	if exists {
		return errors.PermissionExists(perm.Code)
	}

	// 2. 创建权限
	if err := s.permRepo.Create(ctx, perm); err != nil {
		return errors.PermissionCreateFailed(err)
	}

	// 3. 发布权限创建事件
	err = s.eventBus.Publish(ctx, events.NewPermissionEvent(actx.GetTenantId(ctx), perm.ID, events.PermissionCreated))
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(ctx context.Context, perm *model.Permissions) herrors.Herr {
	// 1. 检查权限是否存在
	oldPerm, err := s.permRepo.FindByID(ctx, perm.ID)
	if err != nil {
		return errors.PermissionQueryFailed(err)
	}
	if oldPerm == nil {
		return errors.PermissionNotFound(perm.ID)
	}

	// 2. 更新权限
	if err := s.permRepo.Update(ctx, perm); err != nil {
		return errors.PermissionUpdateFailed(err)
	}

	// 3. 发布权限更新事件
	err = s.eventBus.Publish(ctx, events.NewPermissionEvent(actx.GetTenantId(ctx), perm.ID, events.PermissionUpdated))
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx context.Context, id int64) herrors.Herr {
	// 1. 检查权限是否存在
	perm, err := s.permRepo.FindByID(ctx, id)
	if err != nil {
		return errors.PermissionQueryFailed(err)
	}
	if perm == nil {
		return errors.PermissionNotFound(id)
	}

	// 2. 检查是否有子权限
	if len(perm.Children) > 0 {
		return errors.HasChildPermission(id)
	}

	// 3. 删除权限
	if err := s.permRepo.Delete(ctx, id); err != nil {
		return errors.PermissionDeleteFailed(err)
	}

	// 4. 发布权限删除事件
	err = s.eventBus.Publish(ctx, events.NewPermissionEvent(actx.GetTenantId(ctx), perm.ID, events.PermissionDeleted))
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// FindByID 根据ID查询权限
func (s *PermissionService) FindByID(ctx context.Context, id int64) (*model.Permissions, herrors.Herr) {
	perm, err := s.permRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.PermissionQueryFailed(err)
	}
	if perm == nil {
		return nil, errors.PermissionNotFound(id)
	}
	return perm, nil
}

// UpdatePermissionStatus 更新权限状态
func (s *PermissionService) UpdatePermissionStatus(ctx context.Context, id int64, status int8) herrors.Herr {
	// 1. 检查权限是否存在
	perm, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 更新状态
	if err := perm.UpdateStatus(status); err != nil {
		return err
	}

	// 3. 保存更新
	if err := s.permRepo.Update(ctx, perm); err != nil {
		return errors.PermissionUpdateFailed(err)
	}

	// 4. 发布权限更新事件
	err1 := s.eventBus.Publish(ctx, events.NewPermissionEvent(actx.GetTenantId(ctx), perm.ID, events.PermissionStatusChange))
	if err1 != nil {
		return herrors.NewServerHError(err1)
	}

	return nil
}

// ValidatePermission 验证权限信息
func (s *PermissionService) ValidatePermission(ctx context.Context, perm *model.Permissions) herrors.Herr {
	// 1. 基础验证
	if err := perm.Validate(); err != nil {
		return err
	}

	// 2. 如果有父权限,检查父权限是否存在且有效
	if perm.ParentID > 0 {
		parent, err := s.FindByID(ctx, perm.ParentID)
		if err != nil {
			return err
		}
		if !parent.IsEnabled() {
			return errors.PermissionInvalidField("parent_id", "parent permission is disabled")
		}
	}

	return nil
}

// CheckPermissionExists 检查权限是否存在
func (s *PermissionService) CheckPermissionExists(ctx context.Context, code string) (bool, herrors.Herr) {
	exists, err := s.permRepo.ExistsByCode(ctx, code)
	if err != nil {
		return false, errors.PermissionQueryFailed(err)
	}
	return exists, nil
}

// CheckHasChildren 检查是否有子权限
func (s *PermissionService) CheckHasChildren(ctx context.Context, id int64) (bool, herrors.Herr) {
	perm, err := s.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return perm.HasChildren(), nil
}
