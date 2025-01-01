package model

import (
	"fmt"
)

// DataScope 数据权限范围
type DataScope int8

const (
	DataScopeAll      DataScope = 1 // 全部数据
	DataScopeDeptTree DataScope = 2 // 部门及以下数据
	DataScopeDept     DataScope = 3 // 本部门数据
	DataScopeCustom   DataScope = 4 // 自定义部门数据
	DataScopeSelf     DataScope = 5 // 仅本人数据
)

// DataPermission 数据权限领域模型
type DataPermission struct {
	ID        string    `json:"id"`
	RoleID    int64     `json:"role_id"`   // 角色ID
	Scope     DataScope `json:"scope"`     // 数据范围
	DeptIDs   []string  `json:"dept_ids"`  // 部门ID列表(自定义数据权限时使用)
	TenantID  string    `json:"tenant_id"` // 租户ID
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

// Validate 验证数据权限
func (d *DataPermission) Validate() error {
	if d.RoleID <= 0 {
		return fmt.Errorf("角色ID不能为空")
	}
	if d.Scope < DataScopeAll || d.Scope > DataScopeSelf {
		return fmt.Errorf("无效的数据范围")
	}
	if d.Scope == DataScopeCustom && len(d.DeptIDs) == 0 {
		return fmt.Errorf("自定义数据权限必须指定部门")
	}
	return nil
}

func NewDataPermission(roleID int64, scope DataScope, deptIDs []string) *DataPermission {
	return &DataPermission{
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
