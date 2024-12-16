package data

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

// sysMenuRepo ， 菜单数据层
type sysMenuRepo struct {
	*baserepo.BaseRepo[entity.Permissions, int64]
}

// NewSysMenuRepo ， 菜单数据层工厂方法
// 参数：
//
//	data ： desc
//
// 返回值：
//
//	biz.ISysMenuRepo ：desc
func NewSysMenuRepo(data database.IDataBase) repository.IPermissionsRepo {
	model := new(entity.Permissions)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.RolePermissions{}, &entity.PermissionsResource{}); err != nil {
		hlog.Fatalf("sync sys user tables to db error: %v", err)
	}
	return &sysMenuRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.Permissions, int64](data, entity.Permissions{}),
	}
}
func (r *sysMenuRepo) DelByPermissionsId(ctx context.Context, permissionsId int64) error {
	return r.Db(ctx).Unscoped().Where("permissions_id = ? ", permissionsId).Delete(&entity.PermissionsResource{}).Error
}

func (r *sysMenuRepo) SavePermissionsResource(ctx context.Context, permissionsResource *entity.PermissionsResource) error {
	return r.Db(ctx).Create(permissionsResource).Error
}

func (r *sysMenuRepo) GetByPermissionsId(ctx context.Context, permissionsId int64) ([]*entity.PermissionsResource, error) {
	var permissionsResource []*entity.PermissionsResource
	if err := r.Db(ctx).Find(&permissionsResource, permissionsId).Error; err != nil {
		return nil, err
	}
	return permissionsResource, nil
}
func (r *sysMenuRepo) GetResourceByPermissionsIds(ctx context.Context, permissionsId []int64) ([]*entity.PermissionsResource, error) {
	var permissionsResources []*entity.PermissionsResource
	if err := r.Db(ctx).Where("permissions_id IN ?", permissionsId).Find(&permissionsResources).Error; err != nil {
		return nil, err
	}
	return permissionsResources, nil
}

// GetByCode 根据编码获取权限资源
func (r *sysMenuRepo) GetByCode(ctx context.Context, code string) (*entity.Permissions, error) {
	var perm entity.Permissions

	// 先���询权限基本信息获取ID
	err := r.Db(ctx).Where("code = ?", code).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByRoleID 根据角色ID获取权限列表
func (r *sysMenuRepo) GetByRoleID(ctx context.Context, roleID int64) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
	var rolePerms []*entity.RolePermissions
	var permEntities []*entity.Permissions

	// 查询角色权限关联
	err := r.Db(ctx).Where("role_id = ?", roleID).Find(&rolePerms).Error
	if err != nil {
		return nil, nil, err
	}

	if len(rolePerms) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 获取权限ID列表
	permIDs := make([]int64, 0, len(rolePerms))
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	// 查询权限信息
	err = r.Db(ctx).Where("id IN ?", permIDs).Find(&permEntities).Error
	if err != nil {
		return nil, nil, err
	}

	// 查询权限资源
	resources, err := r.GetResourceByPermissionsIds(ctx, permIDs)
	if err != nil {
		return nil, nil, err
	}

	return permEntities, resources, nil
}

// GetAllTree 获取所有权限树
func (r *sysMenuRepo) GetAllTree(ctx context.Context) ([]*entity.Permissions, []int64, error) {
	var permissions []*entity.Permissions

	// 只查询需要的字段
	err := r.Db(ctx).Select("id, code, name, localize, icon, parent_id").
		Order("sequence desc").
		Find(&permissions).Error
	if err != nil {
		return nil, nil, err
	}

	// 收集所有ID
	ids := make([]int64, 0, len(permissions))
	for _, p := range permissions {
		ids = append(ids, p.ID)
	}

	return permissions, ids, nil
}

// GetTreeByType 根据类型获取权限树
func (r *sysMenuRepo) GetTreeByType(ctx context.Context, permType int8) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
	var permEntities []*entity.Permissions

	// 查询指定类型的权限
	err := r.Db(ctx).Where("type = ?", permType).Order("sequence desc").Find(&permEntities).Error
	if err != nil {
		return nil, nil, err
	}

	if len(permEntities) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 查询权限资源
	resources, err := r.GetResourceByPermissionsIds(ctx, getPermissionIDs(permEntities))
	if err != nil {
		return nil, nil, err
	}

	return permEntities, resources, nil
}

// GetTreeByQuery 根据查询条件获取权限树
func (r *sysMenuRepo) GetTreeByQuery(ctx context.Context, qb *query.QueryBuilder) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
	var permEntities []*entity.Permissions

	db := r.Db(ctx)

	// 添加查询条件
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}

	// 添加排序
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	} else {
		db = db.Order("sequence desc")
	}

	// 执行查询
	err := db.Find(&permEntities).Error
	if err != nil {
		return nil, nil, err
	}

	if len(permEntities) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 查询权限资源
	resources, err := r.GetResourceByPermissionsIds(ctx, getPermissionIDs(permEntities))
	if err != nil {
		return nil, nil, err
	}

	return permEntities, resources, nil
}

// GetTreeByUserAndType 根据用户和类型获取权限树
func (r *sysMenuRepo) GetTreeByUserAndType(ctx context.Context, userID string, permType int8) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
	var rolePerms []*entity.RolePermissions
	var userRoles []*entity.Role

	// 查询用户角色
	err := r.Db(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, nil, err
	}

	if len(userRoles) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 获取角色ID列表
	roleIDs := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.ID)
	}

	// 查询角色权限关联
	err = r.Db(ctx).Where("role_id IN ?", roleIDs).Find(&rolePerms).Error
	if err != nil {
		return nil, nil, err
	}

	if len(rolePerms) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 获取权限ID列表
	permIDs := make([]int64, 0, len(rolePerms))
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	// 查询权限信息
	var permEntities []*entity.Permissions
	err = r.Db(ctx).Where("id IN ? AND type = ?", permIDs, permType).
		Order("sequence desc").Find(&permEntities).Error
	if err != nil {
		return nil, nil, err
	}

	if len(permEntities) == 0 {
		return []*entity.Permissions{}, nil, nil
	}

	// 查询权限资源
	resources, err := r.GetResourceByPermissionsIds(ctx, getPermissionIDs(permEntities))
	if err != nil {
		return nil, nil, err
	}

	return permEntities, resources, nil
}

// 辅助函数：获取权限ID列表
func getPermissionIDs(permissions []*entity.Permissions) []int64 {
	ids := make([]int64, 0, len(permissions))
	for _, p := range permissions {
		ids = append(ids, p.ID)
	}
	return ids
}
