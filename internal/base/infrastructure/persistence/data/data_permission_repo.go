package data

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

type dataPermissionRepo struct {
	db database.IDataBase
}

func NewDataPermissionRepo(db database.IDataBase) repository.IDataPermissionRepo {
	return &dataPermissionRepo{db: db}
}

func (d *dataPermissionRepo) FindByRoleID(ctx context.Context, roleID int64) (*entity.DataPermission, error) {
	var e entity.DataPermission
	err := d.db.DB(ctx).Where("role_id = ?", roleID).First(&e).Error
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (d *dataPermissionRepo) FindByRoleIDs(ctx context.Context, roleIDs []int64) ([]*entity.DataPermission, error) {
	var entities []*entity.DataPermission
	err := d.db.DB(ctx).Where("role_id IN ?", roleIDs).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *dataPermissionRepo) Save(ctx context.Context, e *entity.DataPermission) error {
	// 查询是否已存在
	exists, err := d.ExistsByRoleID(ctx, e.RoleID)
	if err != nil {
		return err
	}

	if exists {
		// 更新
		return d.db.DB(ctx).Model(&entity.DataPermission{}).
			Where("role_id = ?", e.RoleID).
			Updates(map[string]interface{}{
				"scope":    e.Scope,
				"dept_ids": e.DeptIDs,
			}).Error
	}

	// 创建
	return d.db.DB(ctx).Create(e).Error
}

func (d *dataPermissionRepo) DeleteByRoleID(ctx context.Context, roleID int64) error {
	return d.db.DB(ctx).Where("role_id = ?", roleID).Delete(&entity.DataPermission{}).Error
}

func (d *dataPermissionRepo) ExistsByRoleID(ctx context.Context, roleID int64) (bool, error) {
	var count int64
	err := d.db.DB(ctx).Model(&entity.DataPermission{}).Where("role_id = ?", roleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
