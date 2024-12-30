package model

import (
	"regexp"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

const (
	StatusEnabled  = 1 // 启用
	StatusDisabled = 2 // 禁用/锁定
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
	LockReason  string // 禁用原因
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

// Validate 验证租户模型
func (t *Tenant) Validate() herrors.Herr {
	// 验证租户编码
	if t.Code == "" {
		return errors.TenantInvalidField("code", "cannot be empty")
	}
	if len(t.Code) > 50 {
		return errors.TenantInvalidField("code", "too long, max length is 50")
	}
	// 验证租户编码格式
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(t.Code) {
		return errors.TenantInvalidField("code", "only lowercase letters, numbers and hyphens are allowed")
	}

	// 验证租户名称
	if t.Name == "" {
		return errors.TenantInvalidField("name", "cannot be empty")
	}
	if len(t.Name) > 100 {
		return errors.TenantInvalidField("name", "too long, max length is 100")
	}

	// 验证域名
	if t.Domain != "" {
		if len(t.Domain) > 255 {
			return errors.TenantDomainInvalid(t.Domain, "too long, max length is 255")
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$`).MatchString(t.Domain) {
			return errors.TenantDomainInvalid(t.Domain, "invalid format")
		}
	}

	// 验证管理员用户
	if t.AdminUser == nil {
		return errors.TenantAdminInvalid("admin user is required")
	}
	if err := t.AdminUser.Validate(); err != nil {
		return errors.TenantAdminInvalid(err.Error())
	}

	// 验证过期时间
	if t.ExpireTime > 0 && t.ExpireTime < time.Now().Unix() {
		return errors.TenantExpired()
	}

	// 验证状态
	if t.Status != StatusEnabled && t.Status != StatusDisabled {
		return errors.TenantStatusInvalid(t.Status)
	}

	return nil
}

// UpdateBasicInfo 更新基本信息
func (t *Tenant) UpdateBasicInfo(name, description string) herrors.Herr {
	if name == "" {
		return errors.TenantInvalidField("name", "cannot be empty")
	}
	if len(name) > 100 {
		return errors.TenantInvalidField("name", "too long, max length is 100")
	}

	t.Name = name
	t.Description = description
	return nil
}

// UpdateStatus 更新状态
func (t *Tenant) UpdateStatus(status int8) herrors.Herr {
	if status != 1 && status != 2 && status != 3 {
		return errors.TenantStatusInvalid(status)
	}
	t.Status = status
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateIsDefault 更新是否为默认租户
func (t *Tenant) UpdateIsDefault(isDefault int8) herrors.Herr {
	if isDefault != 1 && isDefault != 2 {
		return errors.TenantInvalidField("is_default", "must be 1(default) or 2(not default)")
	}
	t.IsDefault = isDefault
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateExpireTime 更新过期时间
func (t *Tenant) UpdateExpireTime(expireTime int64) herrors.Herr {
	if expireTime > 0 && expireTime < time.Now().Unix() {
		return errors.TenantExpired()
	}
	t.ExpireTime = expireTime
	return nil
}

// IsActive 检查租户是否有效
func (t *Tenant) IsActive() (bool, herrors.Herr) {
	if t.Status == StatusDisabled {
		return false, errors.TenantDisabled(t.LockReason)
	}
	if time.Now().After(time.Unix(t.ExpireTime, 0)) {
		return false, errors.TenantExpired()
	}
	return true, nil
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

// CheckQuota 检查资源配额
func (t *Tenant) CheckQuota(resource string, current, limit int64) herrors.Herr {
	if current >= limit {
		return errors.TenantQuotaExceeded(resource, limit)
	}
	return nil
}

// Lock 锁定租户
func (t *Tenant) Lock(reason string) herrors.Herr {
	if t.Status == StatusDisabled {
		return errors.TenantInvalidField("status", "tenant is already disabled")
	}
	t.Status = StatusDisabled
	t.LockReason = reason
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// Unlock 解锁租户
func (t *Tenant) Unlock() herrors.Herr {
	if t.Status != StatusDisabled {
		return errors.TenantInvalidField("status", "tenant is not disabled")
	}
	t.Status = StatusEnabled
	t.LockReason = ""
	t.UpdatedAt = time.Now().Unix()
	return nil
}

// IsLocked 检查租户是否被禁用
func (t *Tenant) IsLocked() (bool, string) {
	if t.Status == StatusDisabled {
		return true, t.LockReason
	}
	return false, ""
}
