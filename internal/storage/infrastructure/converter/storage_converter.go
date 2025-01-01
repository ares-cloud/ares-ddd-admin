package converter

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/entity"
)

// ToFileDto 将 File 模型转换为 FileDto
func ToFileDto(file *entity.File) *dto.FileDto {
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
func ToFolderDto(folder *entity.Folder) *dto.FolderDto {
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
func ToFolderTreeDto(folder *dto.FolderDto) *dto.FolderTreeDto {
	return &dto.FolderTreeDto{
		FolderDto: *folder,
		Children:  make([]*dto.FolderTreeDto, 0),
	}
}

// ToFileShareDto 将 FileShare 模型转换为 FileShareDto
func ToFileShareDto(share *entity.FileShare) *dto.FileShareDto {
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
func ToFileDtoList(files []*entity.File) []*dto.FileDto {
	dtos := make([]*dto.FileDto, len(files))
	for i, file := range files {
		dtos[i] = ToFileDto(file)
	}
	return dtos
}

// ToFolderDtoList 将 Folder 模型列表转换为 FolderDto 列表
func ToFolderDtoList(folders []*entity.Folder) []*dto.FolderDto {
	dtos := make([]*dto.FolderDto, len(folders))
	for i, folder := range folders {
		dtos[i] = ToFolderDto(folder)
	}
	return dtos
}
