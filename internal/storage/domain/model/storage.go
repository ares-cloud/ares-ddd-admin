package model

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal   StorageType = "local"   // 本地存储
	StorageTypeMinio   StorageType = "minio"   // Minio存储
	StorageTypeAliyun  StorageType = "aliyun"  // 阿里云存储
	StorageTypeTencent StorageType = "tencent" // 腾讯云存储
	StorageTypeQiniu   StorageType = "qiniu"   // 七牛云存储
)

// File 文件模型
type File struct {
	ID            string      `json:"id"`             // ID
	Name          string      `json:"name"`           // 原始文件名
	GeneratedName string      `json:"generated_name"` // 生成的唯一文件名
	Path          string      `json:"path"`           // 文件路径
	OriginalPath  string      `json:"original_path"`  // 原始路径(用于回收站恢复)
	FolderID      string      `json:"folder_id"`      // 文件夹ID
	Size          int64       `json:"size"`           // 文件大小(字节)
	Type          string      `json:"type"`           // 文件类型
	StorageType   StorageType `json:"storage_type"`   // 存储类型
	URL           string      `json:"url"`            // 访问URL
	TenantID      string      `json:"tenant_id"`      // 租户ID
	CreatedBy     string      `json:"created_by"`     // 创建人
	CreatedAt     int64       `json:"created_at"`     // 创建时间
	UpdatedAt     int64       `json:"updated_at"`     // 更新时间
	DeletedBy     string      `json:"deleted_by"`     // 删除人
	DeletedAt     int64       `json:"deleted_at"`     // 删除时间
	IsDeleted     bool        `json:"is_deleted"`     // 是否已删除
}

// Folder 文件夹模型
type Folder struct {
	ID        string `json:"id"`         // ID
	Name      string `json:"name"`       // 文件夹名
	ParentID  string `json:"parent_id"`  // 父文件夹ID
	Path      string `json:"path"`       // 文件夹路径
	TenantID  string `json:"tenant_id"`  // 租户ID
	CreatedBy string `json:"created_by"` // 创建人
	CreatedAt int64  `json:"created_at"` // 创建时间
	UpdatedAt int64  `json:"updated_at"` // 更新时间
}

// FolderTree 文件夹树形结构
type FolderTree struct {
	*Folder
	Children []*FolderTree `json:"children"` // 子文件夹
}

// FileShare 文件分享
type FileShare struct {
	ID         string `json:"id"`          // ID
	FileID     string `json:"file_id"`     // 文件ID
	ShareCode  string `json:"share_code"`  // 分享码
	Password   string `json:"password"`    // 访问密码
	ExpireTime int64  `json:"expire_time"` // 过期时间
	CreatedBy  string `json:"created_by"`  // 创建人
	CreatedAt  int64  `json:"created_at"`  // 创建时间
}
