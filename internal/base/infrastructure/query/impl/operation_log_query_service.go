package impl

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type OperationLogQueryService struct {
	repo repository.IOperationLogRepo
}

func NewOperationLogQueryService(
	repo repository.IOperationLogRepo,
) *OperationLogQueryService {
	return &OperationLogQueryService{
		repo: repo,
	}
}

func (s *OperationLogQueryService) Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*dto.OperationLogDto, error) {
	// 查询操作日志
	logs, err := s.repo.Find(ctx, tenantID, month, qb)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return dto.ToOperationLogDtoList(logs), nil
}

func (s *OperationLogQueryService) Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error) {
	return s.repo.Count(ctx, tenantID, month, qb)
}
