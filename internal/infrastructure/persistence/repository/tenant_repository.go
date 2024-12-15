package repository

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ISysTenantRepo interface {
	baserepo.IBaseRepo[entity.Tenant, string]
	GetByCode(ctx context.Context, code string) (*entity.Tenant, error)
	DelById(ctx context.Context, id string) error
	Create(ctx context.Context, tenant *entity.Tenant) error
	Update(ctx context.Context, tenant *entity.Tenant) error

	// 权限相关方法
	AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error
	GetPermissionsByTenantID(ctx context.Context, tenantID string) ([]*entity.Permissions, []*entity.PermissionsResource, error)
	HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error)

	// 用户相关方法
	FindByIds(ctx context.Context, ids []string) ([]*entity.Tenant, error)
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
		err := r.repo.Create(ctx, tenantEntity)
		if err != nil {
			return fmt.Errorf("create tenant failed: %w", err)
		}

		return nil
	})
}

func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	tenantEntity := r.mapper.ToEntity(tenant)
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 更新租户信息
		err := r.repo.Update(ctx, tenantEntity)
		if err != nil {
			return fmt.Errorf("update tenant failed: %w", err)
		}

		return nil
	})
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelById(ctx, id)
}

func (r *tenantRepository) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenantEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
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
	tenantEntity, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
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

func (r *tenantRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	_, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *tenantRepository) Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.Tenant, error) {
	tenants, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	// 获取所有管理员用户ID
	adminUserIDs := make([]string, 0, len(tenants))
	for _, t := range tenants {
		if t.AdminUserID != "" {
			adminUserIDs = append(adminUserIDs, t.AdminUserID)
		}
	}

	// 查询管理员用户
	var adminUsers []*entity.SysUser
	if len(adminUserIDs) > 0 {
		if adminUsers, err = r.userRepo.FindByIds(ctx, adminUserIDs); err != nil {
			return nil, err
		}
	}

	// 构建管理员用户映射
	adminUserMap := make(map[string]*entity.SysUser)
	for _, u := range adminUsers {
		adminUserMap[u.ID] = u
	}

	// 转换为领域模型
	result := make([]*model.Tenant, len(tenants))
	for i, t := range tenants {
		result[i] = r.mapper.ToDomain(t, adminUserMap[t.AdminUserID])
	}

	return result, nil
}

func (r *tenantRepository) Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, qb)
}
func (r *tenantRepository) AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error {
	return r.repo.AssignPermissions(ctx, tenantID, permissionIDs)
}

func (r *tenantRepository) GetPermissions(ctx context.Context, tenantID string) ([]*model.Permissions, error) {
	// 获取租户权限和资源
	permissions, resources, err := r.repo.GetPermissionsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("get tenant permissions failed: %w", err)
	}

	// 转换为领域模型
	return r.permMapper.ToDomainList(permissions, resources), nil
}

func (r *tenantRepository) HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error) {
	return r.repo.HasPermission(ctx, tenantID, permissionID)
}
