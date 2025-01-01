package impl

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type StorageQueryService struct {
	repo    repository.IStorageRepos
	storage storage.StorageFactory
}

func NewStorageQueryService(repo repository.IStorageRepos, storage storage.StorageFactory) query.IStorageQueryService {
	return &StorageQueryService{
		repo:    repo,
		storage: storage,
	}
}

// ListFiles 查询文件列表
func (s *StorageQueryService) ListFiles(ctx context.Context, folderID string, qb *db_query.QueryBuilder) ([]*dto.FileDto, int64, error) {
	// 1. 查询数据
	files, total, err := s.repo.ListFiles(ctx, folderID, qb)
	if err != nil {
		return nil, 0, err
	}

	// 2. 转换为DTO
	dtos := converter.ToFileDtoList(files)

	return dtos, total, nil
}

// ListFolders 查询文件夹列表
func (s *StorageQueryService) ListFolders(ctx context.Context, parentID string, qb *db_query.QueryBuilder) ([]*dto.FolderDto, int64, error) {
	// 1. 查询数据
	folders, total, err := s.repo.ListFolders(ctx, parentID, qb)
	if err != nil {
		return nil, 0, err
	}

	// 2. 转换为DTO
	dtos := converter.ToFolderDtoList(folders)

	return dtos, total, nil
}

// PreviewFile 预览文件

// GetShareFile 获取分享文件
func (s *StorageQueryService) GetShareFile(ctx context.Context, shareCode string, password string) (*dto.FileDto, error) {
	// 1. 获取分享信息
	share, err := s.repo.GetFileShare(ctx, shareCode)
	if err != nil {
		return nil, errors.ErrShareNotFound
	}

	// 2. 验证密码
	if share.Password != "" && share.Password != password {
		return nil, errors.ErrSharePasswordIncorrect
	}

	// 3. 验证是否过期
	if share.ExpireTime > 0 && share.ExpireTime < time.Now().Unix() {
		return nil, errors.ErrShareExpired
	}

	// 4. 获取文件信息
	file, err := s.repo.GetFile(ctx, share.FileID)
	if err != nil {
		return nil, err
	}

	// 5. 转换为DTO
	return converter.ToFileDto(file), nil
}
