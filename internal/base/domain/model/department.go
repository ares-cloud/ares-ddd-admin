package model

import "time"

// Department 部门实体
type Department struct {
	ID          string        `json:"id"`
	TenantID    string        `json:"tenant_id"`
	ParentID    string        `json:"parent_id"`
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Sort        int           `json:"sort"`
	AdminID     string        `json:"admin_id"`
	Leader      string        `json:"leader"`
	Phone       string        `json:"phone"`
	Email       string        `json:"email"`
	Status      int8          `json:"status"`
	Description string        `json:"description"`
	Children    []*Department `json:"children"`
	CreatedAt   int64         `json:"created_at"`
	UpdatedAt   int64         `json:"updated_at"`
}

// NewDepartment 创建部门
func NewDepartment(code, name string, sort int) *Department {
	return &Department{
		Code:      code,
		Name:      name,
		Sort:      sort,
		Status:    1,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

// UpdateBasicInfo 更新基本信息
func (d *Department) UpdateBasicInfo(name, code string, sort int) {
	d.Name = name
	d.Code = code
	d.Sort = sort
	d.UpdatedAt = time.Now().Unix()
}

// UpdateContactInfo 更新联系信息
func (d *Department) UpdateContactInfo(leader, phone, email string) {
	d.Leader = leader
	d.Phone = phone
	d.Email = email
	d.UpdatedAt = time.Now().Unix()
}

// UpdateStatus 更新状态
func (d *Department) UpdateStatus(status int8) {
	d.Status = status
	d.UpdatedAt = time.Now().Unix()
}

// UpdateParent 更新父部门
func (d *Department) UpdateParent(parentID string) {
	d.ParentID = parentID
	d.UpdatedAt = time.Now().Unix()
}

// AddChild 添加子部门
func (d *Department) AddChild(child *Department) {
	child.ParentID = d.ID
	if d.Children == nil {
		d.Children = make([]*Department, 0)
	}
	d.Children = append(d.Children, child)
}

// IsEnabled 是否启用
func (d *Department) IsEnabled() bool {
	return d.Status == 1
}

// SetAdmin 设置部门管理员
func (d *Department) SetAdmin(adminID string, leaderName string, phone string) {
	d.AdminID = adminID
	d.Leader = leaderName
	d.Phone = phone
	d.UpdatedAt = time.Now().Unix()
}
