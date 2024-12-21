package model

import "time"

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

func (p *Permissions) UpdateBasicInfo(name, description string, sequence int) {
	p.Name = name
	p.Description = description
	p.Sequence = sequence
	p.UpdatedAt = time.Now().Unix()
}

func (p *Permissions) UpdateStatus(status int8) {
	p.Status = status
	p.UpdatedAt = time.Now().Unix()
}
func (p *Permissions) ChangeType(tp int8) {
	if tp > 0 && tp < 4 {
		p.Type = tp
	}
}
func (p *Permissions) ChangeParentID(pid int64) {
	p.ParentID = pid
}
func (p *Permissions) AddResource(method, path string) {
	p.Resources = append(p.Resources, &PermissionsResource{
		Method: method,
		Path:   path,
	})
}

// UpdateResources 更新资源列表
func (p *Permissions) UpdateResources(resources []*PermissionsResource) {
	p.Resources = resources
	p.UpdatedAt = time.Now().Unix()
}
