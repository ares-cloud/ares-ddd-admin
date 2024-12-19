package casbin

type ApiPermissions struct {
	Id     string `json:"id"`
	Method string `json:"method" ` // HTTP 方法
	Path   string `json:"path"`    // API 请求路径（例如 /api/v1/users/:id）
}

type Role struct {
	Id          string           `json:"id"`
	Code        string           `json:"code"`        // 角色代码（唯一）
	TenantID    string           `json:"tenant_id"`   // 租户ID
	Permissions []ApiPermissions `json:"permissions"` //权限
}
