package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// AssignDataPermissionCommand 分配数据权限命令
type AssignDataPermissionCommand struct {
	RoleID  int64    `json:"roleId" validate:"required" label:"角色ID"`                 // 修改为int64
	Scope   int8     `json:"scope" validate:"omitempty,oneof=1 2 3 4 5" label:"数据范围"` // 数据范围
	DeptIDs []string `json:"deptIds" validate:"required_if=Scope 5" label:"部门ID列表"`   // 部门ID列表(自定义数据权限时使用)
}

func (c *AssignDataPermissionCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// RemoveDataPermissionCommand 移除数据权限命令
type RemoveDataPermissionCommand struct {
	RoleID int64 `json:"roleId"` // 修改为int64
}
