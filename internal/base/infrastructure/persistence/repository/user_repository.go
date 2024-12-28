package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

type ISysUserRepo interface {
	baserepo.IBaseRepo[entity.SysUser, string]
	GetByUsername(ctx context.Context, username string) (*entity.SysUser, error)
	DeleteRoleByUserId(ctx context.Context, userId string) error
	BelongsToDepartment(ctx context.Context, userID string, deptID string) bool
}

type userRepository struct {
	repo       ISysUserRepo
	roleRepo   ISysRoleRepo
	mapper     *mapper.UserMapper
	roleMapper *mapper.RoleMapper
}

func NewUserRepository(repo ISysUserRepo, roleRepo ISysRoleRepo) drepository.IUserRepository {
	return &userRepository{
		repo:       repo,
		roleRepo:   roleRepo,
		mapper:     &mapper.UserMapper{},
		roleMapper: &mapper.RoleMapper{},
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	userEntity := r.mapper.ToEntity(user)
	err := r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		userEntity.ID = r.repo.GenStringId()
		_, err := r.repo.Add(ctx, userEntity)
		if err != nil {
			return err
		}
		if len(user.Roles) > 0 {
			// 创建用户角色关联
			for _, role := range user.Roles {
				userRole := &entity.SysUserRole{
					UserID: userEntity.ID,
					RoleID: role.ID,
				}
				if err = r.repo.Db(ctx).Create(userRole).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	userEntity := r.mapper.ToEntity(user)
	err := r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		err := r.repo.EditById(ctx, userEntity.ID, userEntity)
		if err != nil {
			return err
		}
		err = r.repo.DeleteRoleByUserId(ctx, userEntity.ID)
		if err != nil {
			return err
		}
		if len(user.Roles) > 0 {
			// 创建用户角色关联
			userRoles := make([]*entity.SysUserRole, 0, len(user.Roles))
			for _, role := range user.Roles {
				userRoles = append(userRoles, &entity.SysUserRole{
					UserID: userEntity.ID,
					RoleID: role.ID,
				})
			}
			if err = r.repo.Db(ctx).Create(&userRoles).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	// 查询用户基本信息
	userEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	// 查询用户角色关联
	userRoles, err := r.roleRepo.GetByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	roles := make([]*model.Role, 0)
	if len(userRoles) > 0 {
		roleIds := make([]int64, 0)
		for _, role := range userRoles {
			roleIds = append(roleIds, role.RoleID)
		}
		rs, err1 := r.roleRepo.FindByIds(ctx, roleIds)
		if err1 != nil {
			return nil, err1
		}
		roles = r.roleMapper.ToDomainList(rs)
	}
	// 转换为领域模型
	return r.mapper.ToDomain(userEntity, roles), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	// 查询用户基本信息
	userEntity, err := r.repo.GetByUsername(ctx, username)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	// 查询用户角色关联
	userRoles, err := r.roleRepo.GetByUserId(ctx, userEntity.ID)
	if err != nil {
		return nil, err
	}
	roles := make([]*model.Role, 0)
	if len(userRoles) > 0 {
		roleIds := make([]int64, 0)
		for _, role := range userRoles {
			roleIds = append(roleIds, role.RoleID)
		}
		rs, err1 := r.roleRepo.FindByIds(ctx, roleIds)
		if err1 != nil {
			return nil, err1
		}
		roles = r.roleMapper.ToDomainList(rs)
	}
	// 转换为领域模型
	return r.mapper.ToDomain(userEntity, roles), nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, err := r.repo.GetByUsername(ctx, username)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Find 实现动态查询
func (r *userRepository) Find(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error) {
	records, err := r.repo.Find(ctx, qb)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}
	users := r.mapper.ToDomainList(records)
	return users, nil
}

// Count 实现��询
func (r *userRepository) Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, qb)
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelById(ctx, id)
}

// BelongsToDepartment 检查用户是否属于指定部门
func (r *userRepository) BelongsToDepartment(ctx context.Context, userID string, deptID string) bool {
	return r.repo.BelongsToDepartment(ctx, userID, deptID)
}

// FindByDepartment 查询部门下的用户
func (r *userRepository) FindByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *query.QueryBuilder) ([]*model.User, error) {
	var users []*entity.SysUser

	// 构建查询
	db := r.repo.Db(ctx).Model(&entity.SysUser{}).
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

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users), nil
}

// CountByDepartment 统计部门下的用户数量
func (r *userRepository) CountByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *query.QueryBuilder) (int64, error) {
	var count int64

	// 构建查询
	db := r.repo.Db(ctx).Model(&entity.SysUser{}).
		Joins("JOIN sys_user_dept ud ON ud.user_id = sys_user.id").
		Where("ud.dept_id = ?", deptID)

	// 排除管理员
	if excludeAdminID != "" {
		db = db.Where("sys_user.id != ?", excludeAdminID)
	}

	if err := qb.Build(db); err != nil {
		return 0, err
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindUnassignedUsers 查询未分配部门的用户
func (r *userRepository) FindUnassignedUsers(ctx context.Context, qb *query.QueryBuilder) ([]*model.User, error) {
	var users []*entity.SysUser

	// 构建子查询
	subQuery := r.repo.Db(ctx).Model(&entity.UserDepartment{}).
		Select("user_id").
		Group("user_id")

	// 主查询
	db := r.repo.Db(ctx).Model(&entity.SysUser{}).
		Where("id NOT IN (?)", subQuery)

	// 应用查询条件
	if err := qb.Build(db); err != nil {
		return nil, err
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(users), nil
}

// AssignUsers 分配用户到部门(修改为先删除原有关系)
func (r *userRepository) AssignUsers(ctx context.Context, deptID string, userIDs []string) error {
	// 开启事务
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 删除用户原有的部门关系
		if err := r.repo.Db(ctx).Where("user_id IN ?", userIDs).
			Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		// 2. 创建新的部门关系
		userDepts := make([]*entity.UserDepartment, 0, len(userIDs))
		for _, userID := range userIDs {
			userDepts = append(userDepts, &entity.UserDepartment{
				ID:     r.repo.GenStringId(),
				UserID: userID,
				DeptID: deptID,
			})
		}
		return r.repo.Db(ctx).Create(&userDepts).Error
	})
}

// CountUnassignedUsers 统计未分配部门的用户数量
func (r *userRepository) CountUnassignedUsers(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	var count int64

	// 构建子查询
	subQuery := r.repo.Db(ctx).Model(&entity.UserDepartment{}).
		Select("user_id").
		Group("user_id")

	// 主查询
	db := r.repo.Db(ctx).Model(&entity.SysUser{}).
		Where("id NOT IN (?)", subQuery)

	if err := qb.Build(db); err != nil {
		return 0, err
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// TransferUser 调动用户部门
func (r *userRepository) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 删除原部门关系
		if err := r.repo.Db(ctx).Where("user_id = ? AND dept_id = ?", userID, fromDeptID).
			Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		// 2. 创建新部门关系
		userDept := &entity.UserDepartment{
			ID:     r.repo.GenStringId(),
			UserID: userID,
			DeptID: toDeptID,
		}
		return r.repo.Db(ctx).Create(userDept).Error
	})
}
