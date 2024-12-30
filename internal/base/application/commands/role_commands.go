package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// CreateRoleCommand 创建角色命令
type CreateRoleCommand struct {
	Code        string  `json:"code" validate:"required" label:"角色编码"`
	Name        string  `json:"name" validate:"required" label:"角色名称"`
	Localize    string  `json:"localize" validate:"omitempty" label:"多语言标识"`
	Description string  `json:"description" validate:"omitempty,max=200" label:"描述"`
	Sequence    int     `json:"sequence" validate:"gte=0" label:"排序"`
	Type        int8    `json:"type" validate:"required,oneof=1 2" label:"角色类型"`
	PermIDs     []int64 `json:"permIds" validate:"omitempty,dive,gt=0" label:"权限ID列表"`
}

func (c *CreateRoleCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdateRoleCommand 更新角色命令
type UpdateRoleCommand struct {
	ID          int64   `json:"id" validate:"required" label:"角色ID"`
	Name        string  `json:"name" validate:"omitempty" label:"角色名称"`
	Description string  `json:"description" validate:"omitempty,max=200" label:"描述"`
	Sequence    int     `json:"sequence" validate:"omitempty,gte=0" label:"排序"`
	Type        int8    `json:"type" validate:"omitempty,oneof=1 2" label:"角色类型"`
	Localize    string  `json:"localize" validate:"omitempty" label:"多语言标识"`
	Status      *int8   `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
	PermIDs     []int64 `json:"permIds" validate:"omitempty,dive,gt=0" label:"权限ID列表"`
}

func (c *UpdateRoleCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// DeleteRoleCommand 删除角色命令
type DeleteRoleCommand struct {
	ID int64 `json:"id" validate:"required" label:"角色ID"`
}

func (c *DeleteRoleCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
