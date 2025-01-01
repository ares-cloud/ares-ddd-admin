package impl

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type RoleQueryService struct {
	roleRepo             repository.ISysRoleRepo
	converter            *converter.RoleConverter
	userConverter        *converter.UserConverter
	permissionsConverter *converter.PermissionsConverter
}

func NewRoleQueryService(
	roleRepo repository.ISysRoleRepo,
	converter *converter.RoleConverter,
	userConverter *converter.UserConverter,
	permissionsConverter *converter.PermissionsConverter,
) *RoleQueryService {
	return &RoleQueryService{
		roleRepo:             roleRepo,
		converter:            converter,
		permissionsConverter: permissionsConverter,
		userConverter:        userConverter,
	}
}

func (r *RoleQueryService) GetRole(ctx context.Context, id int64) (*dto.RoleDto, error) {
	role, err := r.roleRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// 获取角色权限ID列表
	permIds, err := r.roleRepo.GetPermissionsByRoleID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.converter.ToDTO(role, permIds), nil
}

func (r *RoleQueryService) FindRoles(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.RoleDto, error) {
	roles, err := r.roleRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}
	return r.converter.ToDTOList(roles), nil
}

func (r *RoleQueryService) CountRoles(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return r.roleRepo.Count(ctx, qb)
}

func (r *RoleQueryService) GetRolePermissions(ctx context.Context, roleID int64) ([]*dto.PermissionsDto, error) {
	perms, err := r.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return r.permissionsConverter.ToDTOList(perms), nil
}

func (r *RoleQueryService) FindByType(ctx context.Context, roleType int8) ([]*dto.RoleDto, error) {
	roles, err := r.roleRepo.FindByType(ctx, roleType)
	if err != nil {
		return nil, err
	}
	return r.converter.ToDTOList(roles), nil
}
func (r *RoleQueryService) GetRoleByCode(ctx context.Context, code string) (*dto.RoleDto, error) {
	role, err := r.roleRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.converter.ToDTO(role, nil), nil
}

// GetRoleUsers 获取角色下的用户列表
func (r *RoleQueryService) GetRoleUsers(ctx context.Context, roleID int64) ([]*dto.UserDto, error) {
	users, err := r.roleRepo.GetUsersByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return r.userConverter.ToDTOList(users), nil
}

// GetTenantRoles 获取租户下的角色列表
func (r *RoleQueryService) GetTenantRoles(ctx context.Context, tenantID string) ([]*dto.RoleDto, error) {
	roles, err := r.roleRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return r.converter.ToDTOList(roles), nil
}
