package model

import (
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type Permissions struct {
	ID          int64
	Code        string
	Name        string
	Localize    string
	Icon        string
	Description string
	Sequence    int
	Type        int8
	Path        string
	Properties  string
	Status      int8
	ParentID    int64
	ParentPath  string
	Resources   []*PermissionsResource
	CreatedAt   int64
	UpdatedAt   int64
	Children    []*Permissions
}

func NewPermissions(code, name string, permType int8, sequence int) *Permissions {
	return &Permissions{
		Code:      code,
		Name:      name,
		Type:      permType,
		Sequence:  sequence,
		Status:    1,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

// Validate 验证权限信息
func (p *Permissions) Validate() herrors.Herr {
	if p.Code == "" {
		return errors.PermissionInvalidField("code", "code cannot be empty")
	}
	if p.Name == "" {
		return errors.PermissionInvalidField("name", "name cannot be empty")
	}
	if p.Type <= 0 || p.Type > 3 {
		return errors.PermissionInvalidField("type", "invalid permission type")
	}
	return nil
}

// UpdateBasicInfo 更新基本信息
func (p *Permissions) UpdateBasicInfo(name, description string, sequence int) herrors.Herr {
	if name == "" {
		return errors.PermissionInvalidField("name", "name cannot be empty")
	}
	p.Name = name
	p.Description = description
	p.Sequence = sequence
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateStatus 更新状态
func (p *Permissions) UpdateStatus(status int8) herrors.Herr {
	if status != 0 && status != 1 {
		return errors.PermissionInvalidField("status", "invalid status value")
	}
	p.Status = status
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// ChangeType 修改权限类型
func (p *Permissions) ChangeType(tp int8) herrors.Herr {
	if tp <= 0 || tp > 3 {
		return errors.PermissionInvalidField("type", "invalid permission type")
	}
	p.Type = tp
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// ChangeParentID 修改父权限
func (p *Permissions) ChangeParentID(pid int64) herrors.Herr {
	if pid < 0 {
		return errors.PermissionInvalidField("parent_id", "invalid parent id")
	}
	p.ParentID = pid
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// AddResource 添加资源
func (p *Permissions) AddResource(method, path string) herrors.Herr {
	if method == "" || path == "" {
		return errors.PermissionInvalidField("resource", "method and path cannot be empty")
	}
	p.Resources = append(p.Resources, &PermissionsResource{
		Method: method,
		Path:   path,
	})
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// UpdateResources 更新资源列表
func (p *Permissions) UpdateResources(resources []*PermissionsResource) herrors.Herr {
	for _, r := range resources {
		if r.Method == "" || r.Path == "" {
			return errors.PermissionInvalidField("resource", "method and path cannot be empty")
		}
	}
	p.Resources = resources
	p.UpdatedAt = time.Now().Unix()
	return nil
}

// IsEnabled 是否启用
func (p *Permissions) IsEnabled() bool {
	return p.Status == 1
}

// HasChildren 是否有子权限
func (p *Permissions) HasChildren() bool {
	return len(p.Children) > 0
}
