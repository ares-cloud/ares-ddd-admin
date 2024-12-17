package repository

import (
	"context"
	"errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/mapper"
)

type authRepository struct {
	userRepo repository.IUserRepository
	rdb      *redis.Client
	mapper   *mapper.UserMapper
}

func NewAuthRepository(userRepo repository.IUserRepository, rdb *h_redis.RedisClient) repository.IAuthRepository {
	return &authRepository{
		userRepo: userRepo,
		rdb:      rdb.GetClient(),
		mapper:   &mapper.UserMapper{},
	}
}

func (r *authRepository) FindByUsername(ctx context.Context, username string) (*model.Auth, error) {
	user, err := r.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return model.NewAuth(user, "web"), nil
}

func (r *authRepository) SaveCaptcha(ctx context.Context, key, code string, expiration time.Duration) error {
	return r.rdb.Set(ctx, "captcha:"+key, code, expiration).Err()
}

func (r *authRepository) ValidateCaptcha(ctx context.Context, key, code string) (bool, error) {
	storedCode, err := r.rdb.Get(ctx, "captcha:"+key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	// 验证后删除验证码
	r.rdb.Del(ctx, "captcha:"+key)

	return storedCode == code, nil
}

func (r *authRepository) FindByUserID(ctx context.Context, userID string) (*model.Auth, error) {
	user, err := r.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.NewAuth(user, "web"), nil
}
