package dto

// RoleDto 角色DTO
type RoleDto struct {
	ID          int64   `json:"id"`          // 角色ID
	Code        string  `json:"code"`        // 角色代码
	Name        string  `json:"name"`        // 角色名称
	Type        int8    `json:"type"`        // 角色类型(1:资源角色 2:数据权限角色)
	Localize    string  `json:"localize"`    // 国际化key
	Description string  `json:"description"` // 描述
	Sequence    int     `json:"sequence"`    // 排序
	Status      int8    `json:"status"`      // 状态
	PermIds     []int64 `json:"permIds"`     // 权限id
	TenantID    string  `json:"tenantId"`    // 租户ID
	CreatedAt   int64   `json:"createdAt"`   // 创建时间
	UpdatedAt   int64   `json:"updatedAt"`   // 更新时间
}
