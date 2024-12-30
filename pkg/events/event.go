package events

import (
	"context"
	"sync"
	"time"
)

// Event 事件接口
type Event interface {
	// EventName 事件名称
	EventName() string
	// EventTime 事件发生时间
	EventTime() int64
}

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event Event) error
}

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (bus *EventBus) Subscribe(eventName string, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
}

// Publish 发布事件
func (bus *EventBus) Publish(ctx context.Context, event Event) error {
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

// BaseEvent 基础事件
type BaseEvent struct {
	eventName string
	eventTime int64
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		eventName: name,
		eventTime: time.Now().UnixNano(),
	}
}

func (e *BaseEvent) EventName() string {
	return e.eventName
}

func (e *BaseEvent) EventTime() int64 {
	return e.eventTime
}
