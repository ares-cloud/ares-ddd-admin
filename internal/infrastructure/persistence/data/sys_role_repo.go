package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
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
func (r *sysRoleRepo) GetByUserId(ctx context.Context, userId string) ([]*entity.SysUserRole, error) {
	var list []*entity.SysUserRole
	err := r.Db(ctx).Where("user_id = ?", userId).Find(&list).Error
	return list, err
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
