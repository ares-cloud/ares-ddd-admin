package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type RoleService struct {
	roleRepo repository.IRoleRepository
	userRepo repository.IUserRepository
	eventBus *pkgEvent.EventBus
}

func NewRoleService(roleRepo repository.IRoleRepository, userRepo repository.IUserRepository, eventBus *pkgEvent.EventBus) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// CheckRolePermissions 检查角色是否具有特定权限
func (s *RoleService) CheckRolePermissions(ctx context.Context, roleID int64, permissionCode string) (bool, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return false, err
	}

	for _, perm := range role.Permissions {
		if perm.Code == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

// GetRolePermissionCodes 获取角色的所有权限代码
func (s *RoleService) GetRolePermissionCodes(ctx context.Context, roleID int64) ([]string, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(role.Permissions))
	for i, perm := range role.Permissions {
		codes[i] = perm.Code
	}

	return codes, nil
}

// GetAllDataPermissionRoles 获取所有数据权限角色
func (s *RoleService) GetAllDataPermissionRoles(ctx context.Context) ([]*model.Role, error) {
	// 构建查询条件
	qb := query.NewQueryBuilder().
		Where("type", query.Eq, int8(entity.RoleTypeData)).
		Where("status", query.Eq, 1) // 只查询启用状态的角色

	// 查询数据权限角色
	roles, err := s.roleRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return err
	}

	return s.publishRoleEvent(ctx, role.ID, role.TenantID)
}

func (s *RoleService) UpdateRolePermissions(ctx context.Context, roleID int64, permIDs []int64) error {
	if err := s.roleRepo.UpdatePermissions(ctx, roleID, permIDs); err != nil {
		return err
	}

	// 获取角色信息
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	return s.publishRoleEvent(ctx, roleID, role.TenantID)
}

// publishRoleEvent 发布角色相关事件
func (s *RoleService) publishRoleEvent(ctx context.Context, roleID int64, tenantID string) error {
	//// 获取角色关联的用户
	//users, err := s.userRepo.FindByRoleID(ctx, roleID)
	//if err != nil {
	//	return err
	//}
	//
	//// 构建用户ID列表
	//userIDs := make([]string, len(users))
	//for i, user := range users {
	//	userIDs[i] = user.ID
	//}
	//
	//// 发布事件
	//event := &events.RoleEvent{
	//	BaseEvent: events.BaseEvent{TenantID: tenantID},
	//	RoleID:    roleID,
	//	UserIDs:   userIDs,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}
