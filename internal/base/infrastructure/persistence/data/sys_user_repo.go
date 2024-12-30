package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// sysUserRepo ， 用户数据层
type sysUserRepo struct {
	*baserepo.BaseRepo[entity.SysUser, string]
}

// NewSysUserRepo ， 用户数据层工厂方法
// 参数：
//
//	data ： desc
//
// 返回值：
//
//	biz.ISysUserRepo ：desc
func NewSysUserRepo(data database.IDataBase) repository.ISysUserRepo {
	model := new(entity.SysUser)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.SysUserRole{}); err != nil {
		hlog.Fatalf("sync sys user tables to db error: %v", err)
	}
	return &sysUserRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.SysUser, string](data, entity.SysUser{}),
	}
}
func (r *sysUserRepo) GetByUsername(ctx context.Context, username string) (*entity.SysUser, error) {
	var result *entity.SysUser
	err := r.Db(ctx).Where("username = ?", username).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *sysUserRepo) DeleteRoleByUserId(ctx context.Context, userId string) error {
	return r.Db(ctx).Where("user_id = ?", userId).Delete(&entity.SysUserRole{}).Error
}
func (r *sysUserRepo) BelongsToDepartment(ctx context.Context, userID string, deptID string) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.UserDepartment{}).
		Where("user_id = ? AND dept_id = ?", userID, deptID).
		Count(&count).Error
	return count > 0, err
}

// GetUserPermissionCodes 获取用户权限代码列表
func (r *sysUserRepo) GetUserPermissionCodes(ctx context.Context, userID string) ([]string, error) {
	var codes []string
	err := r.Db(ctx).Model(&entity.Permissions{}).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permissions.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Pluck("sys_permissions.code", &codes).Error
	return codes, err
}

// GetUserMenus 获取用户菜单权限
func (r *sysUserRepo) GetUserMenus(ctx context.Context, userID string) ([]*entity.Permissions, error) {
	var permissions []*entity.Permissions
	err := r.Db(ctx).Model(&entity.Permissions{}).
		Joins("JOIN sys_role_permissions ON sys_role_permissions.permission_id = sys_permissions.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permissions.role_id").
		Where("sys_user_role.user_id = ? AND sys_permissions.type = ?", userID, 1).
		Find(&permissions).Error
	return permissions, err
}

// FindByDepartment 查询部门下的用户
func (r *sysUserRepo) FindByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*entity.SysUser, error) {
	// 构建查询
	db := r.Db(ctx).Model(&entity.SysUser{}).
		Joins("JOIN sys_user_dept ud ON ud.user_id = sys_user.id").
		Where("ud.dept_id = ?", deptID)

	// 排除管理员
	if excludeAdminID != "" {
		db = db.Where("sys_user.id != ?", excludeAdminID)
	}

	// 应用查询条件
	if err := qb.Build(db); err != nil {
		return nil, err
	}

	var users []*entity.SysUser
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CountByDepartment 统计部门下的用户数量
func (r *sysUserRepo) CountByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	var count int64
	db := r.Db(ctx).Model(&entity.SysUser{}).
		Joins("JOIN sys_user_dept ud ON ud.user_id = sys_user.id").
		Where("ud.dept_id = ?", deptID)

	if excludeAdminID != "" {
		db = db.Where("sys_user.id != ?", excludeAdminID)
	}

	if err := qb.Build(db); err != nil {
		return 0, err
	}

	err := db.Count(&count).Error
	return count, err
}

// FindUnassignedUsers 查询未分配部门的用户
func (r *sysUserRepo) FindUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*entity.SysUser, error) {
	subQuery := r.Db(ctx).Model(&entity.UserDepartment{}).
		Select("user_id").
		Group("user_id")

	db := r.Db(ctx).Model(&entity.SysUser{}).
		Where("id NOT IN (?)", subQuery)

	if err := qb.Build(db); err != nil {
		return nil, err
	}

	var users []*entity.SysUser
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CountUnassignedUsers 统计未分配部门的用户数量
func (r *sysUserRepo) CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	var count int64
	subQuery := r.Db(ctx).Model(&entity.UserDepartment{}).
		Select("user_id").
		Group("user_id")

	db := r.Db(ctx).Model(&entity.SysUser{}).
		Where("id NOT IN (?)", subQuery)

	if err := qb.Build(db); err != nil {
		return 0, err
	}

	err := db.Count(&count).Error
	return count, err
}

// FindByRoleID 根据角色ID查找用户
func (r *sysUserRepo) FindByRoleID(ctx context.Context, roleID int64) ([]*entity.SysUser, error) {
	var users []*entity.SysUser
	err := r.Db(ctx).Model(&entity.SysUser{}).
		Joins("JOIN sys_user_role ON sys_user_role.user_id = sys_user.id").
		Where("sys_user_role.role_id = ?", roleID).
		Find(&users).Error
	return users, err
}

// AssignUsersToDepartment 分配用户到部门
func (r *sysUserRepo) AssignUsersToDepartment(ctx context.Context, deptID string, userIDs []string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		if err := r.Db(ctx).Where("user_id IN ?", userIDs).
			Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		userDepts := make([]*entity.UserDepartment, len(userIDs))
		for i, userID := range userIDs {
			userDepts[i] = &entity.UserDepartment{
				ID:     r.GenInt64Id(),
				UserID: userID,
				DeptID: deptID,
			}
		}
		return r.Db(ctx).Create(&userDepts).Error
	})
}

// TransferUserDepartment 调动用户部门
func (r *sysUserRepo) TransferUserDepartment(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 删除原部门关联
		if err := r.Db(ctx).Where("user_id = ? AND dept_id = ?", userID, fromDeptID).
			Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		// 2. 创建新部门关联
		userDept := &entity.UserDepartment{
			UserID: userID,
			DeptID: toDeptID,
		}
		if err := r.Db(ctx).Create(userDept).Error; err != nil {
			return err
		}

		// 3. 更新用户主部门信息
		if err := r.Db(ctx).Model(&entity.SysUser{}).
			Where("id = ?", userID).
			Update("dept_id", toDeptID).Error; err != nil {
			return err
		}

		return nil
	})
}
