package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

// 部门事件类型定义
const (
	DepartmentCreated = "department.created"
	DepartmentUpdated = "department.updated"
	DepartmentDeleted = "department.deleted"
	DepartmentMoved   = "department.moved"
	UserAssigned      = "department.user.assigned"
	UserRemoved       = "department.user.removed"
	UserTransferred   = "department.user.transferred"
)

// DepartmentEvent 部门事件基类
type DepartmentEvent struct {
	events.BaseEvent
	TenantID string `json:"tenant_id"`
	DeptID   string `json:"dept_id"`
}

// NewDepartmentEvent 创建部门事件
func NewDepartmentEvent(tenantID, deptID string, eventName string) *DepartmentEvent {
	return &DepartmentEvent{
		BaseEvent: events.NewBaseEvent(eventName),
		TenantID:  tenantID,
		DeptID:    deptID,
	}
}

// DepartmentMovedEvent 部门移动事件
type DepartmentMovedEvent struct {
	DepartmentEvent
	FromParentID string `json:"from_parent_id"`
	ToParentID   string `json:"to_parent_id"`
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
	UserIDs []string `json:"user_Ids"`
}

// NewUserAssignedEvent 创建用户分配事件
func NewUserAssignedEvent(tenantID, deptID string, userIDs []string) *UserAssignedEvent {
	return &UserAssignedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, UserAssigned),
		UserIDs:         userIDs,
	}
}

// UserRemovedEvent 用户移除事件
type UserRemovedEvent struct {
	DepartmentEvent
	UserIDs []string `json:"user_Ids"`
}

// NewUserRemovedEvent 创建用户移除事件
func NewUserRemovedEvent(tenantID, deptID string, userIDs []string) *UserRemovedEvent {
	return &UserRemovedEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, deptID, UserRemoved),
		UserIDs:         userIDs,
	}
}

// UserTransferredEvent 用户调动事件
type UserTransferredEvent struct {
	DepartmentEvent
	UserID     string `json:"user_id"`
	FromDeptID string `json:"from_dept_id"`
	ToDeptID   string `json:"to_dept_id"`
}

// NewUserTransferredEvent 创建用户调动事件
func NewUserTransferredEvent(tenantID, userID, fromDeptID, toDeptID string) *UserTransferredEvent {
	return &UserTransferredEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, toDeptID, UserTransferred),
		UserID:          userID,
		FromDeptID:      fromDeptID,
		ToDeptID:        toDeptID,
	}
}
