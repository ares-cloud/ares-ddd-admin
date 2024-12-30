package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 用户事件类型定义
const (
	UserCreated         = "user.created"
	UserUpdated         = "user.updated"
	UserDeleted         = "user.deleted"
	UserStatusChanged   = "user.status.changed"
	UserRoleChanged     = "user.role.changed"
	UserLoggedIn        = "user.logged_in"
	UserPasswordChanged = "user.password.changed"
)

// UserEvent 用户事件
type UserEvent struct {
	events.BaseEvent
	tenantID string
	userID   string
}

// NewUserEvent 创建用户事件
func NewUserEvent(tenantID, userID string, eventName string) *UserEvent {
	return &UserEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		tenantID:  tenantID,
		userID:    userID,
	}
}

func (e *UserEvent) UserID() string {
	return e.userID
}

func (e *UserEvent) TenantID() string {
	return e.tenantID
}

// UserLoginEvent 用户登录事件
type UserLoginEvent struct {
	UserEvent
}

// NewUserLoginEvent 创建用户登录事件
func NewUserLoginEvent(tenantID, userID string) *UserLoginEvent {
	return &UserLoginEvent{
		UserEvent: *NewUserEvent(tenantID, userID, UserLoggedIn),
	}
}

// UserRoleEvent 用户角色变更事件
type UserRoleEvent struct {
	UserEvent
	roleIDs []int64
}

// NewUserRoleEvent 创建用户角色变更事件
func NewUserRoleEvent(tenantID, userID string, roleIDs []int64) *UserRoleEvent {
	return &UserRoleEvent{
		UserEvent: *NewUserEvent(tenantID, userID, UserRoleChanged),
		roleIDs:   roleIDs,
	}
}

func (e *UserRoleEvent) RoleIDs() []int64 {
	return e.roleIDs
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	UserEvent
}

// NewUserPasswordChangedEvent 创建用户密码修改事件
func NewUserPasswordChangedEvent(tenantID, userID string) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		UserEvent: *NewUserEvent(tenantID, userID, UserPasswordChanged),
	}
}
