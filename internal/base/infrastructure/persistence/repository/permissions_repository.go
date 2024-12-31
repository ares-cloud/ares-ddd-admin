package repository

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IPermissionsRepo ， 系统菜单
type IPermissionsRepo interface {
	baserepo.IBaseRepo[entity.Permissions, int64]
	DelByPermissionsId(ctx context.Context, permissionsId int64) error
	SavePermissionsResource(ctx context.Context, permissionsResource *entity.PermissionsResource) error
	GetByPermissionsId(ctx context.Context, permissionsId int64) ([]*entity.PermissionsResource, error)
	GetResourceByPermissionsIds(ctx context.Context, permissionsId []int64) ([]*entity.PermissionsResource, error)
	GetByCode(ctx context.Context, code string) (*entity.Permissions, error)
	GetByRoleID(ctx context.Context, roleID int64) ([]*entity.Permissions, []*entity.PermissionsResource, error)
	GetAllTree(ctx context.Context) ([]*entity.Permissions, []int64, error)
	GetTreeByType(ctx context.Context, permType int8) ([]*entity.Permissions, []*entity.PermissionsResource, error)
	GetTreeByQuery(ctx context.Context, qb *db_query.QueryBuilder) ([]*entity.Permissions, []*entity.PermissionsResource, error)
	GetTreeByUserAndType(ctx context.Context, userID string, permType int8) ([]*entity.Permissions, []*entity.PermissionsResource, error)
	GetResourcesByRoles(ctx context.Context, roles []int64) ([]*entity.PermissionsResource, error)
	GetByRoles(ctx context.Context, roles []int64) ([]*entity.Permissions, error)
	GetResourcesByRolesGrouped(ctx context.Context, roles []int64) (map[int64][]*entity.PermissionsResource, error)
	ExistsById(ctx context.Context, permissionID int64) (bool, error)
}

type permissionsRepository struct {
	repo   IPermissionsRepo
	mapper *mapper.PermissionsMapper
}

func NewPermissionsRepository(repo IPermissionsRepo) drepository.IPermissionsRepository {
	return &permissionsRepository{
		repo:   repo,
		mapper: &mapper.PermissionsMapper{},
	}
}

func (r *permissionsRepository) Create(ctx context.Context, permissions *model.Permissions) error {
	permEntity, resources := r.mapper.ToEntity(permissions)
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 创建权限基本信息
		dm, err := r.repo.Add(ctx, permEntity)
		if err != nil {
			return err
		}

		// 创建权限资源
		if len(resources) > 0 {
			for _, resource := range resources {
				resource.PermissionsID = dm.ID
				if err := r.repo.SavePermissionsResource(ctx, resource); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *permissionsRepository) Update(ctx context.Context, permissions *model.Permissions) error {
	permEntity, resources := r.mapper.ToEntity(permissions)
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 更新权限基本信息
		err := r.repo.EditById(ctx, permEntity.ID, permEntity)
		if err != nil {
			return err
		}

		// 删除原有资源
		if err = r.repo.DelByPermissionsId(ctx, permEntity.ID); err != nil {
			return err
		}

		// 创建新的资源
		if len(resources) > 0 {
			for _, resource := range resources {
				resource.PermissionsID = permEntity.ID
				if err = r.repo.SavePermissionsResource(ctx, resource); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *permissionsRepository) Delete(ctx context.Context, id int64) error {
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 删除权限资源
		if err := r.repo.DelByPermissionsId(ctx, id); err != nil {
			return err
		}

		// 删除角色权限关联
		if err := r.repo.Db(ctx).Where("permission_id = ?", id).Delete(&entity.RolePermissions{}).Error; err != nil {
			return err
		}

		// 删除权限
		return r.repo.DelByIdUnScoped(ctx, id)
	})
}

func (r *permissionsRepository) FindByCode(ctx context.Context, code string) (*model.Permissions, error) {
	// 查询权限基本信息
	permEntity, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	// 查询权限资源
	resource, err := r.repo.GetByPermissionsId(ctx, permEntity.ID)
	if err != nil && !database.IfErrorNotFound(err) {
		return nil, err
	}

	return r.mapper.ToDomain(permEntity, resource), nil
}

func (r *permissionsRepository) FindByRoleID(ctx context.Context, roleID int64) ([]*model.Permissions, error) {
	// 获取权限和资源数据
	permEntities, resources, err := r.repo.GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if len(permEntities) == 0 {
		return []*model.Permissions{}, nil
	}

	// 转换为领域模型
	return r.mapper.ToDomainList(permEntities, resources), nil
}

func (r *permissionsRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	// 查询权限基本信息
	permEntity, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return permEntity != nil, nil
}

func (r *permissionsRepository) FindByType(ctx context.Context, permType int8) ([]*model.Permissions, error) {
	// 查询权限基本信息
	permEntities, err := r.repo.Find(ctx, db_query.NewQueryBuilder().Where("type", db_query.Eq, permType))
	if err != nil {
		return nil, err
	}

	if len(permEntities) == 0 {
		return []*model.Permissions{}, nil
	}

	// 获取权限ID列表
	permIDs := make([]int64, 0, len(permEntities))
	for _, p := range permEntities {
		permIDs = append(permIDs, p.ID)
	}

	// 查询权限资源
	var resources []*entity.PermissionsResource
	resources, err = r.repo.GetResourceByPermissionsIds(ctx, permIDs)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomainList(permEntities, resources), nil
}

func (r *permissionsRepository) FindByID(ctx context.Context, id int64) (*model.Permissions, error) {
	// 查询权限基本信息
	permEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// 查询权限资源
	resource, err := r.repo.GetByPermissionsId(ctx, id)
	if err != nil && !database.IfErrorNotFound(err) {
		return nil, err
	}

	return r.mapper.ToDomain(permEntity, resource), nil
}

func (r *permissionsRepository) Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*model.Permissions, error) {
	permissions, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("find permissions failed: %w", err)
	}
	// 获取所有权限的资源
	var resources []*entity.PermissionsResource
	if len(permissions) > 0 {
		permIDs := make([]int64, 0, len(permissions))
		for _, p := range permissions {
			permIDs = append(permIDs, p.ID)
		}
		resources, err = r.repo.GetResourceByPermissionsIds(ctx, permIDs)
		if err != nil {
			return nil, err
		}
	}
	return r.mapper.ToDomainList(permissions, resources), nil
}

func (r *permissionsRepository) Count(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, qb)
}

func (r *permissionsRepository) FindAllTree(ctx context.Context) ([]*model.Permissions, []int64, error) {
	// 获取权限数据
	permEntities, ids, err := r.repo.GetAllTree(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(permEntities) == 0 {
		return []*model.Permissions{}, []int64{}, nil
	}

	// 构建ID到权限的映射
	permMap := make(map[int64]*model.Permissions)
	var roots []*model.Permissions

	// 先转换所有权限
	for _, entity := range permEntities {
		perm := &model.Permissions{
			ID:       entity.ID,
			Code:     entity.Code,
			Name:     entity.Name,
			Localize: entity.Localize,
			Icon:     entity.Icon,
			ParentID: entity.ParentID,
			Children: make([]*model.Permissions, 0),
		}
		permMap[perm.ID] = perm

		// 收集根节点
		if perm.ParentID == 0 {
			roots = append(roots, perm)
		}
	}

	// 构建树结构
	for _, perm := range permMap {
		if perm.ParentID != 0 {
			if parent, ok := permMap[perm.ParentID]; ok {
				parent.Children = append(parent.Children, perm)
			}
		}
	}

	return roots, ids, nil
}

func (r *permissionsRepository) FindTreeByType(ctx context.Context, permType int8) ([]*model.Permissions, error) {
	var permissions []*entity.Permissions
	var resources []*entity.PermissionsResource
	var err error
	permissions, resources, err = r.repo.GetTreeByType(ctx, permType)
	if err != nil {
		return nil, err
	}

	// return buildPermissionTree(r.mapper.ToDomainList(permissions, resources)), nil
	return r.mapper.ToDomainList(permissions, resources), nil
}

func (r *permissionsRepository) FindTreeByQuery(ctx context.Context, qb *db_query.QueryBuilder) ([]*model.Permissions, error) {
	permEntities, resources, err := r.repo.GetTreeByQuery(ctx, qb)
	if err != nil {
		return nil, err
	}

	return buildPermissionTree(r.mapper.ToDomainList(permEntities, resources)), nil
}

func (r *permissionsRepository) FindTreeByUserAndType(ctx context.Context, userID string, permType int8) ([]*model.Permissions, error) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	if permType > 0 {
		qb.Where("type", db_query.Eq, permType)
	}
	qb.Where("status", db_query.Eq, 1)
	qb.OrderBy("sequence", false)

	// 查询用户权限
	permissions, resources, err := r.repo.GetTreeByUserAndType(ctx, userID, permType)
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	return buildPermissionTree(r.mapper.ToDomainList(permissions, resources)), nil
}

// 辅助函数：构建权限树
func buildPermissionTree(permissions []*model.Permissions) []*model.Permissions {
	// 创建ID到权限的映射
	permMap := make(map[int64]*model.Permissions)
	for _, p := range permissions {
		permMap[p.ID] = p
	}

	// 构建树结构
	var roots []*model.Permissions
	for _, p := range permissions {
		if p.ParentID == 0 {
			roots = append(roots, p)
			continue
		}
		if parent, ok := permMap[p.ParentID]; ok {
			if parent.Children == nil {
				parent.Children = make([]*model.Permissions, 0)
			}
			parent.Children = append(parent.Children, p)
		}
	}

	return roots
}

func (r *permissionsRepository) FindAllEnabled(ctx context.Context) ([]*model.Permissions, error) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("status", db_query.Eq, 1)
	qb.OrderBy("sequence", true)

	// 查询权限
	permissions, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	// 查询权限资源
	var resources []*entity.PermissionsResource
	if len(permissions) > 0 {
		permissionIds := make([]int64, 0, len(permissions))
		for _, p := range permissions {
			permissionIds = append(permissionIds, p.ID)
		}
		resources, err = r.repo.GetResourceByPermissionsIds(ctx, permissionIds)
		if err != nil {
			return nil, err
		}
	}

	return r.mapper.ToDomainList(permissions, resources), nil
}

func (r *permissionsRepository) GetResourcesByRoles(ctx context.Context, roles []string) ([]*entity.PermissionsResource, error) {
	var resources []*entity.PermissionsResource

	// 1. 先查询角色对应的权限ID
	var permissionIDs []int64
	err := r.repo.Db(ctx).Model(&entity.RolePermissions{}).
		Joins("JOIN sys_role ON sys_role.id = sys_role_permissions.role_id").
		Where("sys_role.code IN ? AND sys_role.status = ?", roles, 1).
		Pluck("permission_id", &permissionIDs).Error
	if err != nil {
		return nil, err
	}

	if len(permissionIDs) == 0 {
		return []*entity.PermissionsResource{}, nil
	}

	// 2. 查询启用的权限资源
	err = r.repo.Db(ctx).Model(&entity.PermissionsResource{}).
		Joins("JOIN sys_permissions ON sys_permissions.id = sys_permissions_resource.permissions_id").
		Where("sys_permissions.id IN ? AND sys_permissions.status = ? AND sys_permissions.type = ?",
			permissionIDs, 1, 3). // type=3 表示API类型
		Find(&resources).Error
	if err != nil {
		return nil, err
	}

	return resources, nil
}
