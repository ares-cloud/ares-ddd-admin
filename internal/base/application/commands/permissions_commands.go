package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// CreatePermissionsCommand 创建权限命令
type CreatePermissionsCommand struct {
	Code        string                             `json:"code" validate:"required" label:"权限编码"`
	Name        string                             `json:"name" validate:"required" label:"权限名称"`
	Localize    string                             `json:"localize" validate:"omitempty" label:"多语言标识"`
	Icon        string                             `json:"icon" validate:"omitempty" label:"图标"`
	Description string                             `json:"description" validate:"omitempty,max=200" label:"描述"`
	Sequence    int                                `json:"sequence" validate:"gte=0" label:"排序"`
	Type        int8                               `json:"type" validate:"required,oneof=1 2 3" label:"权限类型"`
	Path        string                             `json:"path" validate:"omitempty" label:"路径"`
	Properties  string                             `json:"properties" validate:"omitempty" label:"扩展属性"`
	ParentID    int64                              `json:"parentId" validate:"omitempty" label:"父权限ID"`
	Resources   []CreatePermissionsResourceCommand `json:"resources" validate:"omitempty,dive" label:"资源列表"`
}

func (c *CreatePermissionsCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// CreatePermissionsResourceCommand 创建权限资源命令
type CreatePermissionsResourceCommand struct {
	Method string `json:"method" validate:"required,oneof=GET POST PUT DELETE" label:"请求方法"`
	Path   string `json:"path" validate:"required" label:"资源路径"`
}

func (c *CreatePermissionsResourceCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdatePermissionsCommand 更新权限命令
type UpdatePermissionsCommand struct {
	ID          int64                              `json:"id" validate:"required" label:"权限ID"`
	Name        string                             `json:"name" validate:"omitempty" label:"权限名称"`
	Icon        string                             `json:"icon" validate:"omitempty" label:"图标"`
	Description string                             `json:"description" validate:"omitempty,max=200" label:"描述"`
	Sequence    int                                `json:"sequence" validate:"omitempty,gte=0" label:"排序"`
	Path        string                             `json:"path" validate:"omitempty" label:"路径"`
	Properties  string                             `json:"properties" validate:"omitempty" label:"扩展属性"`
	Status      *int8                              `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
	ParentID    int64                              `json:"parentId" validate:"omitempty" label:"父权限ID"`
	Localize    string                             `json:"localize" validate:"omitempty" label:"多语言标识"`
	Type        int8                               `json:"type" validate:"omitempty,oneof=1 2 3" label:"权限类型"`
	Resources   []CreatePermissionsResourceCommand `json:"resources" validate:"omitempty,dive" label:"资源列表"`
}

func (c *UpdatePermissionsCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// DeletePermissionsCommand 删除权限命令
type DeletePermissionsCommand struct {
	ID int64 `json:"id" validate:"required" label:"权限ID"`
}

func (c *DeletePermissionsCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
