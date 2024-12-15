package commands

type CreatePermissionsCommand struct {
	Code        string                             `json:"code" binding:"required"`
	Name        string                             `json:"name" binding:"required"`
	Localize    string                             `json:"localize"`
	Icon        string                             `json:"icon"`
	Description string                             `json:"description"`
	Sequence    int                                `json:"sequence"`
	Type        int8                               `json:"type" binding:"required"`
	Path        string                             `json:"path"`
	Properties  string                             `json:"properties"`
	ParentID    int64                              `json:"parentId"`
	Resources   []CreatePermissionsResourceCommand `json:"resources"`
}

type CreatePermissionsResourceCommand struct {
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

type UpdatePermissionsCommand struct {
	ID          int64  `json:"id" binding:"required"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
	Path        string `json:"path"`
	Properties  string `json:"properties"`
	Status      *int8  `json:"status"`
}

// DeletePermissionsCommand 删除权限命令
type DeletePermissionsCommand struct {
	ID int64 `json:"id" path:"id"` // 权限ID
}
