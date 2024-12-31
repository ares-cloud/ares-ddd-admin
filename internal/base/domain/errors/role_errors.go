package errors

import (
	"fmt"
	"net/http"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 角色相关错误码定义
const (
	ReasonRoleNotFound      = "ROLE_NOT_FOUND"      // 角色不存在
	ReasonRoleExists        = "ROLE_EXISTS"         // 角色已存在
	ReasonRoleInUse         = "ROLE_IN_USE"         // 角色正在使用中
	ReasonRoleInvalid       = "ROLE_INVALID"        // 角色无效
	ReasonRoleStatusInvalid = "ROLE_STATUS_INVALID" // 角色状态无效
	ReasonRoleDisabled      = "ROLE_DISABLED"       // 角色已禁用
	ReasonRolePermDenied    = "ROLE_PERM_DENIED"    // 角色权限不足
	ReasonPermNotFound      = "PERM_NOT_FOUND"      // 权限不存在
)

// RoleNotFound 角色不存在
func RoleNotFound(id int64) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonRoleNotFound,
		fmt.Sprintf("role not found: %d", id))
}

// RoleExists 角色已存在
func RoleExists(code string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonRoleExists,
		fmt.Sprintf("role already exists: %s", code))
}

// RoleInUse 角色正在使用中
func RoleInUse(id int64) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonRoleInUse,
		fmt.Sprintf("role is in use: %d", id))
}

// RoleInvalidField 角色字段无效
func RoleInvalidField(field, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonRoleInvalid,
		fmt.Sprintf("invalid role %s: %s", field, reason))
}

// RoleStatusInvalid 角色状态无效
func RoleStatusInvalid(status int8) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonRoleStatusInvalid,
		fmt.Sprintf("invalid role status: %d, must be 1(enabled) or 2(disabled)", status))
}

// RoleDisabled 角色已禁用
func RoleDisabled(id int64, reason string) herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonRoleDisabled,
		fmt.Sprintf("role is disabled: %d, reason: %s", id, reason))
}

// RolePermissionDenied 角色权限不足
func RolePermissionDenied(roleID int64, permID int64) herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonRolePermDenied,
		fmt.Sprintf("role %d does not have permission %d", roleID, permID))
}

// PermissionNotFound 权限不存在
func PermissionNotFound(id int64) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonPermNotFound,
		fmt.Sprintf("permission not found: %d", id))
}
