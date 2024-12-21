package entity

import "github.com/ares-cloud/ares-ddd-admin/pkg/database"

// File 文件实体
type File struct {
	database.BaseIntTime
	ID            string `gorm:"primaryKey;type:varchar(32)"` // ID
	Name          string `gorm:"type:varchar(255)"`           // 原始文件名
	GeneratedName string `gorm:"type:varchar(255)"`           // 生成的唯一文件名
	Path          string `gorm:"type:varchar(255)"`           // 文件路径
	OriginalPath  string `gorm:"type:varchar(255)"`           // 原始路径
	FolderID      string `gorm:"type:varchar(32)"`            // 文件夹ID
	Size          int64  `gorm:"type:bigint"`                 // 文件大小
	Type          string `gorm:"type:varchar(32)"`            // 文件类型
	StorageType   string `gorm:"type:varchar(32)"`            // 存储类型
	URL           string `gorm:"type:varchar(255)"`           // 访问URL
	TenantID      string `gorm:"type:varchar(32)"`            // 租户ID
	CreatedBy     string `gorm:"type:varchar(32)"`            // 创建人
	CreatedAt     int64  `gorm:"type:bigint"`                 // 创建时间
	UpdatedAt     int64  `gorm:"type:bigint"`                 // 更新时间
	DeletedBy     string `gorm:"type:varchar(32)"`            // 删除人
	DeletedAt     int64  `gorm:"type:bigint"`                 // 删除时间
	IsDeleted     bool   `gorm:"type:boolean"`                // 是否已删除
}

// Folder 文件夹实体
type Folder struct {
	database.BaseIntTime
	ID        string `gorm:"primaryKey;type:varchar(32)"`
	Name      string `gorm:"type:varchar(255)"`
	ParentID  string `gorm:"type:varchar(32)"`
	Path      string `gorm:"type:varchar(255)"`
	TenantID  string `gorm:"type:varchar(32)"`
	CreatedBy string `gorm:"type:varchar(32)"`
	CreatedAt int64  `gorm:"type:bigint"`
	UpdatedAt int64  `gorm:"type:bigint"`
}

// FileShare 文件分享实体
type FileShare struct {
	database.BaseIntTime
	ID         string `gorm:"primaryKey;type:varchar(32)"`
	FileID     string `gorm:"type:varchar(32)"`
	ShareCode  string `gorm:"type:varchar(32)"`
	Password   string `gorm:"type:varchar(32)"`
	ExpireTime int64  `gorm:"type:bigint"`
	CreatedBy  string `gorm:"type:varchar(32)"`
	CreatedAt  int64  `gorm:"type:bigint"`
}
