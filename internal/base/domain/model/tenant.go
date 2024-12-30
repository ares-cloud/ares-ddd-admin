package model

import (
	"errors"
	"time"
)

// Tenant 租户领域模型
type Tenant struct {
	ID          string
	Code        string // 租户编码(唯一)
	Name        string // 租户名称
	Domain      string // 域名
	AdminUser   *User  // 管理员用户
	Status      int8   // 状态(1:启用 2:禁用)
	IsDefault   int8   // 是否默认租户(1:是 2:否)
	ExpireTime  int64  // 过期时间
	Description string // 描述
	CreatedAt   int64
	UpdatedAt   int64
	Permissions []*Permissions // 租户拥有的权限
}

// NewTenant 创建新租户
func NewTenant(code, name string, adminUser *User) *Tenant {
	now := time.Now()
	return &Tenant{
		Code:       code,
		Name:       name,
		AdminUser:  adminUser,
		Status:     1,
		IsDefault:  2,                           // 默认为非默认租户
		ExpireTime: now.AddDate(1, 0, 0).Unix(), // 默认一年有效期
		CreatedAt:  now.Unix(),
		UpdatedAt:  now.Unix(),
	}
}

// UpdateBasicInfo 更新基本信息
func (t *Tenant) UpdateBasicInfo(name, description string) {
	t.Name = name
	t.Description = description
	t.UpdatedAt = time.Now().Unix()
}

// UpdateStatus 更新状态
func (t *Tenant) UpdateStatus(status int8) error {
	if status != 1 && status != 2 {
		return errors.New("invalid status: must be 1(enabled) or 2(disabled)")
	}
	t.Status = status
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateIsDefault 更新是否为默认租户
func (t *Tenant) UpdateIsDefault(isDefault int8) error {
	if isDefault != 1 && isDefault != 2 {
		return errors.New("invalid isDefault value: must be 1(default) or 2(not default)")
	}
	t.IsDefault = isDefault
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateExpireTime 更新过期时间
func (t *Tenant) UpdateExpireTime(expireTime int64) {
	t.ExpireTime = expireTime
	t.UpdatedAt = time.Now().Unix()
}

// IsActive 检查租户是否有效
func (t *Tenant) IsActive() bool {
	return t.Status == 1 && time.Now().Before(time.Unix(t.ExpireTime, 0))
}

// IsDefaultTenant 是否为默认租户
func (t *Tenant) IsDefaultTenant() bool {
	return t.IsDefault == 1
}

// AssignPermissions 分配权限给租户
func (t *Tenant) AssignPermissions(permissions []*Permissions) {
	t.Permissions = permissions
	t.UpdatedAt = time.Now().Unix()
}

// HasPermission 检查租户是否拥有指定权限
func (t *Tenant) HasPermission(permissionID int64) bool {
	for _, perm := range t.Permissions {
		if perm.ID == permissionID {
			return true
		}
	}
	return false
}

// GetPermissionIDs 获取权限ID列表
func (t *Tenant) GetPermissionIDs() []int64 {
	ids := make([]int64, 0, len(t.Permissions))
	for _, perm := range t.Permissions {
		ids = append(ids, perm.ID)
	}
	return ids
}
