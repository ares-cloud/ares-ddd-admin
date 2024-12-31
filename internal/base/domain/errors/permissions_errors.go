package errors

import (
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 权限错误码定义
const (
	ReasonPermissionNotFound     = "PERMISSION_NOT_FOUND"
	ReasonPermissionExists       = "PERMISSION_EXISTS"
	ReasonPermissionInvalid      = "PERMISSION_INVALID"
	ReasonPermissionDisabled     = "PERMISSION_DISABLED"
	ReasonHasChildPermission     = "HAS_CHILD_PERMISSION"
	ReasonPermissionCreateFailed = "PERMISSION_CREATE_FAILED"
	ReasonPermissionUpdateFailed = "PERMISSION_UPDATE_FAILED"
	ReasonPermissionDeleteFailed = "PERMISSION_DELETE_FAILED"
	ReasonPermissionQueryFailed  = "PERMISSION_QUERY_FAILED"
)

// PermissionExists 权限已存在
func PermissionExists(code string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionExists,
		fmt.Errorf("permission already exists: %s", code))
}

// PermissionInvalidField 权限字段无效
func PermissionInvalidField(field, reason string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionInvalid,
		fmt.Errorf("invalid permission field %s: %s", field, reason))
}

// HasChildPermission 存在子权限
func HasChildPermission(id int64) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonHasChildPermission,
		fmt.Errorf("permission has children: %d", id))
}

// PermissionCreateFailed 权限创建失败
func PermissionCreateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionCreateFailed,
		fmt.Errorf("failed to create permission: %v", err))
}

// PermissionUpdateFailed 权限更新失败
func PermissionUpdateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionUpdateFailed,
		fmt.Errorf("failed to update permission: %v", err))
}

// PermissionDeleteFailed 权限删除失败
func PermissionDeleteFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionDeleteFailed,
		fmt.Errorf("failed to delete permission: %v", err))
}

// PermissionQueryFailed 权限查询失败
func PermissionQueryFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonPermissionQueryFailed,
		fmt.Errorf("failed to query permission: %v", err))
}
