package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// CreateTenantCommand 创建租户命令
type CreateTenantCommand struct {
	Code        string            `json:"code" validate:"required" label:"租户编码"`
	Name        string            `json:"name" validate:"required" label:"租户名称"`
	Description string            `json:"description" validate:"omitempty,max=200" label:"描述"`
	IsDefault   int8              `json:"isDefault" validate:"oneof=0 1" label:"是否默认租户"`
	ExpireTime  int64             `json:"expireTime" validate:"omitempty" label:"过期时间"`
	AdminUser   CreateUserCommand `json:"adminUser" validate:"required" label:"管理员信息"`
}

func (c *CreateTenantCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdateTenantCommand 更新租户命令
type UpdateTenantCommand struct {
	ID          string `json:"id" validate:"required" label:"租户ID"`
	Name        string `json:"name" validate:"omitempty" label:"租户名称"`
	Description string `json:"description" validate:"omitempty,max=200" label:"描述"`
	IsDefault   int8   `json:"isDefault" validate:"omitempty,oneof=0 1" label:"是否默认租户"`
	ExpireTime  int64  `json:"expireTime" validate:"omitempty" label:"过期时间"`
}

func (c *UpdateTenantCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// DeleteTenantCommand 删除租户命令
type DeleteTenantCommand struct {
	ID string `json:"id" validate:"required" label:"租户ID"`
}

func (c *DeleteTenantCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// AssignTenantPermissionsCommand 分配租户权限命令
type AssignTenantPermissionsCommand struct {
	TenantID      string  `json:"tenantId" validate:"required" label:"租户ID"`
	PermissionIDs []int64 `json:"permissionIds" validate:"required,dive,gt=0" label:"权限ID列表"`
}

func (c *AssignTenantPermissionsCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
