package query

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

// IRoleQueryService 角色查询服务接口
type IRoleQueryService interface {
	// GetRole 获取角色详情
	GetRole(ctx context.Context, id int64) (*dto.RoleDto, error)

	// GetRoleByCode 根据编码获取角色
	GetRoleByCode(ctx context.Context, code string) (*dto.RoleDto, error)

	// FindRoles 查询角色列表
	FindRoles(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.RoleDto, error)

	// CountRoles 统计角色数量
	CountRoles(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)

	// GetRolePermissions 获取角色权限
	GetRolePermissions(ctx context.Context, roleID int64) ([]*dto.PermissionsDto, error)
}
