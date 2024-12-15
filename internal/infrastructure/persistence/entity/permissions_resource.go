package entity

// PermissionsResource 用于基于角色的访问控制 (RBAC) 的菜单资源管理
type PermissionsResource struct {
	ID            int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"`             // 唯一ID
	PermissionsID int64  `json:"permissions_id" gorm:"index;comment:来源于 Menu.ID"`             // 来源于 Permissions.ID
	Method        string `json:"method" gorm:"size:20;comment:HTTP 方法"`                       // HTTP 方法
	Path          string `json:"path" gorm:"size:255;comment:API 请求路径（例如 /api/v1/users/:id）"` // API 请求路径（例如 /api/v1/users/:id）
}

func (a *PermissionsResource) TableName() string {
	return "sys_permissions_resource"
}
