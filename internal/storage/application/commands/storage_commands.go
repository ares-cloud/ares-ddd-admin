package commands

// CreateFolderCommand 创建文件夹命令
type CreateFolderCommand struct {
	Name     string `json:"name"`      // 文件夹名称
	ParentID string `json:"parent_id"` // 父文件夹ID
}

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	FolderID string `form:"folder_id"`                    // 文件夹ID
	File     []byte `form:"file" binding:"required"`      // 文件内容
	FileName string `form:"file_name" binding:"required"` // 文件名
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
	FolderID string `json:"folder_id"`    // 目标文件夹ID
}

// BatchDeleteFilesCommand 批量删除文件命令
type BatchDeleteFilesCommand struct {
	FileIDs []string `json:"file_ids"` // 文件ID列表
}

// BatchMoveFilesCommand 批量移动文件命令
type BatchMoveFilesCommand struct {
	FileIDs        []string `json:"file_ids"`         // 文件ID列表
	TargetFolderID string   `json:"target_folder_id"` // 目标文件夹ID
}

// ShareFileCommand 分享文件命令
type ShareFileCommand struct {
	FileID     string `json:"file_id"`     // 文件ID
	Password   string `json:"password"`    // 访问密码
	ExpireTime int64  `json:"expire_time"` // 过期时间
	CreatedBy  string `json:"-"`           // 创建人
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
	ID             string `json:"id" path:"id"`     // 文件夹ID
	TargetParentID string `json:"target_parent_id"` // 目标父文件夹ID
}
