package model

import (
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

const (
	RoleStatusEnabled  = 1 // 启用
	RoleStatusDisabled = 2 // 禁用

	RoleTypeSystem = 1 // 系统角色
	RoleTypeCustom = 2 // 自定义角色
)

// Role 角色领域模型
type Role struct {
	ID          int64          `json:"id"`          // 角色ID
	TenantID    string         `json:"tenant_id"`   // 租户ID
	Code        string         `json:"code"`        // 角色编码
	Name        string         `json:"name"`        // 角色名称
	Localize    string         `json:"localize"`    // 多语言标识
	Description string         `json:"description"` // 描述
	Sequence    int            `json:"sequence"`    // 排序
	Type        int8           `json:"type"`        // 类型
	Status      int8           `json:"status"`      // 状态
	Permissions []*Permissions `json:"permissions"` // 权限列表
	CreatedAt   int64          `json:"created_at"`  // 创建时间
	UpdatedAt   int64          `json:"updated_at"`  // 更新时间
}

// NewRole 创建角色
func NewRole(tenantID, code, name string) *Role {
	now := time.Now().Unix()
	return &Role{
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		Status:    RoleStatusEnabled,
		Type:      RoleTypeCustom,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate 验证角色
func (r *Role) Validate() herrors.Herr {
	// 验证编码
	if !validator.ValidateRequired(r.Code) {
		return errors.RoleInvalidField("code", "cannot be empty")
	}
	if !validator.ValidateLength(r.Code, 0, 50) {
		return errors.RoleInvalidField("code", "too long")
	}

	// 验证名称
	if !validator.ValidateRequired(r.Name) {
		return errors.RoleInvalidField("name", "cannot be empty")
	}
	if !validator.ValidateLength(r.Name, 0, 50) {
		return errors.RoleInvalidField("name", "too long")
	}

	// 验证状态
	if r.Status != RoleStatusEnabled && r.Status != RoleStatusDisabled {
		return errors.RoleStatusInvalid(r.Status)
	}

	return nil
}

// UpdateBasicInfo 更新基本信息
func (r *Role) UpdateBasicInfo(name, localize, description string, sequence int) {
	r.Name = name
	r.Localize = localize
	r.Description = description
	r.Sequence = sequence
	r.UpdatedAt = time.Now().Unix()
}

// UpdateStatus 更新状态
func (r *Role) UpdateStatus(status int8) herrors.Herr {
	if status != RoleStatusEnabled && status != RoleStatusDisabled {
		return errors.RoleStatusInvalid(status)
	}
	r.Status = status
	r.UpdatedAt = time.Now().Unix()
	return nil
}

// AssignPermissions 分配权限
func (r *Role) AssignPermissions(permissions []*Permissions) {
	r.Permissions = permissions
	r.UpdatedAt = time.Now().Unix()
}

// HasPermission 检查是否拥有权限
func (r *Role) HasPermission(permissionID int64) bool {
	for _, p := range r.Permissions {
		if p.ID == permissionID {
			return true
		}
	}
	return false
}
