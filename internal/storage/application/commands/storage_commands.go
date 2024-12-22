package commands

// CreateFolderCommand 创建文件夹命令
type CreateFolderCommand struct {
	Name     string `json:"name"`     // 文件夹名称
	ParentID string `json:"parentId"` // 父文件夹ID
}

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	FolderID string `form:"folderId"`                    // 文件夹ID
	File     []byte `form:"file" binding:"required"`     // 文件内容
	FileName string `form:"fileName" binding:"required"` // 文件名
}

// DeleteFileCommand 删除文件命令
type DeleteFileCommand struct {
	ID string `json:"id" path:"id"` // 文件ID
}

// DeleteFolderCommand 删除文件夹命令
type DeleteFolderCommand struct {
	ID string `json:"id" path:"id"` // 文件夹ID
}

// RenameFolderCommand 重命名文件夹命令
type RenameFolderCommand struct {
	ID   string `json:"id" path:"id"` // 文件夹ID
	Name string `json:"name"`         // 新名称
}

// MoveFileCommand 移动文件命令
type MoveFileCommand struct {
	ID       string `json:"id" path:"id"` // 文件ID
	FolderID string `json:"folderId"`     // 目标文件夹ID
}

// BatchDeleteFilesCommand 批量删除文件命令
type BatchDeleteFilesCommand struct {
	FileIDs []string `json:"fileIds"` // 文件ID列表
}

// BatchMoveFilesCommand 批量移动文件命令
type BatchMoveFilesCommand struct {
	FileIDs        []string `json:"fileIds"`        // 文件ID列表
	TargetFolderID string   `json:"targetFolderId"` // 目标文件夹ID
}

// ShareFileCommand 分享文件命令
type ShareFileCommand struct {
	FileID     string `json:"fileId"`     // 文件ID
	Password   string `json:"password"`   // 访问密码
	ExpireTime int64  `json:"expireTime"` // 过期时间
	CreatedBy  string `json:"-"`          // 创建人
}

// RecycleFileCommand 回收文件命令
type RecycleFileCommand struct {
	ID        string `json:"id" path:"id"` // 文件ID
	DeletedBy string `json:"-"`            // 删除人
}

// RestoreFileCommand 恢复文件命令
type RestoreFileCommand struct {
	ID string `json:"id" path:"id"` // 文件ID
}

// MoveFolderCommand 移动文件夹命令
type MoveFolderCommand struct {
	ID             string `json:"id" path:"id"`   // 文件夹ID
	TargetParentID string `json:"targetParentId"` // 目标父文件夹ID
}
