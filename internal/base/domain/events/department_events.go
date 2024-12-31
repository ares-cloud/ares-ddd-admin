package events

import (
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 部门事件类型定义
const (
	DepartmentCreated  = "department.created"
	DepartmentUpdated  = "department.updated"
	DepartmentDeleted  = "department.deleted"
	DepartmentDisabled = "department.disabled"
	DepartmentEnabled  = "department.enabled"
	DepartmentMoved    = "department.moved"
	UserAssigned       = "department.user.assigned"
	UserRemoved        = "department.user.removed"
	UserTransferred    = "department.user.transferred"
)

// DepartmentEvent 部门事件基类
type DepartmentEvent struct {
	events.BaseEvent
	tenantID string
	deptID   string
}

// NewDepartmentEvent 创建部门事件
func NewDepartmentEvent(tenantID, deptID string, eventName string) *DepartmentEvent {
	return &DepartmentEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		tenantID:  tenantID,
		deptID:    deptID,
	}
}

func (e *DepartmentEvent) DepartmentID() string {
	return e.deptID
}

func (e *DepartmentEvent) TenantID() string {
	return e.tenantID
}

// DepartmentCreatedEvent 部门创建事件
type DepartmentCreatedEvent struct {
	DepartmentEvent
}

// NewDepartmentCreatedEvent 创建部门创建事件
func NewDepartmentCreatedEvent(tenantID, deptID string) *DepartmentCreatedEvent {
	return &DepartmentCreatedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, DepartmentCreated),
	}
}

// DepartmentUpdatedEvent 部门更新事件
type DepartmentUpdatedEvent struct {
	DepartmentEvent
}

// NewDepartmentUpdatedEvent 创建部门更新事件
func NewDepartmentUpdatedEvent(tenantID, deptID string) *DepartmentUpdatedEvent {
	return &DepartmentUpdatedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, DepartmentUpdated),
	}
}

// DepartmentDeletedEvent 部门删除事件
type DepartmentDeletedEvent struct {
	DepartmentEvent
}

// NewDepartmentDeletedEvent 创建部门删除事件
func NewDepartmentDeletedEvent(tenantID, deptID string) *DepartmentDeletedEvent {
	return &DepartmentDeletedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, DepartmentDeleted),
	}
}

// DepartmentMovedEvent 部门移动事件
type DepartmentMovedEvent struct {
	DepartmentEvent
	FromParentID string
	ToParentID   string
}

// NewDepartmentMovedEvent 创建部门移动事件
func NewDepartmentMovedEvent(tenantID string, deptID string, fromParentID string, toParentID string) *DepartmentMovedEvent {
	return &DepartmentMovedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, DepartmentMoved),
		FromParentID:    fromParentID,
		ToParentID:      toParentID,
	}
}

// UserAssignedEvent 用户分配事件
type UserAssignedEvent struct {
	DepartmentEvent
	userIDs []string
}

// NewUserAssignedEvent 创建用户分配事件
func NewUserAssignedEvent(tenantID, deptID string, userIDs []string) *UserAssignedEvent {
	return &UserAssignedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, UserAssigned),
		userIDs:         userIDs,
	}
}

func (e *UserAssignedEvent) UserIDs() []string {
	return e.userIDs
}

// UserRemovedEvent 用户移除事件
type UserRemovedEvent struct {
	DepartmentEvent
	userIDs []string
}

// NewUserRemovedEvent 创建用户移除事件
func NewUserRemovedEvent(tenantID, deptID string, userIDs []string) *UserRemovedEvent {
	return &UserRemovedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, UserRemoved),
		userIDs:         userIDs,
	}
}

func (e *UserRemovedEvent) UserIDs() []string {
	return e.userIDs
}

// UserTransferredEvent 用户调动事件
type UserTransferredEvent struct {
	DepartmentEvent
	userID     string
	fromDeptID string
	toDeptID   string
	timestamp  int64
}

// NewUserTransferredEvent 创建用户调动事件
func NewUserTransferredEvent(tenantID, userID, fromDeptID, toDeptID string) *UserTransferredEvent {
	return &UserTransferredEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, toDeptID, UserTransferred),
		userID:          userID,
		fromDeptID:      fromDeptID,
		toDeptID:        toDeptID,
		timestamp:       time.Now().Unix(),
	}
}

func (e *UserTransferredEvent) GetTopic() string {
	return "department.user.transfer"
}

// DepartmentAdminSetEvent 部门管理员设置事件
type DepartmentAdminSetEvent struct {
	DepartmentEvent
	AdminID string
}

// NewDepartmentAdminSetEvent 创建部门管理员设置事件
func NewDepartmentAdminSetEvent(tenantID string, deptID string, adminID string) *DepartmentAdminSetEvent {
	return &DepartmentAdminSetEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, UserTransferred),
		AdminID:         adminID,
	}
}
