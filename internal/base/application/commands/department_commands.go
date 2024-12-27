package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// CreateDepartmentCommand 创建部门命令
type CreateDepartmentCommand struct {
	ParentID    string `json:"parentId"`                                            // 父部门ID
	Name        string `json:"name" validate:"required" label:"部门名称"`               // 必须是中文
	Code        string `json:"code" validate:"required,min=2,max=50" label:"部门编码"`  // 部门编码
	Sort        int    `json:"sort" validate:"gte=0,lte=999" label:"显示顺序"`          // 排序
	Leader      string `json:"leader" validate:"omitempty" label:"负责人"`             // 必须是中文
	Phone       string `json:"phone" validate:"omitempty,mobile" label:"联系电话"`      // 必须是手机号
	Email       string `json:"email" validate:"omitempty,email,max=100" label:"邮箱"` // 邮箱
	Status      int8   `json:"status" validate:"oneof=0 1" label:"部门状态"`            // 部门状态(0停用 1启用)
	Description string `json:"description" validate:"omitempty,max=200" label:"描述"` // 描述
	TenantID    string `json:"-"`                                                   // 租户ID
}

func (c *CreateDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// UpdateDepartmentCommand 更新部门命令
type UpdateDepartmentCommand struct {
	ID          string `json:"id" validate:"required" label:"部门ID"`
	ParentID    string `json:"parentId"`
	Name        string `json:"name" validate:"required" label:"部门名称"`
	Code        string `json:"code" validate:"required,min=2,max=50" label:"部门编码"`
	Sort        int    `json:"sort" validate:"gte=0,lte=999" label:"显示顺序"`
	Leader      string `json:"leader" validate:"omitempty,max=50" label:"负责人"`
	Phone       string `json:"phone" validate:"omitempty,e164" label:"联系电话"`
	Email       string `json:"email" validate:"omitempty,email,max=100" label:"邮箱"`
	Status      int8   `json:"status" validate:"oneof=0 1" label:"部门状态"`
	Description string `json:"description" validate:"omitempty,max=200" label:"描述"`
}

func (c *UpdateDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// DeleteDepartmentCommand 删除部门命令
type DeleteDepartmentCommand struct {
	ID string `json:"id" validate:"required" label:"部门ID"`
}

func (c *DeleteDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// MoveDepartmentCommand 移动部门命令
type MoveDepartmentCommand struct {
	ID           string `json:"id" validate:"required" label:"部门ID"`
	TargetParent string `json:"targetParent" validate:"required" label:"目标父部门ID"`
}

func (c *MoveDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// SetDepartmentAdminCommand 设置部门管理员命令
type SetDepartmentAdminCommand struct {
	DeptID  string `json:"deptId" validate:"required" label:"部门ID"`
	AdminID string `json:"adminId" validate:"required" label:"管理员ID"`
}

func (c *SetDepartmentAdminCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// AssignUsersToDepartmentCommand 分配用户到部门命令
type AssignUsersToDepartmentCommand struct {
	DeptID  string   `json:"deptId" validate:"required" label:"部门ID"`
	UserIDs []string `json:"userIds" validate:"required,min=1" label:"用户ID列表"`
}

func (c *AssignUsersToDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// RemoveUsersFromDepartmentCommand 从部门移除用户命令
type RemoveUsersFromDepartmentCommand struct {
	DeptID  string   `json:"deptId" validate:"required" label:"部门ID"`
	UserIDs []string `json:"userIds" validate:"required,min=1" label:"用户ID列表"`
}

func (c *RemoveUsersFromDepartmentCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
