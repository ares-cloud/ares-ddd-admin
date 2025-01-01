package impl

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type LoginLogQueryService struct {
	repo repository.ILoginLogRepo
}

func NewLoginLogQueryService(
	repo repository.ILoginLogRepo,
) *LoginLogQueryService {
	return &LoginLogQueryService{
		repo: repo,
	}
}

func (s *LoginLogQueryService) Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*dto.LoginLogDto, error) {
	// 查询登录日志
	logs, err := s.repo.Find(ctx, tenantID, month, qb)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return dto.ToLoginLogDtoList(logs), nil
}

func (s *LoginLogQueryService) Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error) {
	return s.repo.Count(ctx, tenantID, month, qb)
}
