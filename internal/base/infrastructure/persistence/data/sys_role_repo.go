package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
)

type sysRoleRepo struct {
	*baserepo.BaseRepo[entity.Role, int64]
}

func NewSysRoleRepo(data database.IDataBase) repository.ISysRoleRepo {
	model := new(entity.Role)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.RolePermissions{}); err != nil {
		hlog.Fatalf("sync sys user tables to db error: %v", err)
	}
	return &sysRoleRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.Role, int64](data, entity.Role{}),
	}
}

// GetByCode 根据编码获取角色
func (r *sysRoleRepo) GetByCode(ctx context.Context, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.Db(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByRoleId 获取角色的权限关联
func (r *sysRoleRepo) GetByRoleId(ctx context.Context, roleId int64) ([]*entity.RolePermissions, error) {
	var rolePerms []*entity.RolePermissions
	err := r.Db(ctx).Where("role_id = ?", roleId).Find(&rolePerms).Error
	if err != nil {
		return nil, err
	}
	return rolePerms, nil
}

// DeletePermissionsByRoleId 删除角色的权限关联
func (r *sysRoleRepo) DeletePermissionsByRoleId(ctx context.Context, roleId int64) error {
	return r.Db(ctx).Unscoped().Where("role_id = ?", roleId).Delete(&entity.RolePermissions{}).Error
}

// GetByUserId 根据用户ID获取角色列表
func (r *sysRoleRepo) GetByUserId(ctx context.Context, userId string) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.Db(ctx).Model(&entity.Role{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userId, 1).
		Order("sys_role.sequence").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserRoles 获取用户角色关联
func (r *sysRoleRepo) GetUserRoles(ctx context.Context, userId string) ([]*entity.SysUserRole, error) {
	var userRoles []*entity.SysUserRole
	err := r.Db(ctx).Where("user_id = ?", userId).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

// FindByIds 根据ID列表查询角色
func (r *sysRoleRepo) FindByIds(ctx context.Context, ids []int64) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.Db(ctx).Where("id IN ?", ids).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// DelById 删除角色（包括关联关系）
func (r *sysRoleRepo) DelById(ctx context.Context, id int64) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		if err := r.Db(ctx).Where("role_id = ?", id).Delete(&entity.RolePermissions{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := r.Db(ctx).Where("role_id = ?", id).Delete(&entity.SysUserRole{}).Error; err != nil {
			return err
		}

		// 删除角色
		return r.Db(ctx).Delete(&entity.Role{}, "id = ?", id).Error
	})
}

func (r *sysRoleRepo) FindAllEnabled(ctx context.Context) ([]*entity.Role, error) {
	var list []*entity.Role
	err := r.Db(ctx).Where("status = 1 and deleted_at = 0").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdatePermissions 更新角色权限
func (r *sysRoleRepo) UpdatePermissions(ctx context.Context, roleID int64, permIDs []int64) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 删除原有权限
		if err := r.Db(ctx).Where("role_id = ?", roleID).
			Delete(&entity.RolePermissions{}).Error; err != nil {
			return err
		}

		// 2. 添加新权限
		if len(permIDs) > 0 {
			rolePerms := make([]*entity.RolePermissions, len(permIDs))
			for i, permID := range permIDs {
				rolePerms[i] = &entity.RolePermissions{
					RoleID:       roleID,
					PermissionID: permID,
				}
			}
			if err := r.Db(ctx).Create(&rolePerms).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRoleDataPermission 获取角色数据权限
func (r *sysRoleRepo) GetRoleDataPermission(ctx context.Context, roleID int64) (*entity.DataPermission, error) {
	var dataPerm *entity.DataPermission
	err := r.Db(ctx).Model(&entity.DataPermission{}).
		Where("role_id = ?", roleID).
		First(&dataPerm).Error
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return dataPerm, nil
}

// GetRolePermissions 获取角色权限
func (r *sysRoleRepo) GetRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permissions, error) {
	var permissions []*entity.Permissions
	err := r.Db(ctx).Model(&entity.Permissions{}).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Where("sys_role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// FindByPermissionID 根据权限ID查找角色
func (r *sysRoleRepo) FindByPermissionID(ctx context.Context, permissionID int64) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.Db(ctx).Model(&entity.Role{}).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.role_id = sys_role.id").
		Where("sys_role_permissions.permission_id = ?", permissionID).
		Find(&roles).Error
	return roles, err
}

// FindByType 根据角色类型查询角色列表
func (r *sysRoleRepo) FindByType(ctx context.Context, roleType int8) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.Db(ctx).Where("type = ? AND status = ?", roleType, 1).
		Order("sequence").
		Find(&roles).Error
	return roles, err
}

// GetPermissionsByRoleID 获取角色的权限ID列表
func (r *sysRoleRepo) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var permIDs []int64
	err := r.Db(ctx).Model(&entity.RolePermissions{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permIDs).Error
	return permIDs, err
}

// HasPermission 检查是否有权限
func (r *sysRoleRepo) HasPermission(ctx context.Context, roleID int64, permissionID int64) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.RolePermissions{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error
	return count > 0, err
}

// GetIdsByUserId 获取用户的角色ID列表
func (r *sysRoleRepo) GetIdsByUserId(ctx context.Context, userId string) ([]int64, error) {
	var roleIds []int64
	err := r.Db(ctx).Model(&entity.SysUserRole{}).
		Where("user_id = ?", userId).
		Pluck("role_id", &roleIds).Error
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}

// GetUserCountByRoleID 获取角色关联的用户数量
func (r *sysRoleRepo) GetUserCountByRoleID(ctx context.Context, roleID int64) (int64, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.SysUserRole{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}
