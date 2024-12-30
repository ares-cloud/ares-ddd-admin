package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/events"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	domanevent "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
)

type UserCommandService struct {
	userRepo repository.IUserRepository
	eventBus *events.EventBus
}

func NewUserCommandService(
	userRepo repository.IUserRepository,
	eventBus *events.EventBus,
) *UserCommandService {
	return &UserCommandService{
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// CreateUser 创建用户
func (s *UserCommandService) CreateUser(ctx context.Context, user *model.User) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUsernameExists
	}

	// 创建用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// 发布用户创建事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserCreated)
	return s.eventBus.Publish(ctx, event)
}

// UpdateUser 更新用户
func (s *UserCommandService) UpdateUser(ctx context.Context, user *model.User) error {
	// 检查用户是否存在
	exists, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return errors.ErrUserNotFound
	}

	// 更新用户
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 发布用户更新事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserUpdated)
	return s.eventBus.Publish(ctx, event)
}

// AssignRoles 分配角色
func (s *UserCommandService) AssignRoles(ctx context.Context, userID string, roleIDs []int64) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 分配角色
	if err := s.userRepo.AssignRoles(ctx, userID, roleIDs); err != nil {
		return err
	}

	// 发布角色分配事件
	event := domanevent.NewUserRoleEvent(user.TenantID, userID, roleIDs)
	return s.eventBus.Publish(ctx, event)
}

// DeleteUser 删除用户
func (s *UserCommandService) DeleteUser(ctx context.Context, userID string) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// 发布用户删除事件
	event := domanevent.NewUserEvent(user.TenantID, userID, domanevent.UserDeleted)
	return s.eventBus.Publish(ctx, event)
}
