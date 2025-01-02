package events

import (
	"context"
	"sync"
)

// DefEventBus 事件总线
type DefEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() IEventBus {
	return &DefEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (bus *DefEventBus) Subscribe(eventName string, handler EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
	return nil
}

// Publish 发布事件
func (bus *DefEventBus) Publish(ctx context.Context, event Event) error {
	bus.mu.RLock()
	handlers := bus.handlers[event.EventName()]
	bus.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
