package model

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// Department 部门领域模型
type Department struct {
	ID          string
	ParentID    string
	Name        string
	Code        string
	Sequence    int32
	Leader      string
	Phone       string
	Email       string
	Status      int8
	Description string
	AdminID     string
	TenantID    string
	Children    []*Department
	CreatedAt   int64
	UpdatedAt   int64
}

// NewDepartment 创建部门
func NewDepartment(code string, name string, sequence int32) *Department {
	return &Department{
		Code:     code,
		Name:     name,
		Sequence: sequence,
		Status:   1, // 默认启用
	}
}

// Validate 验证部门信息
func (d *Department) Validate() herrors.Herr {
	if d.Name == "" {
		return errors.DepartmentInvalidField("name", "cannot be empty")
	}
	if d.Code == "" {
		return errors.DepartmentInvalidField("code", "cannot be empty")
	}
	if d.Status != 0 && d.Status != 1 {
		return errors.DepartmentStatusInvalid(d.Status)
	}
	return nil
}

// UpdateBasicInfo 更新基本信息
func (d *Department) UpdateBasicInfo(name string, code string, sequence int32) {
	d.Name = name
	d.Code = code
	d.Sequence = sequence
}

// UpdateContactInfo 更新联系信息
func (d *Department) UpdateContactInfo(leader string, phone string, email string) {
	d.Leader = leader
	d.Phone = phone
	d.Email = email
}

// UpdateStatus 更新状态
func (d *Department) UpdateStatus(status int8) {
	d.Status = status
}

// UpdateParent 更新父部门
func (d *Department) UpdateParent(parentID string) {
	d.ParentID = parentID
}

// SetAdmin 设置管理员
func (d *Department) SetAdmin(adminID string) {
	d.AdminID = adminID
}

// IsEnabled 检查部门是否启用
func (d *Department) IsEnabled() bool {
	return d.Status == 1
}

// HasParent 检查是否有父部门
func (d *Department) HasParent() bool {
	return d.ParentID != ""
}

// AddChild 添加子部门
func (d *Department) AddChild(child *Department) {
	if d.Children == nil {
		d.Children = make([]*Department, 0)
	}
	d.Children = append(d.Children, child)
}

// RemoveChild 移除子部门
func (d *Department) RemoveChild(childID string) {
	if d.Children == nil {
		return
	}
	for i, child := range d.Children {
		if child.ID == childID {
			d.Children = append(d.Children[:i], d.Children[i+1:]...)
			return
		}
	}
}

// HasChildren 检查是否有子部门
func (d *Department) HasChildren() bool {
	return len(d.Children) > 0
}

// IsAdmin 检查用户是否是部门管理员
func (d *Department) IsAdmin(userID string) bool {
	return d.AdminID == userID
}
