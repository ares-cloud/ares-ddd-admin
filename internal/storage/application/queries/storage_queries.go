package queries

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
)

// ListFoldersQuery 查询文件夹列表
type ListFoldersQuery struct {
	query.Page
	ParentID string `query:"parent_id"` // 父文件夹ID
	Name     string `query:"name"`      // 文件夹名称
	TenantID string `query:"tenant_id"` // 租户ID
}

// ListFilesQuery 查询文件列表
type ListFilesQuery struct {
	query.Page
	FolderID    string `query:"folder_id"`    // 文件夹ID
	Name        string `query:"name"`         // 文件名
	Type        string `query:"type"`         // 文件类型
	StorageType string `query:"storage_type"` // 存储类型
	TenantID    string `query:"tenant_id"`    // 租户ID
}

// CreateFolderCommand 创建文件夹
type CreateFolderCommand struct {
	Name     string `json:"name"`      // 文件夹名称
	ParentID int64  `json:"parent_id"` // 父文件夹ID
}

// GetFolderTreeQuery 获取文件夹树形结构查询
type GetFolderTreeQuery struct {
	TenantID string `query:"tenant_id"` // 租户ID
}

// ListRecycleFilesQuery 查询回收站文件列表查询
type ListRecycleFilesQuery struct {
	query.Page
	Name        string `query:"name"`         // 文件名
	Type        string `query:"type"`         // 文件类型
	StorageType string `query:"storage_type"` // 存储类型
	TenantID    string `query:"-"`            // 租户ID
}

// GetFileShareQuery 获取文件分享信息
type GetFileShareQuery struct {
	ShareCode string `query:"share_code"` // 分享码
	Password  string `query:"password"`   // 访问密码
}

// GetFilePreviewQuery 获取文件预览查询
type GetFilePreviewQuery struct {
	ID string `query:"id" path:"id"` // 文件ID
}

// GetFileDownloadQuery 获取文件下载查询
type GetFileDownloadQuery struct {
	ID string `query:"id" path:"id"` // 文件ID
}
