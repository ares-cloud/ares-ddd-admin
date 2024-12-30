package repository

import (
	"context"
	"encoding/json"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
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
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*entity.Role, error)
	UpdatePermissions(ctx context.Context, roleID int64, permIDs []int64) error
	GetRoleDataPermission(ctx context.Context, roleID int64) (*entity.DataPermission, error)
	GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permissions, error)
	FindByPermissionID(ctx context.Context, permissionID int64) ([]*entity.Role, error)
	FindByType(ctx context.Context, roleType int8) ([]*entity.Role, error)
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
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
	return r.mapper.ToDomainList(userRoles), nil
}

func (r *roleRepository) FindAllEnabled(ctx context.Context) ([]*model.Role, error) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	qb.Where("status", query.Eq, 1)
	qb.OrderBy("sequence", false)

	// 查询角色
	roles, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(roles), nil
}

// FindByType 根据角色类型查询角色列表
func (r *roleRepository) FindByType(ctx context.Context, roleType int8) ([]*model.Role, error) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	qb.Where("type", query.Eq, roleType)
	qb.Where("status", query.Eq, 1) // 只查询启用状态的角色
	qb.OrderBy("sequence", false)   // 按sequence排序

	// 查询角色
	roles, err := r.repo.Find(ctx, qb)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	// 转换为领域模型
	return r.mapper.ToDomainList(roles), nil
}

// GetPermissionsByRoleID 获取角色的权限ID列表
func (r *roleRepository) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var permIDs []int64
	err := r.repo.Db(ctx).Model(&entity.RolePermissions{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permIDs).Error
	return permIDs, err
}

// FindByPermissionID 根据权限ID查找角色
func (r *roleRepository) FindByPermissionID(ctx context.Context, permissionID int64) ([]*model.Role, error) {
	roles, err := r.repo.FindByPermissionID(ctx, permissionID)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomainList(roles), nil
}

// UpdatePermissions 更新角色权限
func (r *roleRepository) UpdatePermissions(ctx context.Context, roleID int64, permIDs []int64) error {
	return r.repo.UpdatePermissions(ctx, roleID, permIDs)
}

// GetRoleDataPermission 获取角色数据权限
func (r *roleRepository) GetRoleDataPermission(ctx context.Context, roleID int64) (*model.DataPermission, error) {
	dataPerm, err := r.repo.GetRoleDataPermission(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if dataPerm == nil {
		return nil, nil
	}
	var deptIDs []string
	if err = json.Unmarshal([]byte(dataPerm.DeptIDs), &deptIDs); err != nil {
		return nil, err
	}
	return &model.DataPermission{
		ID:       dataPerm.ID,
		RoleID:   dataPerm.RoleID,
		Scope:    model.DataScope(dataPerm.Scope),
		DeptIDs:  deptIDs,
		TenantID: dataPerm.TenantID,
	}, nil
}

// GetRolePermissions 获取角色权限
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*model.Permissions, error) {
	permissions, err := r.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return r.permMapper.ToDomainList(permissions, nil), nil
}
