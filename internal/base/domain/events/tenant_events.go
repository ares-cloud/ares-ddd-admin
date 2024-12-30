package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

const (
	// 租户聚合根名称
	TenantAggregate = "tenant"
)

// 租户事件类型定义
const (
	TenantCreated       = "tenant.created"
	TenantUpdated       = "tenant.updated"
	TenantDeleted       = "tenant.deleted"
	TenantStatusChanged = "tenant.status.changed"
)

// TenantEvent 租户事件
type TenantEvent struct {
	events.BaseTenantEvent
}

// NewTenantEvent 创建租户事件
func NewTenantEvent(tenantID string, eventName string) *TenantEvent {
	return &TenantEvent{
		BaseTenantEvent: events.NewBaseTenantEvent(
			eventName,
			Version,
			tenantID,
			TenantAggregate,
			tenantID,
		),
	}
}
