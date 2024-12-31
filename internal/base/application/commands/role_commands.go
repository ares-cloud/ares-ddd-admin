package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// AssignRolePermissionsCommand 分配角色权限命令
type AssignRolePermissionsCommand struct {
	RoleID        int64   `json:"role_id" binding:"required"`        // 角色ID
	PermissionIDs []int64 `json:"permission_ids" binding:"required"` // 权限ID列表
}

// Validate 验证命令
func (c *AssignRolePermissionsCommand) Validate() herrors.Herr {
	// 验证角色ID
	if c.RoleID <= 0 {
		return errors.RoleInvalidField("id", "must be greater than 0")
	}

	// 验证权限ID列表
	if len(c.PermissionIDs) == 0 {
		return errors.RoleInvalidField("permission_ids", "cannot be empty")
	}

	// 验证权限ID有效性
	for _, permID := range c.PermissionIDs {
		if permID <= 0 {
			return herrors.NewBadReqError("permission id must be greater than 0")
		}
	}

	return nil
}

// CreateRoleCommand 创建角色命令
type CreateRoleCommand struct {
	Code        string `json:"code" binding:"required"` // 角色编码
	Name        string `json:"name" binding:"required"` // 角色名称
	Type        int8   `json:"type" binding:"required"` // 角色类型
	Localize    string `json:"localize"`                // 多语言标识
	Description string `json:"description"`             // 描述
	Sequence    int    `json:"sequence"`                // 排序
}

// Validate 验证命令
func (c *CreateRoleCommand) Validate() herrors.Herr {
	// 验证编码
	if c.Code == "" {
		return errors.RoleInvalidField("code", "cannot be empty")
	}
	if len(c.Code) > 50 {
		return errors.RoleInvalidField("code", "too long")
	}

	// 验证名称
	if c.Name == "" {
		return errors.RoleInvalidField("name", "cannot be empty")
	}
	if len(c.Name) > 50 {
		return errors.RoleInvalidField("name", "too long")
	}

	// 验证类型
	if c.Type != 1 && c.Type != 2 {
		return errors.RoleInvalidField("type", "must be 1(system) or 2(custom)")
	}

	return nil
}

// UpdateRoleCommand 更新角色命令
type UpdateRoleCommand struct {
	ID          int64  `json:"id" binding:"required"`   // 角色ID
	Name        string `json:"name" binding:"required"` // 角色名称
	Localize    string `json:"localize"`                // 多语言标识
	Description string `json:"description"`             // 描述
	Sequence    int    `json:"sequence"`                // 排序
	Status      int8   `json:"status"`                  // 状态
}

// Validate 验证命令
func (c *UpdateRoleCommand) Validate() herrors.Herr {
	// 验证ID
	if c.ID <= 0 {
		return errors.RoleInvalidField("id", "must be greater than 0")
	}

	// 验证名称
	if c.Name == "" {
		return errors.RoleInvalidField("name", "cannot be empty")
	}
	if len(c.Name) > 50 {
		return errors.RoleInvalidField("name", "too long")
	}

	// 验证状态
	if c.Status != 0 && c.Status != 1 && c.Status != 2 {
		return errors.RoleStatusInvalid(c.Status)
	}

	return nil
}

// DeleteRoleCommand 删除角色命令
type DeleteRoleCommand struct {
	ID int64 `json:"id" binding:"required"` // 角色ID
}

// Validate 验证命令
func (c *DeleteRoleCommand) Validate() herrors.Herr {
	if c.ID <= 0 {
		return errors.RoleInvalidField("id", "must be greater than 0")
	}
	return nil
}
