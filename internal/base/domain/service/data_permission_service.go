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

type DataPermissionService struct {
	permRepo repository.IDataPermissionRepository
	roleRepo repository.IRoleRepository
	eventBus pkgEvent.IEventBus
}

func NewDataPermissionService(
	permRepo repository.IDataPermissionRepository,
	roleRepo repository.IRoleRepository,
	eventBus pkgEvent.IEventBus,
) *DataPermissionService {
	return &DataPermissionService{
		permRepo: permRepo,
		roleRepo: roleRepo,
		eventBus: eventBus,
	}
}

// AssignDataPermission 分配数据权限
func (s *DataPermissionService) AssignDataPermission(ctx context.Context, perm *model.DataPermission) herrors.Herr {
	// 1. 验证数据权限
	if err := perm.Validate(); err != nil {
		return errors.DataPermissionInvalidField("", err.Error())
	}

	// 2. 验证角色是否存在
	exists, err := s.roleRepo.ExistsById(ctx, perm.RoleID)
	if err != nil {
		return errors.DataPermissionQueryFailed(err)
	}
	if !exists {
		return errors.RoleNotFound(perm.RoleID)
	}

	// 3. 保存数据权限
	if err := s.permRepo.Save(ctx, perm); err != nil {
		return errors.DataPermissionCreateFailed(err)
	}

	// 4. 发布事件
	err1 := s.eventBus.Publish(ctx, events.NewDataPermissionEvent(actx.GetTenantId(ctx), perm, events.DataPermissionAssigned))
	if err1 != nil {
		return herrors.NewServerHError(err1)
	}
	return nil
}

// RemoveDataPermission 移除数据权限
func (s *DataPermissionService) RemoveDataPermission(ctx context.Context, roleID int64) herrors.Herr {
	// 1. 验证角色是否存在
	dataPermission, herr := s.GetByRoleID(ctx, roleID)
	if herrors.HaveError(herr) {
		return herr
	}

	// 2. 删除数据权限
	if err := s.permRepo.DeleteByRoleID(ctx, roleID); err != nil {
		return errors.DataPermissionDeleteFailed(err)
	}

	// 3. 发布事件
	err1 := s.eventBus.Publish(ctx, events.NewDataPermissionEvent(actx.GetTenantId(ctx), dataPermission, events.DataPermissionRemoved))
	if err1 != nil {
		return herrors.NewServerHError(err1)
	}
	return nil
}

// GetByRoleID 获取角色的数据权限
func (s *DataPermissionService) GetByRoleID(ctx context.Context, roleID int64) (*model.DataPermission, herrors.Herr) {
	// 1. 验证角色是否存在
	exists, err := s.roleRepo.ExistsById(ctx, roleID)
	if err != nil {
		return nil, errors.DataPermissionQueryFailed(err)
	}
	if !exists {
		return nil, errors.RoleNotFound(roleID)
	}

	// 2. 获取数据权限
	perm, err := s.permRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, errors.DataPermissionQueryFailed(err)
	}
	if perm == nil {
		return nil, errors.DataPermissionNotFound(roleID)
	}

	return perm, nil
}

// GetByRoleIDs 批量获取角色的数据权限
func (s *DataPermissionService) GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.DataPermission, herrors.Herr) {
	// 1. 验证角色是否都存在
	for _, roleID := range roleIDs {
		exists, err := s.roleRepo.ExistsById(ctx, roleID)
		if err != nil {
			return nil, errors.DataPermissionQueryFailed(err)
		}
		if !exists {
			return nil, errors.RoleNotFound(roleID)
		}
	}

	// 2. 批量获取数据权限
	perms, err := s.permRepo.GetByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, errors.DataPermissionQueryFailed(err)
	}

	return perms, nil
}
