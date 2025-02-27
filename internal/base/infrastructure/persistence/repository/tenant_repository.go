package repository

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type ISysTenantRepo interface {
	baserepo.IBaseRepo[entity.Tenant, string]
	Update(ctx context.Context, tenant *entity.Tenant) error
	DeleteWithRelations(ctx context.Context, id string) error // 删除租户及关联数据

	// 权限相关
	AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error
	GetPermissionsByTenantID(ctx context.Context, tenantID string) ([]*entity.Permissions, error)
	GetTenantIDPermissionsByType(ctx context.Context, tenantID string, int8 int64) ([]*entity.Permissions, error)
	HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error)

	Lock(ctx context.Context, tenantID string, reason string) error
	Unlock(ctx context.Context, tenantID string) error
	GetTenantRoles(ctx context.Context, tenantID string) ([]*entity.Role, error)
}

type tenantRepository struct {
	repo       ISysTenantRepo
	userRepo   ISysUserRepo
	mapper     *mapper.TenantMapper
	userMapper *mapper.UserMapper
	permMapper *mapper.PermissionsMapper
}

func NewTenantRepository(repo ISysTenantRepo, userRepo ISysUserRepo) drepository.ITenantRepository {
	userMapper := &mapper.UserMapper{}
	permMapper := &mapper.PermissionsMapper{}
	return &tenantRepository{
		repo:       repo,
		userRepo:   userRepo,
		mapper:     mapper.NewTenantMapper(userMapper),
		userMapper: userMapper,
		permMapper: permMapper,
	}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	tenantEntity := r.mapper.ToEntity(tenant)
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 生成ID
		tenantEntity.ID = r.repo.GenStringId()

		// 创建管理员用户
		if tenant.AdminUser != nil {
			userEntity := r.userMapper.ToEntity(tenant.AdminUser)
			userEntity.ID = r.userRepo.GenStringId()
			userEntity.TenantID = tenantEntity.ID
			if _, err := r.userRepo.Add(ctx, userEntity); err != nil {
				return fmt.Errorf("create admin user failed: %w", err)
			}
			tenantEntity.AdminUserID = userEntity.ID
		}

		// 创建租户
		if _, err := r.repo.Add(ctx, tenantEntity); err != nil {
			return fmt.Errorf("create tenant failed: %w", err)
		}
		return nil
	})
}

func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	tenantEntity := r.mapper.ToEntity(tenant)
	if err := r.repo.Update(ctx, tenantEntity); err != nil {
		return fmt.Errorf("update tenant failed: %w", err)
	}
	return nil
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	if err := r.repo.DeleteWithRelations(ctx, id); err != nil {
		return fmt.Errorf("delete tenant failed: %w", err)
	}
	return nil
}

func (r *tenantRepository) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenantEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// 查询管理员用户
	adminUser, err := r.userRepo.FindById(ctx, tenantEntity.AdminUserID)
	if err != nil && !database.IfErrorNotFound(err) {
		return nil, err
	}

	return r.mapper.ToDomain(tenantEntity, adminUser), nil
}

func (r *tenantRepository) FindByCode(ctx context.Context, code string) (*model.Tenant, error) {
	qb := db_query.NewQueryBuilder()
	qb.Where("code", db_query.Eq, code)

	tenants, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, nil
	}

	// 查询管理员用户
	adminUser, err := r.userRepo.FindById(ctx, tenants[0].AdminUserID)
	if err != nil && !database.IfErrorNotFound(err) {
		return nil, err
	}

	return r.mapper.ToDomain(tenants[0], adminUser), nil
}

func (r *tenantRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	qb := db_query.NewQueryBuilder()
	qb.Where("code", db_query.Eq, code)

	count, err := r.repo.Count(ctx, qb)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tenantRepository) AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error {
	return r.repo.AssignPermissions(ctx, tenantID, permissionIDs)
}

func (r *tenantRepository) GetPermissions(ctx context.Context, tenantID string) ([]*model.Permissions, error) {
	permissions, err := r.repo.GetPermissionsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return r.permMapper.ToDomainList(permissions, nil), nil
}

func (r *tenantRepository) HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error) {
	return r.repo.HasPermission(ctx, tenantID, permissionID)
}
func (r *tenantRepository) Lock(ctx context.Context, tenantID string, reason string) error {
	return r.repo.Lock(ctx, tenantID, reason)
}

func (r *tenantRepository) Unlock(ctx context.Context, tenantID string) error {
	return r.repo.Unlock(ctx, tenantID)
}
