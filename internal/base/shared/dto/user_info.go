package dto

import "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"

// UserInfoDto 用户信息DTO
type UserInfoDto struct {
	User        *UserDto `json:"user"`
	Permissions []string `json:"permissions"` // 所有权限列表
	Roles       []string `json:"roles"`       // 角色列表
	HomePage    string   `json:"homePage"`    // 首页
}

func ToUserInfoDto(user *model.User, permissions []string, roles []string) *UserInfoDto {
	return &UserInfoDto{
		User:        ToUserDto(user),
		Permissions: permissions,
		Roles:       roles,
	}
}
