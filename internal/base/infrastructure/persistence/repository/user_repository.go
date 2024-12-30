package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type ISysUserRepo interface {
	baserepo.IBaseRepo[entity.SysUser, string]
	GetByUsername(ctx context.Context, username string) (*entity.SysUser, error)
	DeleteRoleByUserId(ctx context.Context, userId string) error
	BelongsToDepartment(ctx context.Context, userID string, deptID string) (bool, error)
	GetUserPermissionCodes(ctx context.Context, userID string) ([]string, error)
	GetUserMenus(ctx context.Context, userID string) ([]*entity.Permissions, error)
	FindByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*entity.SysUser, error)
	CountByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error)
	FindUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*entity.SysUser, error)
	CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
	FindByRoleID(ctx context.Context, roleID int64) ([]*entity.SysUser, error)
	AssignUsersToDepartment(ctx context.Context, deptID string, userIDs []string) error
	TransferUserDepartment(ctx context.Context, userID string, fromDeptID string, toDeptID string) error
}

type userRepository struct {
	repo       ISysUserRepo
	roleRepo   ISysRoleRepo
	mapper     *mapper.UserMapper
	roleMapper *mapper.RoleMapper
	menuMapper *mapper.PermissionsMapper
}

func NewUserRepository(repo ISysUserRepo, roleRepo ISysRoleRepo) drepository.IUserRepository {
	return &userRepository{
		repo:       repo,
		roleRepo:   roleRepo,
		mapper:     &mapper.UserMapper{},
		roleMapper: &mapper.RoleMapper{},
		menuMapper: &mapper.PermissionsMapper{},
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
			roleIds = append(roleIds, role.ID)
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
			roleIds = append(roleIds, role.ID)
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

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelById(ctx, id)
}

// BelongsToDepartment 检查用户是否属于指定部门
func (r *userRepository) BelongsToDepartment(ctx context.Context, userID string, deptID string) (bool, error) {
	return r.repo.BelongsToDepartment(ctx, userID, deptID)
}

// TransferUser 调动用户部门
func (r *userRepository) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	return r.repo.TransferUserDepartment(ctx, userID, fromDeptID, toDeptID)
}

// AssignRoles 分配角色给用户
func (r *userRepository) AssignRoles(ctx context.Context, userID string, roleIDs []int64) error {
	return r.repo.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 删除原有角色关联
		if err := r.repo.DeleteRoleByUserId(ctx, userID); err != nil {
			return err
		}

		// 2. 创建新的角色关联
		if len(roleIDs) > 0 {
			userRoles := make([]*entity.SysUserRole, len(roleIDs))
			for i, roleID := range roleIDs {
				userRoles[i] = &entity.SysUserRole{
					UserID: userID,
					RoleID: roleID,
				}
			}
			if err := r.repo.Db(ctx).Create(&userRoles).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
