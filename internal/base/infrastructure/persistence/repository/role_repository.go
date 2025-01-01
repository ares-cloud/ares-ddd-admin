package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type ISysRoleRepo interface {
	baserepo.IBaseRepo[entity.Role, int64]
	GetByCode(ctx context.Context, code string) (*entity.Role, error)
	GetByRoleId(ctx context.Context, roleId int64) ([]*entity.RolePermissions, error)
	DeletePermissionsByRoleId(ctx context.Context, roleId int64) error
	GetByUserId(ctx context.Context, userId string) ([]*entity.Role, error)
	GetUserRoles(ctx context.Context, userId string) ([]*entity.SysUserRole, error)
	FindByIds(ctx context.Context, ids []int64) ([]*entity.Role, error)
	FindAllEnabled(ctx context.Context) ([]*entity.Role, error)
	Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*entity.Role, error)
	UpdatePermissions(ctx context.Context, roleID int64, permIDs []int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permissions, error)
	HasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error)
	FindByPermissionID(ctx context.Context, permissionID int64) ([]*entity.Role, error)
	FindByType(ctx context.Context, roleType int8) ([]*entity.Role, error)
	GetIdsByUserId(ctx context.Context, userId string) ([]int64, error)
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	GetUserCountByRoleID(ctx context.Context, roleID int64) (int64, error)
	GetUsersByRoleID(ctx context.Context, roleID int64) ([]*entity.SysUser, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Role, error)
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
		permIds := make([]int64, 0)
		for _, perm := range rolePerms {
			permIds = append(permIds, perm.PermissionID)
		}
		perms, err := r.permissionsRepo.FindByIds(ctx, permIds)
		if err != nil {
			return nil, err
		}
		permissions = r.permMapper.ToDomainList(perms, nil)
	}

	// 换为领域模型
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
func (r *roleRepository) ExistsById(ctx context.Context, id int64) (bool, error) {
	_, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

// GetRolePermissions 获取角色权限
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*model.Permissions, error) {
	// 获取权限实体列表
	perms, err := r.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	return r.permMapper.ToDomainList(perms, nil), nil
}

// AssignPermissions 分配权限
func (r *roleRepository) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return r.repo.UpdatePermissions(ctx, roleID, permissionIDs)
}

// ExistsPermission 检查权限是否存在
func (r *roleRepository) ExistsPermission(ctx context.Context, permissionID int64) (bool, error) {
	return r.permissionsRepo.ExistsById(ctx, permissionID)
}

// IsRoleInUse 检查角色是否被使用
func (r *roleRepository) IsRoleInUse(ctx context.Context, roleID int64) (bool, error) {
	// 获取角色关联的用户数量
	userCount, err := r.repo.GetUserCountByRoleID(ctx, roleID)
	if err != nil {
		return false, err
	}
	return userCount > 0, nil
}
