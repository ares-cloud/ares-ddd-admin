package oplog

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"
)

type DbOperationLogWriter struct {
	repo repository.IOperationLogRepository
}

func NewDbOperationLogWriter(repo repository.IOperationLogRepository) oplog.IDbOperationLogWrite {
	return &DbOperationLogWriter{
		repo: repo,
	}
}

func (w *DbOperationLogWriter) Save(ctx context.Context, data *oplog.OperationLog) error {
	log := &model.OperationLog{
		UserID:    data.UserID,
		Username:  data.Username,
		TenantID:  data.TenantID,
		Method:    data.Method,
		Path:      data.Path,
		Query:     data.Query,
		Body:      data.Body,
		IP:        data.IP,
		UserAgent: data.UserAgent,
		Status:    data.Status,
		Error:     data.Error,
		Duration:  data.Duration,
		Module:    data.Module,
		Action:    data.Action,
		CreatedAt: data.CreatedAt.Unix(),
	}
	return w.repo.Create(ctx, log)
}
