package repository

import (
	"context"
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
