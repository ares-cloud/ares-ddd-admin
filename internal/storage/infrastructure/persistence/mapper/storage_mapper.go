package mapper

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/dto"
)

type StorageMapper struct{}

// ToFileEntity 文件领域模型转实体
func (m *StorageMapper) ToFileEntity(domain *model.File) *entity.File {
	return &entity.File{
		ID:           domain.ID,
		Name:         domain.Name,
		Path:         domain.Path,
		OriginalPath: domain.OriginalPath,
		FolderID:     domain.FolderID,
		Size:         domain.Size,
		Type:         domain.Type,
		StorageType:  string(domain.StorageType),
		URL:          domain.URL,
		TenantID:     domain.TenantID,
		CreatedBy:    domain.CreatedBy,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
		DeletedBy:    domain.DeletedBy,
		DeletedAt:    domain.DeletedAt,
		IsDeleted:    domain.IsDeleted,
	}
}

// ToFileDomain 文件实体转领域模型
func (m *StorageMapper) ToFileDomain(ent *entity.File) *model.File {
	return &model.File{
		ID:           ent.ID,
		Name:         ent.Name,
		Path:         ent.Path,
		OriginalPath: ent.OriginalPath,
		FolderID:     ent.FolderID,
		Size:         ent.Size,
		Type:         ent.Type,
		StorageType:  model.StorageType(ent.StorageType),
		URL:          ent.URL,
		TenantID:     ent.TenantID,
		CreatedBy:    ent.CreatedBy,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
		DeletedBy:    ent.DeletedBy,
		DeletedAt:    ent.DeletedAt,
		IsDeleted:    ent.IsDeleted,
	}
}

// ToFileDomainList 文件实体列表转领域模型列表
func (m *StorageMapper) ToFileDomainList(entities []*entity.File) []*model.File {
	if len(entities) == 0 {
		return nil
	}
	domains := make([]*model.File, len(entities))
	for i, ent := range entities {
		domains[i] = m.ToFileDomain(ent)
	}
	return domains
}

// ToFolderEntity 文件夹领域模型转实体
func (m *StorageMapper) ToFolderEntity(domain *model.Folder) *entity.Folder {
	return &entity.Folder{
		ID:        domain.ID,
		Name:      domain.Name,
		ParentID:  domain.ParentID,
		Path:      domain.Path,
		CreatedBy: domain.CreatedBy,
		TenantID:  domain.TenantID,
	}
}

// ToFolderDomain 文件夹实体转领域模型
func (m *StorageMapper) ToFolderDomain(ent *entity.Folder) *model.Folder {
	return &model.Folder{
		ID:        ent.ID,
		Name:      ent.Name,
		ParentID:  ent.ParentID,
		Path:      ent.Path,
		CreatedBy: ent.CreatedBy,
		TenantID:  ent.TenantID,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}

// ToFolderDomainList 文件夹实体列表转领域模型列表
func (m *StorageMapper) ToFolderDomainList(entities []*entity.Folder) []*model.Folder {
	if len(entities) == 0 {
		return nil
	}
	domains := make([]*model.Folder, len(entities))
	for i, ent := range entities {
		domains[i] = m.ToFolderDomain(ent)
	}
	return domains
}

// ToFileShareEntity 文件分享领域模型转实体
func (m *StorageMapper) ToFileShareEntity(domain *model.FileShare) *entity.FileShare {
	return &entity.FileShare{
		ID:         domain.ID,
		FileID:     domain.FileID,
		ShareCode:  domain.ShareCode,
		Password:   domain.Password,
		ExpireTime: domain.ExpireTime,
		CreatedBy:  domain.CreatedBy,
	}
}

// ToFileShareDomain 文件分享实体转领域模型
func (m *StorageMapper) ToFileShareDomain(ent *entity.FileShare) *model.FileShare {
	return &model.FileShare{
		ID:         ent.ID,
		FileID:     ent.FileID,
		ShareCode:  ent.ShareCode,
		Password:   ent.Password,
		ExpireTime: ent.ExpireTime,
		CreatedAt:  ent.CreatedAt,
		CreatedBy:  ent.CreatedBy,
	}
}

// ToFolderTreeDto 将文件夹转换为树形结构DTO
func (m *StorageMapper) ToFolderTreeDto(folder *model.Folder) *dto.FolderTreeDto {
	return &dto.FolderTreeDto{
		FolderDto: dto.FolderDto{
			ID:        folder.ID,
			Name:      folder.Name,
			ParentID:  folder.ParentID,
			Path:      folder.Path,
			CreatedBy: folder.CreatedBy,
			CreatedAt: folder.CreatedAt,
		},
		Children: make([]*dto.FolderTreeDto, 0),
	}
}
