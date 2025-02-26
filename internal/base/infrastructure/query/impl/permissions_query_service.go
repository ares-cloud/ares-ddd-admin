package impl

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type PermissionsQueryService struct {
	permRepo             repository.IPermissionsRepo
	permissionsConverter *converter.PermissionsConverter
	tenantRepo           repository.ISysTenantRepo
}

func NewPermissionsQueryService(
	permRepo repository.IPermissionsRepo,
	tenantRepo repository.ISysTenantRepo,
	permissionsConverter *converter.PermissionsConverter,
) *PermissionsQueryService {
	return &PermissionsQueryService{
		permRepo:             permRepo,
		tenantRepo:           tenantRepo,
		permissionsConverter: permissionsConverter,
	}
}

// FindByID 根据ID查询权限
func (s *PermissionsQueryService) FindByID(ctx context.Context, id int64) (*dto.PermissionsDto, error) {
	perm, err := s.permRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, nil
	}
	resources, err := s.permRepo.GetByPermissionsId(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.permissionsConverter.ToDTO(perm, resources), nil
}

// Find 查询权限列表
func (s *PermissionsQueryService) Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.PermissionsDto, int64, herrors.Herr) {
	perms, err := s.permRepo.Find(ctx, qb)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}

	total, err := s.permRepo.Count(ctx, qb)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}

	return s.permissionsConverter.ToDTOList(perms), total, nil
}

// FindTreeByType 查询权限树
func (s *PermissionsQueryService) FindTreeByType(ctx context.Context, permType int8) ([]*dto.PermissionsDto, error) {
	perms, _, err := s.permRepo.GetTreeByType(ctx, permType)
	if err != nil {
		return nil, err
	}
	return s.permissionsConverter.ToTreeDTOList(perms), nil
}

// FindAllEnabled 查询所有启用的权限
func (s *PermissionsQueryService) FindAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, error) {
	tenantId := actx.GetTenantId(ctx)
	var perms []*entity.Permissions
	var err error
	// 没有租户就是全部,超级管理获取全部
	if tenantId == "" || actx.IsSuperAdmin(ctx) {
		// 构建查询条件
		qb := db_query.NewQueryBuilder()
		qb.Where("status", db_query.Eq, 1)
		qb.OrderBy("sequence", true)

		// 查询数据
		perms, err = s.permRepo.Find(ctx, qb)
	} else {
		perms, err = s.tenantRepo.GetPermissionsByTenantID(ctx, tenantId)
	}
	if err != nil {
		return nil, err
	}
	return s.permissionsConverter.ToDTOList(perms), nil
}

// GetSimplePermissionsTree 获取简化的权限树
func (s *PermissionsQueryService) GetSimplePermissionsTree(ctx context.Context) (*dto.PermissionsTreeResult, error) {
	// 1. 查询所有权限
	perms, _, err := s.permRepo.GetTreeByType(ctx, 1) // 1表示菜单类型
	if err != nil {
		return nil, err
	}
	// 获取ids
	ids := make([]int64, 0, len(perms))
	for _, perm := range perms {
		ids = append(ids, perm.ID)
	}
	tres := s.permissionsConverter.ToSimpleTreeDTOList(perms)
	// 2. 转换为树形结构
	return &dto.PermissionsTreeResult{
		Tree: tres,
		Ids:  ids,
	}, nil
}

// GetPermissionRoles 获取拥有该权限的角色列表
func (s *PermissionsQueryService) GetPermissionRoles(ctx context.Context, permID int64) ([]*dto.RoleDto, error) {
	//// 1. 构建查询条件
	//qb := db_query.NewQueryBuilder()
	//qb.InnerJoin("sys_role_permissions rp", "r.id = rp.role_id")
	//qb.Where("rp.permission_id", db_query.Eq, permID)
	//qb.Where("r.tenant_id", db_query.Eq, actx.GetTenantId(ctx))
	//
	//// 2. 查询数据
	//roles, err := s.permRepo.GetPermissionRoles(ctx, qb)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 3. 转换为DTO
	//return s.permissionsConverter.ToRoleDTOList(roles), nil
	return make([]*dto.RoleDto, 0), nil
}
