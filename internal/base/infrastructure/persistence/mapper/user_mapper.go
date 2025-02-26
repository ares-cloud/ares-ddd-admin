package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDomain(e *entity.SysUser, roles []*model.Role) *model.User {
	return &model.User{
		ID:             e.ID,
		Name:           e.Name,
		Username:       e.Username,
		Avatar:         e.Avatar,
		Password:       e.Password,
		Phone:          e.Phone,
		Email:          e.Email,
		Remark:         e.Remark,
		InvitationCode: e.InvitationCode,
		Status:         e.Status,
		Roles:          roles,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		TenantID:       e.TenantID,
	}
}

func (m *UserMapper) ToEntity(d *model.User) *entity.SysUser {
	return &entity.SysUser{
		ID:             d.ID,
		Username:       d.Username,
		Name:           d.Name,
		Avatar:         d.Avatar,
		Password:       d.Password,
		Phone:          d.Phone,
		Email:          d.Email,
		Remark:         d.Remark,
		InvitationCode: d.InvitationCode,
		Status:         d.Status,
	}
}

func (m *UserMapper) ToDomainList(e []*entity.SysUser) []*model.User {
	if len(e) == 0 {
		return nil
	}
	users := make([]*model.User, len(e))
	for i, user := range e {
		users[i] = m.ToDomain(user, make([]*model.Role, 0))
	}
	return users
}

func (m *UserMapper) ToEntityList(d []*model.User) []*entity.SysUser {
	if len(d) == 0 {
		return nil
	}
	users := make([]*entity.SysUser, len(d))
	for i, user := range d {
		users[i] = m.ToEntity(user)
	}
	return users
}
