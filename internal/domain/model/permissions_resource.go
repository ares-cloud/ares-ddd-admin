package model

// PermissionsResource 权限资源模型
type PermissionsResource struct {
	ID            int64  // 唯一标识
	PermissionsID int64  // 关联的权限ID
	Method        string // HTTP方法
	Path          string // API路径
}

// NewPermissionsResource 创建新的权限资源
func NewPermissionsResource(permissionsID int64, method, path string) *PermissionsResource {
	return &PermissionsResource{
		PermissionsID: permissionsID,
		Method:        method,
		Path:          path,
	}
}

// UpdatePath 更新资源路径
func (p *PermissionsResource) UpdatePath(path string) {
	p.Path = path
}

// UpdateMethod 更新HTTP方法
func (p *PermissionsResource) UpdateMethod(method string) {
	p.Method = method
}

// IsMatch 检查请求是否匹配该资源
func (p *PermissionsResource) IsMatch(method, path string) bool {
	return p.Method == method && p.Path == path
}

// Validate 验证资源是否有效
func (p *PermissionsResource) Validate() bool {
	// 验证HTTP方法
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}

	if !validMethods[p.Method] {
		return false
	}

	// 验证路径
	if p.Path == "" {
		return false
	}

	return true
}

// Clone 克隆资源
func (p *PermissionsResource) Clone() *PermissionsResource {
	return &PermissionsResource{
		ID:            p.ID,
		PermissionsID: p.PermissionsID,
		Method:        p.Method,
		Path:          p.Path,
	}
}

// Equal 比较两个资源是否相等
func (p *PermissionsResource) Equal(other *PermissionsResource) bool {
	if other == nil {
		return false
	}

	return p.Method == other.Method &&
		p.Path == other.Path &&
		p.PermissionsID == other.PermissionsID
}

// String 返回资源的字符串表示
func (p *PermissionsResource) String() string {
	return p.Method + " " + p.Path
}

// ResourceType 资源类型
type ResourceType struct {
	Method string
	Path   string
}

// GetResourceType 获取资源类型
func (p *PermissionsResource) GetResourceType() ResourceType {
	return ResourceType{
		Method: p.Method,
		Path:   p.Path,
	}
}

// IsPublic 判断是否为公开资源
func (p *PermissionsResource) IsPublic() bool {
	// 可以根据实际需求定义公开资源的规则
	publicPaths := map[string]bool{
		"/api/v1/public/*": true,
		"/api/v1/health":   true,
		"/api/v1/login":    true,
	}

	return publicPaths[p.Path]
}

// NormalizeMethod 标准化HTTP方法
func (p *PermissionsResource) NormalizeMethod() {
	// 确保HTTP方法都是大写
	switch p.Method {
	case "get", "Get":
		p.Method = "GET"
	case "post", "Post":
		p.Method = "POST"
	case "put", "Put":
		p.Method = "PUT"
	case "delete", "Delete":
		p.Method = "DELETE"
	case "patch", "Patch":
		p.Method = "PATCH"
	case "head", "Head":
		p.Method = "HEAD"
	case "options", "Options":
		p.Method = "OPTIONS"
	}
}
