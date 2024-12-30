package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// CreateUserCommand 创建用户命令
type CreateUserCommand struct { // 租户ID
	Username string  `json:"username" validate:"required" label:"用户名"`          // 用户名
	Nickname string  `json:"nickname" label:"昵称"`                               // 昵称
	Name     string  `json:"name"  label:"名称"`                                  // 名称
	Password string  `json:"password" validate:"required,min=6" label:"密码"`     // 密码
	Phone    string  `json:"phone" validate:"omitempty,mobile" label:"手机号"`     // 手机号
	Email    string  `json:"email" validate:"omitempty,email" label:"邮箱"`       // 邮箱
	Avatar   string  `json:"avatar" validate:"omitempty,url" label:"头像"`        // 头像
	RoleIDs  []int64 `json:"roleIds" validate:"omitempty,dive,gt=0" label:"角色"` // 角色ID列表
}

func (c *CreateUserCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
	ID       string  `json:"id" validate:"required" label:"用户ID"`
	Nickname string  `json:"nickname" validate:"omitempty" label:"昵称"`
	Phone    string  `json:"phone" validate:"omitempty,mobile" label:"手机号"`
	Email    string  `json:"email" validate:"omitempty,email" label:"邮箱"`
	Avatar   string  `json:"avatar" validate:"omitempty,url" label:"头像"`
	Name     string  `json:"name"  label:"名称"` // 名称
	Status   int8    `json:"status" validate:"omitempty,oneof=1 2" label:"状态"`
	RoleIDs  []int64 `json:"roleIds" validate:"omitempty,dive,gt=0" label:"角色"`
}

func (c *UpdateUserCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// DeleteUserCommand 删除用户命令
type DeleteUserCommand struct {
	ID string `json:"id" validate:"required" label:"用户ID"`
}

func (c *DeleteUserCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdateUserStatusCommand 更新用户状态命令
type UpdateUserStatusCommand struct {
	ID     string `json:"id" validate:"required" label:"用户ID"`
	Status int8   `json:"status" validate:"required,oneof=1 2" label:"状态"`
}

func (c *UpdateUserStatusCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
