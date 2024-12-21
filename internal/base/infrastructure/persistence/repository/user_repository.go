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

// Count 实现计数查询
func (r *userRepository) Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, qb)
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelById(ctx, id)
}
