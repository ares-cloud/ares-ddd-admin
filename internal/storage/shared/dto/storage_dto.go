package dto

// FileDto 文件DTO
type FileDto struct {
	ID          string `json:"id"`          // ID
	Name        string `json:"name"`        // 文件名
	Path        string `json:"path"`        // 文件路径
	FolderID    string `json:"folderId"`    // 文件夹ID
	Size        int64  `json:"size"`        // 文件大小(字节)
	Type        string `json:"type"`        // 文件类型
	StorageType string `json:"storageType"` // 存储类型
	URL         string `json:"url"`         // 访问URL
	CreatedBy   string `json:"createdBy"`   // 创建人
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	DeletedAt   int64  `json:"deletedAt"`   // 删除时间
	DeletedBy   string `json:"deletedBy"`   // 删除人
	IsDeleted   bool   `json:"isDeleted"`   // 是否已删除
}

// FolderDto 文件夹DTO
type FolderDto struct {
	ID        string `json:"id"`        // ID
	Name      string `json:"name"`      // 文件夹名
	ParentID  string `json:"parentId"`  // 父文件夹ID
	Path      string `json:"path"`      // 文件夹路径
	CreatedBy string `json:"createdBy"` // 创建人
	CreatedAt int64  `json:"createdAt"` // 创建时间
}

// FolderTreeDto 文件夹树形结构DTO
type FolderTreeDto struct {
	FolderDto
	Children []*FolderTreeDto `json:"children"` // 子文件夹
}

// FileShareDto 文件分享DTO
type FileShareDto struct {
	ID         string `json:"id"`
	FileID     string `json:"fileId"`
	ShareCode  string `json:"shareCode"`
	ExpireTime int64  `json:"expireTime"`
	CreatedBy  string `json:"createdBy"`
	CreatedAt  int64  `json:"createdAt"`
}
