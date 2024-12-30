package model

import (
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/password"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

const (
	UserStatusEnabled  = 1 // 启用
	UserStatusDisabled = 2 // 禁用

	// 密码相关常量
	MinPasswordLength = 6  // 最小密码长度
	MaxPasswordLength = 32 // 最大密码长度
)

// User 用户领域模型
type User struct {
	ID             string  `json:"id"`              // 用户ID
	TenantID       string  `json:"tenant_id"`       // 租户ID
	Username       string  `json:"username"`        // 用户名
	Password       string  `json:"-"`               // 密码(不序列化)
	Name           string  `json:"name"`            // 姓名
	Nickname       string  `json:"nickname"`        // 昵称
	Avatar         string  `json:"avatar"`          // 头像
	Email          string  `json:"email"`           // 邮箱
	Phone          string  `json:"phone"`           // 手机号
	Remark         string  `json:"remark"`          // 备注
	InvitationCode string  `json:"invitation_code"` // 邀请码
	Status         int8    `json:"status"`          // 状态
	Roles          []*Role `json:"roles"`           // 角色列表
	CreatedAt      int64   `json:"created_at"`      // 创建时间
	UpdatedAt      int64   `json:"updated_at"`      // 更新时间
}

// NewUser 创建新用户
func NewUser(tenantID, username, password string) *User {
	now := time.Now().Unix()
	return &User{
		TenantID:  tenantID,
		Username:  username,
		Password:  password,
		Name:      username,
		Nickname:  username,
		Status:    UserStatusEnabled,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate 验证用户模型
func (u *User) Validate() herrors.Herr {
	// 验证租户ID
	//if !validator.ValidateRequired(u.TenantID) {
	//	return errors.UserInvalidField("tenant_id", "cannot be empty")
	//}

	// 验证用户名
	if !validator.ValidateRequired(u.Username) {
		return errors.UserInvalidField("username", "cannot be empty")
	}
	if !validator.ValidateLength(u.Username, 0, 50) {
		return errors.UserInvalidField("username", "too long, max length is 50")
	}
	if !validator.ValidateUsername(u.Username) {
		return errors.UserInvalidField("username", "only letters, numbers, underscore and hyphen are allowed")
	}

	// 验证密码
	if !validator.ValidateRequired(u.Password) {
		return errors.UserInvalidField("password", "cannot be empty")
	}
	if !validator.ValidatePassword(u.Password) {
		return errors.UserInvalidField("password", "too short, min length is 6")
	}

	// 验证邮箱
	if u.Email != "" {
		if !validator.ValidateLength(u.Email, 0, 100) {
			return errors.UserInvalidField("email", "too long, max length is 100")
		}
		if !validator.ValidateEmail(u.Email) {
			return errors.UserInvalidField("email", "invalid email format")
		}
	}

	// 验证手机号
	if u.Phone != "" {
		if !validator.ValidatePhone(u.Phone) {
			return errors.UserInvalidField("phone", "invalid phone number format")
		}
	}

	// 验证状态
	if u.Status != UserStatusEnabled && u.Status != UserStatusDisabled {
		return errors.UserStatusInvalid(u.Status)
	}

	return nil
}

// HashPassword 加密密码
func (u *User) HashPassword() herrors.Herr {
	// 1. 验证密码长度
	if len(u.Password) < MinPasswordLength {
		return errors.UserInvalidField("password", "password too short")
	}
	if len(u.Password) > MaxPasswordLength {
		return errors.UserInvalidField("password", "password too long")
	}

	// 2. 使用 bcrypt 加密
	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword 比较密码
func (u *User) ComparePassword(pas string) herrors.Herr {
	// 1. 验证密码长度
	if len(pas) < MinPasswordLength {
		return errors.UserInvalidField("password", "password too short")
	}
	if len(pas) > MaxPasswordLength {
		return errors.UserInvalidField("password", "password too long")
	}

	// 2. 比较密码
	if ok := password.CheckPasswordHash(pas, u.Password); !ok {
		return errors.UserInvalidField("password", "password mismatch")
	}

	return nil
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() (bool, string) {
	if u.Status == UserStatusDisabled {
		return true, "user is disabled"
	}
	return false, ""
}

// Lock 锁定用户
func (u *User) Lock(reason string) herrors.Herr {
	if u.Status == UserStatusDisabled {
		return errors.UserInvalidField("status", "user is already disabled")
	}
	u.Status = UserStatusDisabled
	u.UpdatedAt = time.Now().Unix()
	return nil
}

// Unlock 解锁用户
func (u *User) Unlock() herrors.Herr {
	if u.Status != UserStatusDisabled {
		return errors.UserInvalidField("status", "user is not disabled")
	}
	u.Status = UserStatusEnabled
	u.UpdatedAt = time.Now().Unix()
	return nil
}

// AssignRoles 分配角色
func (u *User) AssignRoles(roles []*Role) {
	u.Roles = roles
	u.UpdatedAt = time.Now().Unix()
}

// HasRole 检查是否拥有指定角色
func (u *User) HasRole(roleID int64) bool {
	for _, role := range u.Roles {
		if role.ID == roleID {
			return true
		}
	}
	return false
}

// UpdateBasicInfo 更新基本信息
func (u *User) UpdateBasicInfo(name, nickname, phone, email, avatar, remark string) {
	u.Name = name
	u.Nickname = nickname
	u.Phone = phone
	u.Email = email
	u.Avatar = avatar
	u.Remark = remark
	u.UpdatedAt = time.Now().Unix()
}

// UpdateStatus 更新状态
func (u *User) UpdateStatus(status int8) herrors.Herr {
	if status != UserStatusEnabled && status != UserStatusDisabled {
		return errors.UserStatusInvalid(status)
	}
	u.Status = status
	u.UpdatedAt = time.Now().Unix()
	return nil
}
