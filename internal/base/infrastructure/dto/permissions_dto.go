package dto

// PermissionsDto 权限数据传输对象
type PermissionsDto struct {
	ID          int64                     `json:"id"`                 // ID
	Code        string                    `json:"code"`               // 编码
	Name        string                    `json:"name"`               // 名称
	Localize    string                    `json:"localize"`           // 本地化
	Icon        string                    `json:"icon"`               // 图标
	Description string                    `json:"description"`        // 描述
	Sequence    int                       `json:"sequence"`           // 排序
	Type        int8                      `json:"type"`               // 类型
	Path        string                    `json:"path"`               // 路径
	Properties  string                    `json:"properties"`         // 属性
	Status      int8                      `json:"status"`             // 状态
	ParentID    int64                     `json:"parentId"`           // 父级ID
	ParentPath  string                    `json:"parentPath"`         // 父级路径
	Resources   []*PermissionsResourceDto `json:"resources"`          // 资源列表
	Children    []*PermissionsDto         `json:"children,omitempty"` // 子权限
	CreatedAt   int64                     `json:"createdAt"`          // 创建时间
	UpdatedAt   int64                     `json:"updatedAt"`          // 更新时间
}

// PermissionsResourceDto 权限资源数据传输对象
type PermissionsResourceDto struct {
	Method string `json:"method"` // HTTP方法
	Path   string `json:"path"`   // 资源路径
}
