package repository

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ISysRoleRepo interface {
	baserepo.IBaseRepo[entity.Role, int64]
	GetByCode(ctx context.Context, code string) (*entity.Role, error)
	GetByRoleId(ctx context.Context, roleId int64) ([]*entity.RolePermissions, error)
	DeletePermissionsByRoleId(ctx context.Context, roleId int64) error
	GetByUserId(ctx context.Context, userId string) ([]*entity.SysUserRole, error)
	FindByIds(ctx context.Context, ids []int64) ([]*entity.Role, error)
}

type roleRepository struct {
	repo            ISysRoleRepo
	permissionsRepo IPermissionsRepo
	mapper          *mapper.RoleMapper
	permMapper      *mapper.PermissionsMapper
}

func NewRoleRepository(repo ISysRoleRepo, permissionsRepo IPermissionsRepo) drepository.IRoleRepository {
	return &roleRepository{
		repo:            repo,
		permissionsRepo: permissionsRepo,
		mapper:          &mapper.RoleMapper{},
		permMapper:      &mapper.PermissionsMapper{},
	}
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	roleEntity := r.mapper.ToEntity(role)
	err := r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		dm, err := r.repo.Add(ctx, roleEntity)
		if err != nil {
			return err
		}
		if len(role.Permissions) > 0 {
			// 创建角色权限关联
			for _, perm := range role.Permissions {
				rolePermission := &entity.RolePermissions{
					RoleID:       dm.ID,
					PermissionID: perm.ID,
				}
				if err = r.repo.Db(ctx).Create(rolePermission).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	roleEntity := r.mapper.ToEntity(role)
	err := r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		err := r.repo.EditById(ctx, roleEntity.ID, roleEntity)
		if err != nil {
			return err
		}
		err = r.repo.DeletePermissionsByRoleId(ctx, roleEntity.ID)
		if err != nil {
			return err
		}
		if len(role.Permissions) > 0 {
			// 创建角色权限关联
			for _, perm := range role.Permissions {
				rolePermission := &entity.RolePermissions{
					RoleID:       roleEntity.ID,
					PermissionID: perm.ID,
				}
				if err = r.repo.Db(ctx).Create(rolePermission).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (r *roleRepository) FindByID(ctx context.Context, id int64) (*model.Role, error) {
	// 查询角色基本信息
	roleEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	// 查询角色权限关联
	rolePerms, err := r.repo.GetByRoleId(ctx, id)
	if err != nil {
		return nil, err
	}
	permissions := make([]*model.Permissions, 0)
	if len(rolePerms) > 0 {
		//permIds := make([]int64, 0)
		//for _, perm := range rolePerms {
		//	permIds = append(permIds, perm.PermissionID)
		//}
		//perms, err := r.permissionsRepo.FindByIds(ctx, permIds)
		//if err != nil {
		//	return nil, err
		//}
		//permissions = r.permMapper.ToDomainList(perms, nil)
		for _, perm := range rolePerms {
			permissions = append(permissions, &model.Permissions{
				ID: perm.PermissionID,
			})
		}
	}

	// 转换为领域模型
	return r.mapper.ToDomain(roleEntity, permissions), nil
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	roleEntity, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	return r.mapper.ToDomain(roleEntity, nil), nil
}

func (r *roleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	_, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *roleRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Role, error) {
	roleEntities, err := r.repo.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomainList(roleEntities), nil
}

func (r *roleRepository) Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.Role, error) {
	records, err := r.repo.Find(ctx, qb)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	roles := r.mapper.ToDomainList(records)
	return roles, nil
}

func (r *roleRepository) Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, qb)
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

func (r *roleRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Role, error) {
	// 查询用户角色关联
	userRoles, err := r.repo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []*model.Role{}, nil
	}

	// 获取角色ID列表
	roleIds := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleID)
	}

	// 查询角色信息
	roles, err := r.repo.FindByIds(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(roles), nil
}
