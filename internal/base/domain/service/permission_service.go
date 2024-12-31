package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type PermissionService struct {
	permissionRepo repository.IPermissionsRepository
	roleRepo       repository.IRoleRepository
	tenantRepo     repository.ITenantRepository
	eventBus       *pkgEvent.EventBus
}

func NewPermissionService(permissionRepo repository.IPermissionsRepository, roleRepo repository.IRoleRepository, tenantRepo repository.ITenantRepository, eventBus *pkgEvent.EventBus) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		tenantRepo:     tenantRepo,
		eventBus:       eventBus,
	}
}

// ValidatePermissionCode 验证权限代码是否有效
func (s *PermissionService) ValidatePermissionCode(ctx context.Context, code string) (bool, error) {
	return s.permissionRepo.ExistsByCode(ctx, code)
}

// GetPermissionsByType 根据类型获取权限列表
func (s *PermissionService) GetPermissionsByType(ctx context.Context, permType int8) ([]*model.Permissions, error) {
	return s.permissionRepo.FindByType(ctx, permType)
}

func (s *PermissionService) FindAllEnabled(ctx context.Context) ([]*model.Permissions, error) {
	tenantId := actx.GetTenantId(ctx)
	if tenantId != "" {
		//_, tenant, err := utils.IsTenantAdmin(ctx, nil, s.tenantRepo)
		//if err != nil {
		//	hlog.CtxErrorf(ctx, "isTenantAdmin err: %v", err)
		//	return nil, fmt.Errorf("%w: %v", errors.ErrInvalidTenant, err)
		//}
		//if tenant != nil && !tenant.IsDefaultTenant() {
		//	perms, err := s.tenantRepo.GetPermissions(ctx, tenantId)
		//	if err != nil {
		//		hlog.CtxErrorf(ctx, "GetPermissions err: %v", err)
		//		return nil, err
		//	}
		//	return perms, nil
		//}
	}
	return s.permissionRepo.FindAllEnabled(ctx)
}

func (s *PermissionService) UpdatePermission(ctx context.Context, perm *model.Permissions) error {
	if err := s.permissionRepo.Update(ctx, perm); err != nil {
		return err
	}

	//// 获取使用此权限的角色
	//roles, err := s.roleRepo.FindByPermissionID(ctx, perm.ID)
	//if err != nil {
	//	return err
	//}
	//
	//// 发布权限更新事件
	//roleIDs := make([]int64, len(roles))
	//for i, role := range roles {
	//	roleIDs[i] = role.ID
	//}

	//event := &events.PermissionEvent{
	//	BaseEvent:    events.BaseEvent{TenantID: actx.GetTenantId(ctx)},
	//	PermissionID: perm.ID,
	//	RoleIDs:      roleIDs,
	//}
	//return s.eventBus.Publish(ctx, event)
	return nil
}
