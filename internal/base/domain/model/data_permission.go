package model

import "github.com/google/uuid"

// DataScope 数据权限范围
type DataScope int8

const (
	DataScopeAll      DataScope = 1 // 全部数据
	DataScopeDept     DataScope = 2 // 本部门数据
	DataScopeDeptTree DataScope = 3 // 本部门及下级数据
	DataScopeSelf     DataScope = 4 // 仅本人数据
	DataScopeCustom   DataScope = 5 // 自定义部门数据
)

// DataPermission 数据权限领域模型
type DataPermission struct {
	ID       string    `json:"id"`
	RoleID   int64     `json:"role_id"` // 修改为int64
	Scope    DataScope `json:"scope"`
	DeptIDs  []string  `json:"dept_ids"`
	TenantID string    `json:"tenant_id"`
}

func NewDataPermission(roleID int64, scope DataScope, deptIDs []string) *DataPermission {
	return &DataPermission{
		ID:      uuid.New().String(),
		RoleID:  roleID,
		Scope:   scope,
		DeptIDs: deptIDs,
	}
}

// IsCustomScope 是否为自定义数据范围
func (p *DataPermission) IsCustomScope() bool {
	return p.Scope == DataScopeCustom
}

// HasDeptPermission 是否有指定部门的数据权限
func (p *DataPermission) HasDeptPermission(deptID string) bool {
	if !p.IsCustomScope() {
		return true
	}
	for _, id := range p.DeptIDs {
		if id == deptID {
			return true
		}
	}
	return false
}
