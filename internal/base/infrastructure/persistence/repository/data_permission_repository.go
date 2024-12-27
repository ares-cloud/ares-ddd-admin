package repository

import (
	"context"
	"encoding/json"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
)

type dataPermissionRepository struct {
	db database.IDataBase
}

func NewDataPermissionRepository(db database.IDataBase) repository.IDataPermissionRepository {
	// 同步表结构
	if err := db.DB(context.Background()).AutoMigrate(&entity.DataPermission{}); err != nil {
		panic(err)
	}
	return &dataPermissionRepository{db: db}
}

func (r *dataPermissionRepository) Create(ctx context.Context, perm *model.DataPermission) error {
	deptIDsBytes, err := json.Marshal(perm.DeptIDs)
	if err != nil {
		return err
	}

	return r.db.DB(ctx).Create(&entity.DataPermission{
		ID:       r.db.GenStringId(),
		RoleID:   perm.RoleID,
		Scope:    int8(perm.Scope),
		DeptIDs:  string(deptIDsBytes),
		TenantID: perm.TenantID,
	}).Error
}

func (r *dataPermissionRepository) Update(ctx context.Context, perm *model.DataPermission) error {
	deptIDsBytes, err := json.Marshal(perm.DeptIDs)
	if err != nil {
		return err
	}

	return r.db.DB(ctx).Model(&entity.DataPermission{}).
		Where("id = ?", perm.ID).
		Updates(map[string]interface{}{
			"scope":    int8(perm.Scope),
			"dept_ids": string(deptIDsBytes),
		}).Error
}

func (r *dataPermissionRepository) Delete(ctx context.Context, id string) error {
	return r.db.DB(ctx).Delete(&entity.DataPermission{}, "id = ?", id).Error
}

func (r *dataPermissionRepository) GetByRoleID(ctx context.Context, roleID int64) (*model.DataPermission, error) {
	var e entity.DataPermission
	err := r.db.DB(ctx).Where("role_id = ?", roleID).First(&e).Error
	if err != nil {
		return nil, err
	}
	return r.toDomain(&e)
}

func (r *dataPermissionRepository) GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.DataPermission, error) {
	var entities []*entity.DataPermission
	err := r.db.DB(ctx).Where("role_id IN ?", roleIDs).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	perms := make([]*model.DataPermission, len(entities))
	for i, e := range entities {
		perm, err := r.toDomain(e)
		if err != nil {
			return nil, err
		}
		perms[i] = perm
	}
	return perms, nil
}

func (r *dataPermissionRepository) toDomain(e *entity.DataPermission) (*model.DataPermission, error) {
	var deptIDs []string
	if err := json.Unmarshal([]byte(e.DeptIDs), &deptIDs); err != nil {
		return nil, err
	}

	return &model.DataPermission{
		ID:       e.ID,
		RoleID:   e.RoleID,
		Scope:    model.DataScope(e.Scope),
		DeptIDs:  deptIDs,
		TenantID: e.TenantID,
	}, nil
}
