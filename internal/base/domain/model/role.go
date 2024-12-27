package model

import "time"

// RoleType 角色类型
type RoleType int8

const (
	RoleTypeResource RoleType = 1 // 资源类型角色
	RoleTypeData     RoleType = 2 // 数据权限角色
)

// Role 角色领域模型
type Role struct {
	ID          int64          `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Type        RoleType       `json:"type"` // 角色类型
	Localize    string         `json:"localize"`
	Description string         `json:"description"`
	Sequence    int            `json:"sequence"`
	Status      int8           `json:"status"`
	TenantID    string         `json:"tenant_id"`
	Permissions []*Permissions `json:"permissions"` // 角色拥有的权限列表
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

func NewRole(code, name string, sequence int) *Role {
	return &Role{
		Code:        code,
		Name:        name,
		Type:        RoleTypeResource, // 默认为资源类型角色
		Sequence:    sequence,
		Status:      1,
		Permissions: make([]*Permissions, 0),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
}

// AssignPermissions 分配权限
func (r *Role) AssignPermissions(permissions []*Permissions) {
	r.Permissions = permissions
	r.UpdatedAt = time.Now().Unix()
}

// HasPermission 检查是否有指定权限
func (r *Role) HasPermission(permissionCode string) bool {
	for _, p := range r.Permissions {
		if p.Code == permissionCode {
			return true
		}
	}
	return false
}

// IsResourceRole 是否为资源类型角色
func (r *Role) IsResourceRole() bool {
	return r.Type == RoleTypeResource
}

// IsDataRole 是否为数据权限角色
func (r *Role) IsDataRole() bool {
	return r.Type == RoleTypeData
}

func (r *Role) UpdateBasicInfo(name, description string, sequence int) {
	r.Name = name
	r.Description = description
	r.Sequence = sequence
	r.UpdatedAt = time.Now().Unix()
}

func (r *Role) UpdateLocalize(localize string) {
	if r.Localize != "" {
		r.Localize = localize
	}
}

func (r *Role) UpdateStatus(status int8) {
	r.Status = status
	r.UpdatedAt = time.Now().Unix()
}
