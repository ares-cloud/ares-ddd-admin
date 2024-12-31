package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	// 部门聚合根名称
	DepartmentAggregate = "department"
	// 版本号
)

// 部门事件类型定义
const (
	DepartmentCreated = "department.created"
	DepartmentUpdated = "department.updated"
	DepartmentDeleted = "department.deleted"
	DepartmentMoved   = "department.moved"
	UserDeptChanged   = "department.user.changed" // 用户部门变更事件
)

// DepartmentEvent 部门事件
type DepartmentEvent struct {
	events.BaseTenantEvent
}

// NewDepartmentEvent 创建部门事件
func NewDepartmentEvent(tenantID, deptID string, eventName string) *DepartmentEvent {
	return &DepartmentEvent{
		//BaseTenantEvent: events.NewBaseTenantEvent(
		//	eventName,
		//	Version,
		//	deptID,
		//	DepartmentAggregate,
		//	tenantID,
		//),
	}
}

func (e *DepartmentEvent) DepartmentID() string {
	return e.AggregateID()
}

// UserDeptEvent 用户部门变更事件
type UserDeptEvent struct {
	DepartmentEvent
	userID    string // 用户ID
	oldDeptID string // 原部门ID
	newDeptID string // 新部门ID
}

// NewUserDeptEvent 创建用户部门变更事件
func NewUserDeptEvent(tenantID string, userID string, oldDeptID string, newDeptID string) *UserDeptEvent {
	return &UserDeptEvent{
		DepartmentEvent: *NewDepartmentEvent(tenantID, newDeptID, UserDeptChanged),
		userID:          userID,
		oldDeptID:       oldDeptID,
		newDeptID:       newDeptID,
	}
}

// UserID 获取用户ID
func (e *UserDeptEvent) UserID() string {
	return e.userID
}

// OldDeptID 获取原部门ID
func (e *UserDeptEvent) OldDeptID() string {
	return e.oldDeptID
}

// NewDeptID 获取新部门ID
func (e *UserDeptEvent) NewDeptID() string {
	return e.newDeptID
}
