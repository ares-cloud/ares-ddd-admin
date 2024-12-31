package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

// UserConverter 用户转换器
type UserConverter struct{}

func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToDTO 将领域模型转换为DTO
func (c *UserConverter) ToDTO(user *entity.SysUser, roleIds []int64) *dto.UserDto {
	if user == nil {
		return nil
	}
	return &dto.UserDto{
		ID:             user.ID,
		TenantID:       user.TenantID,
		Username:       user.Username,
		Avatar:         user.Avatar,
		Name:           user.Name,
		Nickname:       user.Nickname,
		Phone:          user.Phone,
		Email:          user.Email,
		Remark:         user.Remark,
		InvitationCode: user.InvitationCode,
		Status:         user.Status,
		RoleIds:        roleIds,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// ToDTOList 将领域模型列表转换为DTO列表
func (c *UserConverter) ToDTOList(users []*entity.SysUser) []*dto.UserDto {
	dos := make([]*dto.UserDto, 0, len(users))
	for _, user := range users {
		if userDto := c.ToDTO(user, nil); userDto != nil {
			dos = append(dos, userDto)
		}
	}
	return dos
}
