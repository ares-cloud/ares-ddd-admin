package impl

import (
	"context"
	"fmt"
	"sort"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"

	dquery "github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
)

type UserQueryService struct {
	userRepo   repository.ISysUserRepo
	roleRepo   repository.ISysRoleRepo
	userMapper *mapper.UserMapper
	roleMapper *mapper.RoleMapper
	permMapper *mapper.PermissionsMapper
}

func NewUserQueryService(
	userRepo repository.ISysUserRepo,
	roleRepo repository.ISysRoleRepo,
	userMapper *mapper.UserMapper,
	roleMapper *mapper.RoleMapper,
	permMapper *mapper.PermissionsMapper,
) *UserQueryService {
	return &UserQueryService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		userMapper: userMapper,
		roleMapper: roleMapper,
		permMapper: permMapper,
	}
}

// GetUser 获取用户信息
func (u UserQueryService) GetUser(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 1. 获取用户基本信息
	user, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 2. 获取用户角色
	roles, err := u.GetUserRoles(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 3. 转换为领域模型
	return u.userMapper.ToDomain(user, roles), nil
}

// FindUsers 查询用户列表
func (u UserQueryService) FindUsers(ctx context.Context, qb *dquery.QueryBuilder) ([]*model.User, error) {
	users, err := u.userRepo.Find(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	return u.userMapper.ToDomainList(users), nil
}
func (u UserQueryService) GetUserRolesCode(ctx context.Context, userID string) ([]string, error) {
	// 在缓存层实现
	return nil, nil
}

// CountUsers 统计用户数量
func (u UserQueryService) CountUsers(ctx context.Context, qb *dquery.QueryBuilder) (int64, error) {
	return u.userRepo.Count(ctx, qb)
}

// GetUserPermissions 获取用户权限
func (u UserQueryService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	permissions, err := u.userRepo.GetUserPermissionCodes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户权限失败: %w", err)
	}

	if permissions == nil {
		permissions = make([]string, 0)
	}
	return permissions, nil
}

// GetUserRoles 获取用户角色
func (u UserQueryService) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	roles, err := u.roleRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	if roles == nil {
		return make([]*model.Role, 0), nil
	}
	return u.roleMapper.ToDomainList(roles), nil
}

// GetUserMenus 获取用户菜单
func (u UserQueryService) GetUserMenus(ctx context.Context, userID string) ([]*model.Permissions, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	menus, err := u.userRepo.GetUserMenus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户菜单失败: %w", err)
	}

	if menus == nil {
		return make([]*model.Permissions, 0), nil
	}
	return u.permMapper.ToDomainList(menus, nil), nil
}

// GetUserTreeMenus 获取用户菜单树
func (u UserQueryService) GetUserTreeMenus(ctx context.Context, userID string) ([]*model.Permissions, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 1. 获取用户的所有菜单权限
	menus, err := u.userRepo.GetUserMenus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户菜单失败: %w", err)
	}

	if menus == nil {
		return make([]*model.Permissions, 0), nil
	}

	// 2. 转换为领域模型
	menuList := u.permMapper.ToDomainList(menus, nil)

	// 3. 构建树形结构
	return buildPermissionTree(menuList), nil
}

// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []*model.Permissions) []*model.Permissions {
	// 1. 创建一个map用于快速查找
	permMap := make(map[int64]*model.Permissions)
	for _, perm := range permissions {
		permMap[perm.ID] = perm
		// 初始化子节点切片
		perm.Children = make([]*model.Permissions, 0)
	}

	// 2. 构建树形结构
	var roots []*model.Permissions
	for _, perm := range permissions {
		if perm.ParentID == 0 { // 根节点
			roots = append(roots, perm)
		} else {
			if parent, ok := permMap[perm.ParentID]; ok {
				parent.Children = append(parent.Children, perm)
			}
		}
	}

	// 3. 对每个节点的子节点进行排序
	for _, perm := range permissions {
		if len(perm.Children) > 0 {
			sortPermissions(perm.Children)
		}
	}

	// 4. 对根节点进行排序
	sortPermissions(roots)

	return roots
}

// sortPermissions 根据Sort字段对权限列表进行排序
func sortPermissions(perms []*model.Permissions) {
	sort.Slice(perms, func(i, j int) bool {
		// 如果Sort相同，则按照ID排序
		if perms[i].Sequence == perms[j].Sequence {
			return perms[i].ID < perms[j].ID
		}
		return perms[i].Sequence < perms[j].Sequence
	})
}

// FindUsersByDepartment 查询部门用户
func (u UserQueryService) FindUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *dquery.QueryBuilder) ([]*model.User, error) {
	if deptID == "" {
		return nil, fmt.Errorf("部门ID不能为空")
	}

	users, err := u.userRepo.FindByDepartment(ctx, deptID, excludeAdminID, qb)
	if err != nil {
		return nil, fmt.Errorf("查询部门用户失败: %w", err)
	}
	return u.toDomainList(users), nil
}

// CountUsersByDepartment 统计部门用户数量
func (u UserQueryService) CountUsersByDepartment(ctx context.Context, deptID string, excludeAdminID string, qb *dquery.QueryBuilder) (int64, error) {
	if deptID == "" {
		return 0, fmt.Errorf("部门ID不能为空")
	}
	return u.userRepo.CountByDepartment(ctx, deptID, excludeAdminID, qb)
}

// FindUnassignedUsers 查询未分配部门的用户
func (u UserQueryService) FindUnassignedUsers(ctx context.Context, qb *dquery.QueryBuilder) ([]*model.User, error) {
	users, err := u.userRepo.FindUnassignedUsers(ctx, qb)
	if err != nil {
		return nil, fmt.Errorf("查询未分配部门用户失败: %w", err)
	}
	return u.toDomainList(users), nil
}

// CountUnassignedUsers 统计未分配部门的用户数量
func (u UserQueryService) CountUnassignedUsers(ctx context.Context, qb *dquery.QueryBuilder) (int64, error) {
	return u.userRepo.CountUnassignedUsers(ctx, qb)
}

func (u UserQueryService) toDomainList(users []*entity.SysUser) []*model.User {
	if users == nil {
		return make([]*model.User, 0)
	}
	return u.userMapper.ToDomainList(users)
}
