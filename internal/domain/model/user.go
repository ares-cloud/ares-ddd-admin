package model

import (
	"errors"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/password"
)

type User struct {
	ID             string
	TenantID       string // 租户ID
	Username       string
	FaceURL        string
	Name           string
	Password       string
	Phone          string
	Email          string
	Remark         string
	InvitationCode string
	Status         int8
	Roles          []*Role
	CreatedAt      int64
	UpdatedAt      int64
}

// NewUser 创建新用户
func NewUser(username, name, password string) *User {
	now := time.Now().Unix()
	return &User{
		Username:  username,
		Name:      name,
		Password:  password,
		Status:    1, // 默认启用
		Roles:     make([]*Role, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HashPassword 密码加密
func (u *User) HashPassword() error {
	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(pas string) bool {
	return password.CheckPasswordHash(pas, u.Password)
}

// UpdateBasicInfo 更新基本信息
func (u *User) UpdateBasicInfo(name, phone, email, faceURL, remark string) {
	u.Name = name
	u.Phone = phone
	u.Email = email
	u.FaceURL = faceURL
	u.Remark = remark
	u.UpdatedAt = time.Now().Unix()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.CheckPassword(oldPassword) {
		return errors.New("old password is incorrect")
	}
	u.Password = newPassword
	return u.HashPassword()
}

// UpdateStatus 更新状态
func (u *User) UpdateStatus(status int8) error {
	if status != 1 && status != 2 {
		return errors.New("invalid status")
	}
	u.Status = status
	u.UpdatedAt = time.Now().Unix()
	return nil
}

// AssignRoles 分配角色
func (u *User) AssignRoles(roles []*Role) {
	u.Roles = roles
	u.UpdatedAt = time.Now().Unix()
}

// AddRole 添加角色
func (u *User) AddRole(role *Role) {
	// 检查角色是否已存在
	for _, r := range u.Roles {
		if r.ID == role.ID {
			return
		}
	}
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now().Unix()
}

// RemoveRole 移除角色
func (u *User) RemoveRole(roleID int64) {
	for i, role := range u.Roles {
		if role.ID == roleID {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			u.UpdatedAt = time.Now().Unix()
			return
		}
	}
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

// IsActive 检查用户是否启用
func (u *User) IsActive() bool {
	return u.Status == 1
}

// GetRoleIDs 获取角色ID列表
func (u *User) GetRoleIDs() []int64 {
	ids := make([]int64, 0, len(u.Roles))
	for _, role := range u.Roles {
		ids = append(ids, role.ID)
	}
	return ids
}

// SetInvitationCode 设置邀请码
func (u *User) SetInvitationCode(code string) {
	u.InvitationCode = code
	u.UpdatedAt = time.Now().Unix()
}
