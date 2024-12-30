package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

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
func (s *UserCommandService) CreateUser(ctx context.Context, user *model.User) herrors.Herr {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if exists {
		return errors.UserExists(user.Username)
	}

	// 创建用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return herrors.NewServerHError(err)
	}

	// 发布用户创建事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserCreated)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// UpdateUser 更新用户
func (s *UserCommandService) UpdateUser(ctx context.Context, user *model.User) herrors.Herr {
	// 检查用户是否存在
	exists, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if exists == nil {
		return errors.UserNotFound(user.ID)
	}

	// 更新用户
	if err := s.userRepo.Update(ctx, user); err != nil {
		return herrors.NewServerHError(err)
	}

	// 发布用户更新事件
	event := domanevent.NewUserEvent(user.TenantID, user.ID, domanevent.UserUpdated)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// AssignRoles 分配角色
func (s *UserCommandService) AssignRoles(ctx context.Context, userID string, roleIDs []int64) herrors.Herr {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if user == nil {
		return errors.UserNotFound(userID)
	}

	// 分配角色
	if err := s.userRepo.AssignRoles(ctx, userID, roleIDs); err != nil {
		return herrors.NewServerHError(err)
	}

	// 发布角色分配事件
	event := domanevent.NewUserRoleEvent(user.TenantID, userID, roleIDs)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// DeleteUser 删除用户
func (s *UserCommandService) DeleteUser(ctx context.Context, userID string) herrors.Herr {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if user == nil {
		return errors.UserNotFound(userID)
	}

	// 删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return herrors.NewServerHError(err)
	}

	// 发布用户删除事件
	event := domanevent.NewUserEvent(user.TenantID, userID, domanevent.UserDeleted)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// BelongsToDepartment 检查用户是否属于指定部门
func (s *UserCommandService) BelongsToDepartment(ctx context.Context, userID string, deptID string) (bool, herrors.Herr) {
	// 1. 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, herrors.NewServerHError(err)
	}
	if user == nil {
		return false, errors.UserNotFound(userID)
	}

	// 2. 检查部门归属关系
	belongs, err := s.userRepo.BelongsToDepartment(ctx, userID, deptID)
	if err != nil {
		return false, herrors.NewServerHError(err)
	}

	return belongs, nil
}

// TransferUser 调动用户部门
func (s *UserCommandService) TransferUser(ctx context.Context, userID, fromDeptID, toDeptID string) herrors.Herr {
	// 1. 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if user == nil {
		return errors.UserNotFound(userID)
	}

	// 2. 检查用户是否被锁定
	if locked, reason := user.IsLocked(); locked {
		return errors.UserDisabled(reason)
	}

	// 3. 检查用户是否属于原部门
	belongs, err := s.userRepo.BelongsToDepartment(ctx, userID, fromDeptID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if !belongs {
		return errors.UserInvalidField("department", "user does not belong to the source department")
	}

	// 4. 执行部门调动
	if err := s.userRepo.TransferUser(ctx, userID, fromDeptID, toDeptID); err != nil {
		return herrors.NewServerHError(err)
	}

	// 5. 发布用户部门调动事件
	event := domanevent.NewUserTransferEvent(user.TenantID, userID, fromDeptID, toDeptID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}

// GetUser 获取用户信息
func (s *UserCommandService) GetUser(ctx context.Context, userID string) (*model.User, herrors.Herr) {
	// 1. 查询用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	if user == nil {
		return nil, errors.UserNotFound(userID)
	}

	// 2. 检查用户状态
	if locked, reason := user.IsLocked(); locked {
		return nil, errors.UserDisabled(reason)
	}

	return user, nil
}
