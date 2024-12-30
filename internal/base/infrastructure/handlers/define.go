package handlers

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type HandlerEvent struct {
	queryCache *cache.EventHandler
	uh         *UserEventHandler
	eventBus   *pkgEvent.EventBus
}

func NewHandlerEvent(eventBus *pkgEvent.EventBus, queryCache *cache.EventHandler, uh *UserEventHandler) *HandlerEvent {
	return &HandlerEvent{
		queryCache: queryCache,
		uh:         uh,
		eventBus:   eventBus,
	}
}

func (h *HandlerEvent) Register() {
	// 注册用户相关事件
	h.eventBus.Subscribe(events.UserCreated, h.uh)
	h.eventBus.Subscribe(events.UserUpdated, h.uh)
	h.eventBus.Subscribe(events.UserDeleted, h.uh)
	h.eventBus.Subscribe(events.UserStatusChanged, h.uh)
	h.eventBus.Subscribe(events.UserRoleChanged, h.uh)
	h.eventBus.Subscribe(events.UserPasswordChanged, h.uh)
	// 注册缓存相关事件
	h.eventBus.Subscribe(events.UserLoggedIn, h.queryCache)
	h.eventBus.Subscribe(events.UserRoleChanged, h.queryCache)
	h.eventBus.Subscribe(events.UserUpdated, h.queryCache)
	h.eventBus.Subscribe(events.UserDeleted, h.queryCache)
	h.eventBus.Subscribe(events.UserStatusChanged, h.queryCache)
}
