package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"gorm.io/gorm"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type UserQueryService struct {
	userRepo             repository.ISysUserRepo
	roleRepo             repository.ISysRoleRepo
	permissionsRepo      repository.IPermissionsRepo
	userConverter        *converter.UserConverter
	roleConverter        *converter.RoleConverter
	permissionsConverter *converter.PermissionsConverter
	deptRepo             repository.ISysDepartmentRepo
	deptConverter        *converter.DepartmentConverter
	conf                 *configs.Bootstrap
}

func NewUserQueryService(
	userRepo repository.ISysUserRepo,
	roleRepo repository.ISysRoleRepo,
	permissionsRepo repository.IPermissionsRepo,
	userConverter *converter.UserConverter,
	roleConverter *converter.RoleConverter,
	permissionsConverter *converter.PermissionsConverter,
	deptRepo repository.ISysDepartmentRepo,
	deptConverter *converter.DepartmentConverter,
	conf *configs.Bootstrap,
) *UserQueryService {
	return &UserQueryService{
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		permissionsRepo:      permissionsRepo,
		userConverter:        userConverter,
		roleConverter:        roleConverter,
		permissionsConverter: permissionsConverter,
		deptRepo:             deptRepo,
		deptConverter:        deptConverter,
		conf:                 conf,
	}
}
func (u *UserQueryService) GetSuperAdmin(_ context.Context) (*dto.UserInfoDto, error) {
	userId := constant.RoleSuperAdmin
	return &dto.UserInfoDto{
		User: &dto.UserDto{
			ID:       userId,
			Username: u.conf.SuperAdmin.Phone,
			Nickname: u.conf.SuperAdmin.Nickname,
			Phone:    u.conf.SuperAdmin.Phone,
			Email:    "",
			Avatar:   "",
			Status:   1,
			TenantID: "",
		},
		Roles:       []string{userId},
		HomePage:    "User",
		Permissions: []string{"*"},
	}, nil
}

// GetUser 获取用户详情
func (u *UserQueryService) GetUser(ctx context.Context, id string) (*dto.UserDto, error) {
	// 1. 获取用户基本信息
	user, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// 2. 获取用户角色ID列表
	roleIds, err := u.roleRepo.GetIdsByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	return u.userConverter.ToDTO(user, roleIds), nil
}

// FindUsers 查询用户列表
func (u *UserQueryService) FindUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	// 1. 获取用户列表
	users, err := u.userRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户角色ID列表
	userDtos := make([]*dto.UserDto, 0, len(users))
	for _, user := range users {
		roleIds, err := u.roleRepo.GetIdsByUserId(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if userDto := u.userConverter.ToDTO(user, roleIds); userDto != nil {
			userDtos = append(userDtos, userDto)
		}
	}

	return userDtos, nil
}

// GetUserRolesCode 获取用户角色编码列表
func (u *UserQueryService) GetUserRolesCode(ctx context.Context, userID string) ([]string, error) {
	// 1. 获取用户角色
	roles, err := u.roleRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 提取角色编码
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Status == 1 { // 只返回启用状态的角色
			roleCodes = append(roleCodes, role.Code)
		}
	}

	return roleCodes, nil
}

// CountUsers 统计用户数量
func (u *UserQueryService) CountUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return u.userRepo.Count(ctx, qb)
}

// GetUserPermissions 获取用户权限
func (u *UserQueryService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	roles, err := u.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取角色对应的权限
	permissions := make([]string, 0)
	for _, role := range roles {
		perms, err := u.roleRepo.GetRolePermissions(ctx, role.RoleID)
		if err != nil {
			return nil, err
		}
		for _, p := range perms {
			permissions = append(permissions, p.Code)
		}
	}
	return permissions, nil
}

// GetUserRoles 获取用户角色
func (u *UserQueryService) GetUserRoles(ctx context.Context, userID string) ([]*dto.RoleDto, error) {
	roles, err := u.roleRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.roleConverter.ToDTOList(roles), nil
}

// GetUserMenus 获取用户菜单
func (u *UserQueryService) GetUserMenus(ctx context.Context, userID string) ([]*dto.PermissionsDto, error) {
	roles, err := u.roleRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取角色对应的菜单权限
	var permissions []*entity.Permissions
	for _, role := range roles {
		perms, err := u.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perms...)
	}

	// 转换为DTO并返回
	return u.permissionsConverter.ToDTOList(permissions), nil
}

// GetUserTreeMenus 获取用户菜单树
func (u *UserQueryService) GetUserTreeMenus(ctx context.Context, userID string) ([]*dto.PermissionsTreeDto, error) {
	var permissions []*entity.Permissions
	var err error
	if actx.IsSuperAdmin(ctx) {
		permissions, _, err = u.permissionsRepo.GetAllTree(ctx)
	} else {
		// 1. 获取用户角色
		roles, err := u.roleRepo.GetIdsByUserId(ctx, userID)
		if err != nil {
			return nil, err
		}

		if len(roles) == 0 {
			return []*dto.PermissionsTreeDto{}, nil
		}

		// 2. 获取角色对应的菜单权限
		permissions, _, err = u.permissionsRepo.GetTreeByUserAndType(context.Background(), userID, 1) // type=1表示菜单类型
	}
	if err != nil {
		return nil, err
	}
	// 3. 转换为DTO
	return u.permissionsConverter.ToSimpleTreeDTOList(permissions), nil
}

// FindUsersByDepartment 查询部门下的用户
func (u *UserQueryService) FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	users, err := u.userRepo.FindByDepartment(ctx, deptID, excludeAdminID, qb)
	if err != nil {
		return nil, err
	}

	// 获取用户角色并转换
	userDtos := make([]*dto.UserDto, 0, len(users))
	for _, user := range users {
		roleIds, err := u.roleRepo.GetIdsByUserId(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if userDto := u.userConverter.ToDTO(user, roleIds); userDto != nil {
			userDtos = append(userDtos, userDto)
		}
	}

	return userDtos, nil
}

// CountUsersByDepartment 统计部门用户数量
func (u *UserQueryService) CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	if deptID == "" {
		return 0, fmt.Errorf("部门ID不能为空")
	}
	return u.userRepo.CountByDepartment(ctx, deptID, excludeAdminID, qb)
}
func (u *UserQueryService) GetUserDepartments(ctx context.Context, userID string) ([]*dto.DepartmentDto, error) {
	departments, err := u.deptRepo.GetDeptByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*dto.DepartmentDto{}, nil
		}
		return nil, err
	}
	return u.deptConverter.ToDTOList(departments), nil
}

// FindUnassignedUsers 查询未分配部门的用户
func (u *UserQueryService) FindUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	users, err := u.userRepo.FindUnassignedUsers(ctx, qb)
	if err != nil {
		return nil, err
	}
	return u.userConverter.ToDTOList(users), nil
}

// CountUnassignedUsers 统计未分配部门的用户数量
func (u *UserQueryService) CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return u.userRepo.CountUnassignedUsers(ctx, qb)
}
