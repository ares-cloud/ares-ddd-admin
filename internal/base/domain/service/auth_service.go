package service

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	domanevent "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
)

type AuthService struct {
	userRepo     repository.IUserRepository
	eventBus     events.IEventBus
	queryService query.IUserQueryService
}

func NewAuthService(
	userRepo repository.IUserRepository,
	eventBus events.IEventBus,
	queryService query.IUserQueryService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		eventBus:     eventBus,
		queryService: queryService,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.User, error) {
	// 获取用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 验证密码
	if her := user.ComparePassword(password); herrors.HaveError(her) {
		return nil, errors.ErrInvalidCredentials
	}

	// 检查用户状态
	if ok, _ := user.IsLocked(); !ok {
		return nil, errors.ErrUserDisabled
	}

	// 发布登录事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserLoggedIn)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if herrors.HaveError(user.ComparePassword(oldPassword)) {
		return errors.ErrInvalidCredentials
	}
	user.Password = newPassword
	// 更新密码
	if err := user.HashPassword(); err != nil {
		return err
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 发布密码修改事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserUpdated)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return err
	}
	return nil
}
