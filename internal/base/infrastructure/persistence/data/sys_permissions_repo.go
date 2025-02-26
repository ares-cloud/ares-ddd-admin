package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
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
	err := r.Db(ctx).
		//Select("id, code, name, localize, icon, parent_id").
		Where("type = 1").
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
func (r *sysMenuRepo) GetTreeByQuery(ctx context.Context, qb *db_query.QueryBuilder) ([]*entity.Permissions, []*entity.PermissionsResource, error) {
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
	var userRoles []*entity.SysUserRole

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
		roleIDs = append(roleIDs, ur.RoleID)
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
func (r *sysMenuRepo) GetResourcesByRoles(ctx context.Context, roles []int64) ([]*entity.PermissionsResource, error) {
	var resources []*entity.PermissionsResource

	// 1. 先查询角色对应的权限ID
	var permissionIDs []int64
	err := r.Db(ctx).Model(&entity.RolePermissions{}).
		Joins("JOIN sys_role ON sys_role.id = sys_role_permissions.role_id").
		Where("sys_role.code IN ? AND sys_role.status = ?", roles, 1).
		Pluck("permission_id", &permissionIDs).Error
	if err != nil {
		return nil, err
	}

	if len(permissionIDs) == 0 {
		return []*entity.PermissionsResource{}, nil
	}

	// 2. 查询启用的权限资源
	err = r.Db(ctx).Model(&entity.PermissionsResource{}).
		Joins("JOIN sys_permissions ON sys_permissions.id = sys_permissions_resource.permissions_id").
		Where("sys_permissions.id IN ? AND sys_permissions.status = ? AND sys_permissions.type = ?",
			permissionIDs, 1, 3). // type=3 表示API类型的权限
		Find(&resources).Error
	if err != nil {
		return nil, err
	}

	return resources, nil
}
func (r *sysMenuRepo) GetByRoles(ctx context.Context, roles []int64) ([]*entity.Permissions, error) {
	var permissions []*entity.Permissions

	// 1. 查询角色关联的权限ID
	var permissionIDs []int64
	err := r.Db(ctx).Model(&entity.RolePermissions{}).
		Where("role_id IN ?", roles).
		Pluck("permission_id", &permissionIDs).Error
	if err != nil {
		return nil, err
	}

	if len(permissionIDs) == 0 {
		return []*entity.Permissions{}, nil
	}

	// 2. 查询启用的权限
	err = r.Db(ctx).Where("id IN ? AND status = ? ",
		permissionIDs, 1). // status=1 表示启用
		Order("sequence desc").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// 辅助函数：获取权限ID列表
func getPermissionIDs(permissions []*entity.Permissions) []int64 {
	ids := make([]int64, 0, len(permissions))
	for _, p := range permissions {
		ids = append(ids, p.ID)
	}
	return ids
}

// GetResourcesByRolesGrouped 根据角色ID获取分组的权限资源
func (r *sysMenuRepo) GetResourcesByRolesGrouped(ctx context.Context, roles []int64) (map[int64][]*entity.PermissionsResource, error) {
	// 结果map: roleID -> resources
	resourceMap := make(map[int64][]*entity.PermissionsResource)

	// 1. 先查询角色权限关联和权限资源
	var results []struct {
		RoleID       int64  `gorm:"column:role_id"`
		Method       string `gorm:"column:method"`
		Path         string `gorm:"column:path"`
		PermissionID int64  `gorm:"column:permissions_id"`
	}

	err := r.Db(ctx).Table("sys_role_permissions").
		Select("sys_role_permissions.role_id, pr.method, pr.path, pr.permissions_id").
		Joins("JOIN sys_permissions p ON p.id = sys_role_permissions.permission_id").
		Joins("JOIN sys_permissions_resource pr ON pr.permissions_id = p.id").
		Where("sys_role_permissions.role_id IN ? AND p.status = ?", roles, 1).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 2. 组织数据到map中
	for _, result := range results {
		resource := &entity.PermissionsResource{
			PermissionsID: result.PermissionID,
			Method:        result.Method,
			Path:          result.Path,
		}
		resourceMap[result.RoleID] = append(resourceMap[result.RoleID], resource)
	}

	// 3. 确保所有角色都有对应的切片(即使是空的)
	for _, roleID := range roles {
		if _, exists := resourceMap[roleID]; !exists {
			resourceMap[roleID] = []*entity.PermissionsResource{}
		}
	}

	return resourceMap, nil
}

func (r *sysMenuRepo) ExistsById(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.Permissions{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
