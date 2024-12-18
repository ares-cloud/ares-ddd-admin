package model

import "time"

type Role struct {
	ID          int64
	TenantID    string
	Code        string
	Name        string
	Localize    string
	Description string
	Sequence    int
	Status      int8
	Permissions []*Permissions
	CreatedAt   int64
	UpdatedAt   int64
}

func NewRole(code, name string, sequence int) *Role {
	return &Role{
		Code:      code,
		Name:      name,
		Sequence:  sequence,
		Status:    1,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (r *Role) AssignPermissions(permissions []*Permissions) {
	r.Permissions = permissions
}

func (r *Role) UpdateBasicInfo(name, description string, sequence int) {
	r.Name = name
	r.Description = description
	r.Sequence = sequence
	r.UpdatedAt = time.Now().Unix()
}

func (r *Role) UpdateLocalize(localize string) {
	if r.Localize != "" {
		r.Localize = localize
	}
}

func (r *Role) UpdateStatus(status int8) {
	r.Status = status
	r.UpdatedAt = time.Now().Unix()
}
