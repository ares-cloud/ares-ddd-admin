package data

import (
	"context"
	"fmt"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type sysTenantRepo struct {
	*baserepo.BaseRepo[entity.Tenant, string]
}

func NewSysTenantRepo(data database.IDataBase) repository.ISysTenantRepo {
	model := new(entity.Tenant)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.TenantPermissions{}); err != nil {
		hlog.Fatalf("sync tenant tables to db error: %v", err)
	}
	return &sysTenantRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.Tenant, string](data, entity.Tenant{}),
	}
}

// GetByCode 根据编码获取租户
func (r *sysTenantRepo) GetByCode(ctx context.Context, code string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.Db(ctx).Where("code = ?", code).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// DelById 删除租户（包括关联关系）
func (r *sysTenantRepo) DelById(ctx context.Context, id string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 检查是否为默认租户
		tenant, err := r.FindById(ctx, id)
		if err != nil {
			return fmt.Errorf("find tenant failed: %w", err)
		}
		if tenant.IsDefault == 1 {
			return fmt.Errorf("cannot delete default tenant")
		}

		// 删除租户下的所有用户
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.SysUser{}).Error; err != nil {
			return err
		}

		// 删除租户下的所有角色
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.Role{}).Error; err != nil {
			return err
		}

		// 删除租户
		return r.Db(ctx).Delete(&entity.Tenant{}, "id = ?", id).Error
	})
}

// Create 创建租户（重写基类方法，处理默认租户）
func (r *sysTenantRepo) Create(ctx context.Context, tenant *entity.Tenant) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		tenant.CreatedAt = time.Now().Unix()
		// 如果是默认租户，需要将其他租户设置为非默认
		if tenant.IsDefault == 1 {
			if err := r.Db(ctx).Model(&entity.Tenant{}).Where("is_default = ?", 1).
				Updates(map[string]interface{}{"is_default": 2}).Error; err != nil {
				return err
			}
		}

		// 创建租户
		return r.Db(ctx).Create(tenant).Error
	})
}

// Update 更新租户（重写基类方法，处理默认租户）
func (r *sysTenantRepo) Update(ctx context.Context, tenant *entity.Tenant) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		tenant.UpdatedAt = time.Now().Unix()
		// 如果是默认租户，需要将其他租户设置为非默认
		if tenant.IsDefault == 1 {
			if err := r.Db(ctx).Model(&entity.Tenant{}).
				Where("id != ? AND is_default = ?", tenant.ID, 1).
				Updates(map[string]interface{}{"is_default": 2}).Error; err != nil {
				return err
			}
		}

		// 更新租户
		return r.Db(ctx).Updates(tenant).Error
	})
}

// AssignPermissions 分配权限给租户
func (r *sysTenantRepo) AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 删除原有权限关联
		if err := r.Db(ctx).Where("tenant_id = ?", tenantID).Delete(&entity.TenantPermissions{}).Error; err != nil {
			return err
		}

		// 创建新的权限关联
		for _, permID := range permissionIDs {
			tp := &entity.TenantPermissions{
				TenantID:     tenantID,
				PermissionID: permID,
			}
			if err := r.Db(ctx).Create(tp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetPermissionsByTenantID 获取租户的权限列表及资源
func (r *sysTenantRepo) GetPermissionsByTenantID(ctx context.Context, tenantID string) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
	var permissions []*entity.Permissions

	// 通过关联表查询权限
	err := r.Db(ctx).
		Joins("JOIN sys_tenant_permissions tp ON tp.permission_id = sys_permissions.id").
		Where("tp.tenant_id = ?", tenantID).
		Find(&permissions).Error
	if err != nil {
		return nil, nil, err
	}

	if len(permissions) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 获取权限ID列表
	permIDs := make([]int64, 0, len(permissions))
	for _, p := range permissions {
		permIDs = append(permIDs, p.ID)
	}

	// 查询权限资源
	var resources []*entity.PermissionsResource
	err = r.Db(ctx).Where("permissions_id IN ?", permIDs).Find(&resources).Error
	if err != nil {
		return nil, nil, err
	}

	return permissions, resources, nil
}

// HasPermission 检查租户是否拥有指定权限
func (r *sysTenantRepo) HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.TenantPermissions{}).
		Where("tenant_id = ? AND permission_id = ?", tenantID, permissionID).
		Count(&count).Error
	return count > 0, err
}
