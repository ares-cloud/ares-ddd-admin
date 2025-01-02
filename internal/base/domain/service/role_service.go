package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	domanevent "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type RoleCommandService struct {
	roleRepo repository.IRoleRepository
	eventBus events.IEventBus
}

func NewRoleCommandService(
	roleRepo repository.IRoleRepository,
	eventBus events.IEventBus,
) *RoleCommandService {
	return &RoleCommandService{
		roleRepo: roleRepo,
		eventBus: eventBus,
	}
}

// CreateRole 创建角色
func (s *RoleCommandService) CreateRole(ctx context.Context, role *model.Role) herrors.Herr {
	// 1. 验证角色
	if hr := role.Validate(); herrors.HaveError(hr) {
		return hr
	}

	// 2. 检查编码是否存在
	exists, err := s.roleRepo.ExistsByCode(ctx, role.Code)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if exists {
		return errors.RoleExists(role.Code)
	}

	// 3. 创建角色
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return herrors.NewServerHError(err)
	}

	// 4. 发布角色创建事件
	err = s.eventBus.Publish(ctx, domanevent.NewRoleEvent(role.TenantID, role.ID, domanevent.RoleCreated))
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// UpdateRole 更新角色
func (s *RoleCommandService) UpdateRole(ctx context.Context, role *model.Role) herrors.Herr {
	// 1. 验证角色
	if hr := role.Validate(); herrors.HaveError(hr) {
		return hr
	}

	// 2. 检查编码是否存在
	exists, err := s.roleRepo.FindByCode(ctx, role.Code)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if exists != nil && exists.ID != role.ID {
		return errors.RoleExists(role.Code)
	}

	// 3. 更新角色
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return herrors.NewServerHError(err)
	}

	// 4. 发布角色更新事件
	err = s.eventBus.Publish(ctx, domanevent.NewRoleEvent(role.TenantID, role.ID, domanevent.RoleUpdated))
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// DeleteRole 删除角色
func (s *RoleCommandService) DeleteRole(ctx context.Context, id int64) herrors.Herr {
	// 1. 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if role == nil {
		return errors.RoleNotFound(id)
	}

	// 2. 检查角色是否被使用
	used, err := s.roleRepo.IsRoleInUse(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if used {
		return errors.RoleInUse(id)
	}

	// 3. 删除角色
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return herrors.NewServerHError(err)
	}

	// 4. 发布角色删除事件
	err = s.eventBus.Publish(ctx, domanevent.NewRoleEvent(role.TenantID, role.ID, domanevent.RoleDeleted))
	if err != nil {
		return herrors.NewErr(err)
	}
	return nil
}

// AssignPermissions 分配权限
func (s *RoleCommandService) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) herrors.Herr {
	// 1. 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if role == nil {
		return errors.RoleNotFound(roleID)
	}

	// 2. 检查权限是否存在
	if err := s.validatePermissions(ctx, permissionIDs); err != nil {
		return herrors.NewServerHError(err)
	}

	// 3. 分配权限
	if err := s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs); err != nil {
		return herrors.NewServerHError(err)
	}

	// 4. 发布权限分配事件
	err = s.eventBus.Publish(ctx, domanevent.NewRolePermissionsAssignedEvent(roleID, permissionIDs))
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// GetRole 获取角色
func (s *RoleCommandService) GetRole(ctx context.Context, id int64) (*model.Role, herrors.Herr) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	if role == nil {
		return nil, errors.RoleNotFound(id)
	}
	return role, nil
}

// validatePermissions 验证权限是否存在
func (s *RoleCommandService) validatePermissions(ctx context.Context, permissionIDs []int64) error {
	for _, id := range permissionIDs {
		exists, err := s.roleRepo.ExistsPermission(ctx, id)
		if err != nil {
			return err
		}
		if !exists {
			return errors.PermissionNotFound(id)
		}
	}
	return nil
}
