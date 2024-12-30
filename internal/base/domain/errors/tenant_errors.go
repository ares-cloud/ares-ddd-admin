package errors

// 租户相关错误
var (
	ErrTenantNotFound     = New("tenant not found")
	ErrTenantCodeExists   = New("tenant code already exists")
	ErrTenantDomainExists = New("tenant domain already exists")
	ErrDefaultTenant      = New("cannot modify default tenant")
)
