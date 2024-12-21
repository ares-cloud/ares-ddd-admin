package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/dto"
)

// ToFileDto 将 File 模型转换为 FileDto
func ToFileDto(file *model.File) *dto.FileDto {
	return &dto.FileDto{
		ID:          file.ID,
		Name:        file.Name,
		Path:        file.Path,
		FolderID:    file.FolderID,
		Size:        file.Size,
		Type:        file.Type,
		StorageType: string(file.StorageType),
		URL:         file.URL,
		CreatedBy:   file.CreatedBy,
		CreatedAt:   file.CreatedAt,
		DeletedAt:   file.DeletedAt,
		DeletedBy:   file.DeletedBy,
		IsDeleted:   file.IsDeleted,
	}
}

// ToFolderDto 将 Folder 模型转换为 FolderDto
func ToFolderDto(folder *model.Folder) *dto.FolderDto {
	return &dto.FolderDto{
		ID:        folder.ID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		Path:      folder.Path,
		CreatedBy: folder.CreatedBy,
		CreatedAt: folder.CreatedAt,
	}
}

// ToFolderTreeDto 将 Folder 模型转换为 FolderTreeDto
func ToFolderTreeDto(folder *model.Folder) *dto.FolderTreeDto {
	return &dto.FolderTreeDto{
		FolderDto: *ToFolderDto(folder),
		Children:  make([]*dto.FolderTreeDto, 0),
	}
}

// ToFileShareDto 将 FileShare 模型转换为 FileShareDto
func ToFileShareDto(share *model.FileShare) *dto.FileShareDto {
	return &dto.FileShareDto{
		ID:         share.ID,
		FileID:     share.FileID,
		ShareCode:  share.ShareCode,
		ExpireTime: share.ExpireTime,
		CreatedBy:  share.CreatedBy,
		CreatedAt:  share.CreatedAt,
	}
}

// ToFileDtoList 将 File 模型列表转换为 FileDto 列表
func ToFileDtoList(files []*model.File) []*dto.FileDto {
	dtos := make([]*dto.FileDto, len(files))
	for i, file := range files {
		dtos[i] = ToFileDto(file)
	}
	return dtos
}

// ToFolderDtoList 将 Folder 模型列表转换为 FolderDto 列表
func ToFolderDtoList(folders []*model.Folder) []*dto.FolderDto {
	dtos := make([]*dto.FolderDto, len(folders))
	for i, folder := range folders {
		dtos[i] = ToFolderDto(folder)
	}
	return dtos
}
