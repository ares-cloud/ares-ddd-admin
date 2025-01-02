package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 用户事件类型定义
const (
	UserCreated     = "user.created"
	UserUpdated     = "user.updated"
	UserDeleted     = "user.deleted"
	UserRoleChanged = "user.role.changed"
	UserLoggedIn    = "user.logged_in"
)

// UserEvent 用户事件
type UserEvent struct {
	events.BaseEvent
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

// NewUserEvent 创建用户事件
func NewUserEvent(tenantID, userID string, eventName string) *UserEvent {
	return &UserEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		TenantID:  tenantID,
		UserID:    userID,
	}
}
