package events

import "context"

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

type IEventBus interface {
	// Subscribe 订阅事件
	Subscribe(eventName string, handler EventHandler) error
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error
}
