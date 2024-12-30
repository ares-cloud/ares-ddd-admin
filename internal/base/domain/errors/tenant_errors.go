package errors

import (
	"fmt"
	"net/http"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 租户相关错误码定义
const (
	ReasonTenantNotFound      = "TENANT_NOT_FOUND"
	ReasonTenantCodeExists    = "TENANT_CODE_EXISTS"
	ReasonTenantIsDefault     = "TENANT_IS_DEFAULT"
	ReasonTenantExpired       = "TENANT_EXPIRED"
	ReasonTenantInvalid       = "TENANT_INVALID"
	ReasonTenantStatusInvalid = "TENANT_STATUS_INVALID"
	ReasonTenantDisabled      = "TENANT_DISABLED"
	ReasonTenantAdminInvalid  = "TENANT_ADMIN_INVALID"
	ReasonTenantDomainInvalid = "TENANT_DOMAIN_INVALID"
	ReasonTenantQuotaExceeded = "TENANT_QUOTA_EXCEEDED"
)

// TenantNotFound 租户不存在
func TenantNotFound(id string) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonTenantNotFound,
		fmt.Sprintf("tenant not found: %s", id))
}

// TenantCodeExists 租户编码已存在
func TenantCodeExists(code string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantCodeExists,
		fmt.Sprintf("tenant code already exists: %s", code))
}

// TenantIsDefault 默认租户不可删除
func TenantIsDefault() herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantIsDefault,
		"cannot delete default tenant")
}

// TenantExpired 租户已过期
func TenantExpired() herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonTenantExpired,
		"tenant has expired")
}

// TenantInvalidField 字段验证错误
func TenantInvalidField(field, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantInvalid,
		fmt.Sprintf("invalid tenant %s: %s", field, reason))
}

// TenantStatusInvalid 状态无效
func TenantStatusInvalid(status int8) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantStatusInvalid,
		fmt.Sprintf("invalid tenant status: %d, must be 1(enabled) or 2(disabled)", status))
}

// TenantDisabled 租户已禁用
func TenantDisabled(reason string) herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonTenantDisabled,
		fmt.Sprintf("tenant is disabled: %s", reason))
}

// TenantPermissionDenied 租户权限不足
func TenantPermissionDenied(permission string) herrors.Herr {
	return herrors.New(http.StatusForbidden, "TENANT_PERMISSION_DENIED",
		fmt.Sprintf("tenant does not have permission: %s", permission))
}

// TenantAdminInvalid 管理员信息无效
func TenantAdminInvalid(reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantAdminInvalid,
		fmt.Sprintf("invalid tenant admin: %s", reason))
}

// TenantDomainInvalid 域名无效
func TenantDomainInvalid(domain, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonTenantDomainInvalid,
		fmt.Sprintf("invalid tenant domain %s: %s", domain, reason))
}

// TenantQuotaExceeded 超出配额限制
func TenantQuotaExceeded(resource string, limit int64) herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonTenantQuotaExceeded,
		fmt.Sprintf("tenant quota exceeded for %s, limit: %d", resource, limit))
}
