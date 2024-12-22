package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/converter"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type StorageQueryHandler struct {
	repo repository.IStorageRepository
}

func NewStorageQueryHandler(repo repository.IStorageRepository) *StorageQueryHandler {
	return &StorageQueryHandler{
		repo: repo,
	}
}

// HandleListFolders 处理查询文件夹列表
func (h *StorageQueryHandler) HandleListFolders(ctx context.Context, q *queries.ListFoldersQuery) (*models.PageRes[dto.FolderDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	if q.TenantID != "" {
		qb.Where("tenant_id", query.Eq, q.TenantID)
	}
	qb.WithPage(&q.Page)

	// 查询数据
	folders, total, err := h.repo.ListFolders(ctx, q.ParentID, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtos := converter.ToFolderDtoList(folders)

	return &models.PageRes[dto.FolderDto]{
		Total: total,
		List:  dtos,
	}, nil
}

// HandleListFiles 处理查询文件列表
func (h *StorageQueryHandler) HandleListFiles(ctx context.Context, q *queries.ListFilesQuery) (*models.PageRes[dto.FileDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	if q.Type != "" {
		qb.Where("type", query.Eq, q.Type)
	}
	if q.StorageType != "" {
		qb.Where("storage_type", query.Eq, q.StorageType)
	}
	if q.TenantID != "" {
		qb.Where("tenant_id", query.Eq, q.TenantID)
	}
	qb.WithPage(&q.Page)

	// 查询数据
	files, total, err := h.repo.ListFiles(ctx, q.FolderID, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtos := converter.ToFileDtoList(files)

	return &models.PageRes[dto.FileDto]{
		List:  dtos,
		Total: total,
	}, nil
}

// HandleGetFolderTree 处理获取文件夹树形结构
func (h *StorageQueryHandler) HandleGetFolderTree(ctx context.Context) ([]*dto.FolderTreeDto, herrors.Herr) {
	// 1. 获取所有文件夹
	qb := query.NewQueryBuilder()
	folders, _, err := h.repo.ListFolders(ctx, "0", qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 2. 构建文件夹映射
	folderMap := make(map[string]*dto.FolderTreeDto)
	var roots []*dto.FolderTreeDto

	// 3. 转换为树形DTO
	for _, folder := range folders {
		dto := converter.ToFolderTreeDto(folder)
		folderMap[folder.ID] = dto

		if folder.ParentID == "0" || folder.ParentID == "" {
			roots = append(roots, dto)
		} else {
			if parent, ok := folderMap[folder.ParentID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		}
	}

	return roots, nil
}

// HandleListRecycleFiles 处理查询回收站文件列表
func (h *StorageQueryHandler) HandleListRecycleFiles(ctx context.Context, q *queries.ListRecycleFilesQuery) (*models.PageRes[dto.FileDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Name != "" {
		qb.Where("name", query.Like, "%"+q.Name+"%")
	}
	if q.Type != "" {
		qb.Where("type", query.Eq, q.Type)
	}
	if q.StorageType != "" {
		qb.Where("storage_type", query.Eq, q.StorageType)
	}
	if q.TenantID != "" {
		qb.Where("tenant_id", query.Eq, q.TenantID)
	}
	// 只查询已删除的文件
	qb.Where("is_deleted", query.Eq, true)
	qb.WithPage(&q.Page)

	// 查询数据
	files, total, err := h.repo.ListFiles(ctx, "", qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtos := converter.ToFileDtoList(files)

	return &models.PageRes[dto.FileDto]{
		List:  dtos,
		Total: total,
	}, nil
}
