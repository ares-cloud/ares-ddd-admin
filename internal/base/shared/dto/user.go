package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

// UserDto 用户数据传输对象
type UserDto struct {
	ID        string  `json:"id"`        // ID
	Username  string  `json:"username"`  // 用户名
	Name      string  `json:"name"`      // 姓名
	Phone     string  `json:"phone"`     // 手机号
	Email     string  `json:"email"`     // 邮箱
	Status    int8    `json:"status"`    // 状态
	RoleIds   []int64 `json:"roleIds"`   // 角色ID列表
	CreatedAt int64   `json:"createdAt"` // 创建时间
	UpdatedAt int64   `json:"updatedAt"` // 更新时间
}

// ToUserDto 领域模型转换为DTO
func ToUserDto(u *model.User) *UserDto {
	if u == nil {
		return nil
	}
	roleIds := make([]int64, 0, len(u.Roles))
	for _, role := range u.Roles {
		roleIds = append(roleIds, role.ID)
	}

	return &UserDto{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Phone:     u.Phone,
		Email:     u.Email,
		Status:    u.Status,
		RoleIds:   roleIds,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserDtoList 领域模型列表转换为DTO列表
func ToUserDtoList(users []*model.User) []*UserDto {
	if users == nil {
		return nil
	}

	dtos := make([]*UserDto, 0, len(users))
	for _, u := range users {
		if dto := ToUserDto(u); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
