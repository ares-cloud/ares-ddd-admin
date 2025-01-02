package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	pkgEvents "github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// UserEventHandler 用户事件处理器
type UserEventHandler struct {
}

func NewUserEventHandler() *UserEventHandler {
	return &UserEventHandler{}
}

// Handle 处理事件
func (h *UserEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.UserEvent:
		return h.handleUserEvent(ctx, e)
	default:
		return nil
	}
}

// handleUserEvent 处理用户基础事件
func (h *UserEventHandler) handleUserEvent(ctx context.Context, event *events.UserEvent) error {
	switch event.EventName() {
	case events.UserCreated:
		hlog.CtxDebugf(ctx, "用户创建事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
		//return h.handleUserCreated(ctx, event)
	case events.UserUpdated:
		hlog.CtxDebugf(ctx, "用户更新事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
		//return h.handleUserUpdated(ctx, event)
	case events.UserDeleted:
		hlog.CtxDebugf(ctx, "用户删除事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
		//return h.handleUserDeleted(ctx, event)
	}
	return nil
}

// handleUserDeleted 处理用户删除事件
func (h *UserEventHandler) handleUserDeleted(ctx context.Context, event *events.UserEvent) error {
	return nil
}

// cleanUserCache 清除用户缓存
func (h *UserEventHandler) cleanUserCache(ctx context.Context, userID string) error {
	return nil
}

// cleanUserRelatedCache 清除用户相关的其他数据缓存
func (h *UserEventHandler) cleanUserRelatedCache(ctx context.Context, userID string) error {

	return nil
}
