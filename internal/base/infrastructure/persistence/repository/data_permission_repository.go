package repository

import (
	"context"
	"encoding/json"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"strings"
)

type IDataPermissionRepo interface {
	// FindByRoleID 获取角色的数据权限
	FindByRoleID(ctx context.Context, roleID int64) (*entity.DataPermission, error)

	// FindByRoleIDs 批量获取角色的数据权限
	FindByRoleIDs(ctx context.Context, roleIDs []int64) ([]*entity.DataPermission, error)

	// Save 保存数据权限(创建或更新)
	Save(ctx context.Context, perm *entity.DataPermission) error

	// DeleteByRoleID 根据角色ID删除数据权限
	DeleteByRoleID(ctx context.Context, roleID int64) error
	ExistsByRoleID(ctx context.Context, roleID int64) (bool, error)
}

type dataPermissionRepository struct {
	repo IDataPermissionRepo
}

func NewDataPermissionRepository(repo IDataPermissionRepo) repository.IDataPermissionRepository {
	return &dataPermissionRepository{repo: repo}
}

// GetByRoleID 获取角色的数据权限
func (r *dataPermissionRepository) GetByRoleID(ctx context.Context, roleID int64) (*model.DataPermission, error) {
	e, err := r.repo.FindByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}
	return r.toDomain(e)
}

// GetByRoleIDs 批量获取角色的数据权限
func (r *dataPermissionRepository) GetByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.DataPermission, error) {
	entities, err := r.repo.FindByRoleIDs(ctx, roleIDs)
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

// Save 保存数据权限(创建或更新)
func (r *dataPermissionRepository) Save(ctx context.Context, perm *model.DataPermission) error {
	deptIds := ""
	if len(perm.DeptIDs) > 0 {
		deptIds = strings.Join(perm.DeptIDs, ",")
	}

	e := &entity.DataPermission{
		RoleID:   perm.RoleID,
		Scope:    int8(perm.Scope),
		DeptIDs:  deptIds,
		TenantID: perm.TenantID,
	}

	return r.repo.Save(ctx, e)
}

// DeleteByRoleID 根据角色ID删除数据权限
func (r *dataPermissionRepository) DeleteByRoleID(ctx context.Context, roleID int64) error {
	return r.repo.DeleteByRoleID(ctx, roleID)
}

// toDomain 将实体转换为领域模型
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
