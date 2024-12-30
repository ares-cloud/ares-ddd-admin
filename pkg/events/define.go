package events

// BaseTenantEvent 租户基础事件
type BaseTenantEvent struct {
	BaseEvent
	version       string // 版本号
	aggregateID   string // 聚合根ID
	aggregateType string // 聚合根类型
	tenantID      string // 租户ID
}

// NewBaseTenantEvent 创建租户基础事件
func NewBaseTenantEvent(
	eventName string,
	version string,
	aggregateID string,
	aggregateType string,
	tenantID string,
) BaseTenantEvent {
	return BaseTenantEvent{
		BaseEvent:     NewBaseEvent(eventName),
		version:       version,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		tenantID:      tenantID,
	}
}

// Version 获取版本号
func (e *BaseTenantEvent) Version() string {
	return e.version
}

// AggregateID 获取聚合根ID
func (e *BaseTenantEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType 获取聚合根类型
func (e *BaseTenantEvent) AggregateType() string {
	return e.aggregateType
}

// TenantID 获取租户ID
func (e *BaseTenantEvent) TenantID() string {
	return e.tenantID
}