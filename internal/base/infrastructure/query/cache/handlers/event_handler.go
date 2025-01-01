package handlers

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	"log"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// EventHandler 缓存事件处理器
type EventHandler struct {
	userCache *cache.UserQueryCache
}

func NewCacheEventHandler(userCache *cache.UserQueryCache) *EventHandler {
	return &EventHandler{
		userCache: userCache,
	}
}

// Handle 处理事件
func (h *EventHandler) Handle(ctx context.Context, event pkgEvent.Event) error {
	switch e := event.(type) {
	case *events.UserEvent:
		return h.handleUserEvent(ctx, e)
	case *events.UserRoleEvent:
		return h.handleUserRoleEvent(ctx, e)
	case *events.UserLoginEvent:
		return h.handleUserLoginEvent(ctx, e)
	default:
		return nil
	}
}

// handleUserEvent 处理用户事件
func (h *EventHandler) handleUserEvent(ctx context.Context, event *events.UserEvent) error {
	log.Printf("处理用户事件: %s, 用户ID=%s", event.EventName(), event.UserID())

	// 所有用户相关事件都清除用户缓存
	if err := h.userCache.InvalidateUserCache(ctx, event.UserID()); err != nil {
		log.Printf("清除用户缓存失败: %v", err)
		return fmt.Errorf("清除用户缓存失败: %w", err)
	}

	// 特定事件的额外处理
	switch event.EventName() {
	case events.UserDeleted, events.UserStatusChanged:
		if err := h.userCache.InvalidateUserListCache(ctx); err != nil {
			log.Printf("清除用户列表缓存失败: %v", err)
			return fmt.Errorf("清除用户列表缓存失败: %w", err)
		}
	}

	return nil
}

// handleUserRoleEvent 处理用户角色事件
func (h *EventHandler) handleUserRoleEvent(ctx context.Context, event *events.UserRoleEvent) error {
	log.Printf("处理用户角色事件: 用户ID=%s, 角色IDs=%v", event.UserID(), event.RoleIDs())

	// 清除用户权限相关缓存
	if err := h.userCache.InvalidateUserPermissionCache(ctx, event.UserID()); err != nil {
		return err
	}

	// 清除角色相关列表缓存
	if err := h.userCache.InvalidateRoleListCache(ctx); err != nil {
		return err
	}

	return nil
}

// handleUserLoginEvent 处理用户登录事件
func (h *EventHandler) handleUserLoginEvent(ctx context.Context, event *events.UserLoginEvent) error {
	log.Printf("处理用户登录事件: 用户ID=%s", event.UserID())

	// 预热用户缓存
	if err := h.userCache.WarmupUserCache(ctx, event.UserID()); err != nil {
		log.Printf("预热用户缓存失败: %v", err)
		return fmt.Errorf("预热用户缓存失败: %w", err)
	}

	return nil
}
