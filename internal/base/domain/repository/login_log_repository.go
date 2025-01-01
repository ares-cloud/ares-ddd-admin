package repository

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type ILoginLogRepository interface {
	Create(ctx context.Context, log *model.LoginLog) error
}
