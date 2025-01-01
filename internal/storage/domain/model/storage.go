package model

import (
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"path/filepath"
	"strings"
	"time"
)

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

// Validate 验证文件属性
func (f *File) Validate() error {
	// 1. 检查必填字段
	if f.Name == "" {
		return fmt.Errorf("file name is required")
	}
	if f.TenantID == "" {
		return fmt.Errorf("tenant id is required")
	}
	if f.CreatedBy == "" {
		return fmt.Errorf("creator is required")
	}

	// 2. 检查文件大小
	if f.Size <= 0 {
		return fmt.Errorf("file size must be greater than 0")
	}

	// 3. 检查存储类型
	if f.StorageType == "" {
		return fmt.Errorf("storage type is required")
	}
	if !isValidStorageType(f.StorageType) {
		return fmt.Errorf("invalid storage type: %s", f.StorageType)
	}

	// 4. 检查文件名格式
	if err := validateFileName(f.Name); err != nil {
		return err
	}

	// 5. 检查文件类型
	if f.Type == "" {
		f.Type = getFileType(f.Name)
	}

	return nil
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

// Validate 验证文件夹属性
func (f *Folder) Validate() error {
	// 1. 检查必填字段
	if f.Name == "" {
		return fmt.Errorf("folder name is required")
	}
	if f.TenantID == "" {
		return fmt.Errorf("tenant id is required")
	}
	if f.CreatedBy == "" {
		return fmt.Errorf("creator is required")
	}

	// 2. 检查文件夹名格式
	if err := validateFolderName(f.Name); err != nil {
		return err
	}

	// 3. 检查父文件夹ID
	if f.ParentID == "" {
		f.ParentID = "0" // 默认为根目录
	}

	return nil
}

func (f *Folder) SetName(name string) error {
	err := validateFolderName(name)
	if herrors.HaveError(err) {
		return err
	}
	f.Name = name
	f.UpdatedAt = time.Now().Unix()
	return nil
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

// Validate 验证文件分享属性
func (s *FileShare) Validate() error {
	// 1. 检查必填字段
	if s.FileID == "" {
		return fmt.Errorf("file id is required")
	}
	if s.ShareCode == "" {
		return fmt.Errorf("share code is required")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("creator is required")
	}

	// 2. 检查过期时间
	if s.ExpireTime <= 0 {
		return fmt.Errorf("expire time must be greater than 0")
	}

	return nil
}

// 辅助函数

// isValidStorageType 检查存储类型是否有效
func isValidStorageType(st StorageType) bool {
	switch st {
	case StorageTypeLocal, StorageTypeMinio, StorageTypeAliyun, StorageTypeTencent, StorageTypeQiniu:
		return true
	default:
		return false
	}
}

// validateFileName 验证文件名
func validateFileName(name string) error {
	// 1. 检查文件名长度
	if len(name) > 255 {
		return fmt.Errorf("file name too long")
	}

	// 2. 检查文件名是否包含非法字符
	if strings.ContainsAny(name, "\\/:*?\"<>|") {
		return fmt.Errorf("file name contains invalid characters")
	}

	return nil
}

// validateFolderName 验证文件夹名
func validateFolderName(name string) error {
	// 1. 检查文件夹名长度
	if len(name) > 255 {
		return fmt.Errorf("folder name too long")
	}

	// 2. 检查文件夹名是否包含非法字符
	if strings.ContainsAny(name, "\\/:*?\"<>|") {
		return fmt.Errorf("folder name contains invalid characters")
	}

	return nil
}

// getFileType 根据文件名获取文件类型
func getFileType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return "unknown"
	}
	return ext[1:] // 去掉点号
}
