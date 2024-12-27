package dto

import "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"

// DataPermissionDto 数据权限DTO
type DataPermissionDto struct {
	ID       string   `json:"id"`       // ID
	RoleID   int64    `json:"roleId"`   // 角色ID
	Scope    int8     `json:"scope"`    // 数据范围(1:全部数据 2:本部门数据 3:本部门及下级数据 4:仅本人数据 5:自定义部门数据)
	DeptIDs  []string `json:"deptIds"`  // 自定义部门ID列表
	TenantID string   `json:"tenantId"` // 租户ID
}

// ToDataPermissionDto 领域模型转DTO
func ToDataPermissionDto(m *model.DataPermission) *DataPermissionDto {
	if m == nil {
		return nil
	}
	return &DataPermissionDto{
		ID:       m.ID,
		RoleID:   m.RoleID,
		Scope:    int8(m.Scope),
		DeptIDs:  m.DeptIDs,
		TenantID: m.TenantID,
	}
}
