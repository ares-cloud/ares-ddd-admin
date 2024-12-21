package repository

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
)

type IAuthRepository interface {
	// FindByUsername 根据用户名查找认证信息
	FindByUsername(ctx context.Context, username string) (*model.Auth, error)

	// SaveCaptcha 保存验证码
	SaveCaptcha(ctx context.Context, key, code string, expiration time.Duration) error

	// ValidateCaptcha 验证验证码
	ValidateCaptcha(ctx context.Context, key, code string) (bool, error)

	// FindByUserID 根据用户ID查找认证信息
	FindByUserID(ctx context.Context, userID string) (*model.Auth, error)
}
