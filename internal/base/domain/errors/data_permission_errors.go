package errors

import (
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 数据权限错误码定义
const (
	ReasonDataPermissionNotFound     = "DATA_PERMISSION_NOT_FOUND"
	ReasonDataPermissionExists       = "DATA_PERMISSION_EXISTS"
	ReasonDataPermissionInvalid      = "DATA_PERMISSION_INVALID"
	ReasonDataPermissionCreateFailed = "DATA_PERMISSION_CREATE_FAILED"
	ReasonDataPermissionUpdateFailed = "DATA_PERMISSION_UPDATE_FAILED"
	ReasonDataPermissionDeleteFailed = "DATA_PERMISSION_DELETE_FAILED"
	ReasonDataPermissionQueryFailed  = "DATA_PERMISSION_QUERY_FAILED"
)

// DataPermissionNotFound 数据权限不存在
func DataPermissionNotFound(roleID int64) herrors.Herr {
	return herrors.NewNotFoundHError(ReasonDataPermissionNotFound,
		fmt.Errorf("data permission not found for role: %d", roleID))
}

// DataPermissionExists 数据权限已存在
func DataPermissionExists(roleID int64) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionExists,
		fmt.Errorf("data permission already exists for role: %d", roleID))
}

// DataPermissionInvalidField 数据权限字段无效
func DataPermissionInvalidField(field, reason string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionInvalid,
		fmt.Errorf("invalid data permission field %s: %s", field, reason))
}

// DataPermissionCreateFailed 数据权限创建失败
func DataPermissionCreateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionCreateFailed,
		fmt.Errorf("failed to create data permission: %v", err))
}

// DataPermissionUpdateFailed 数据权限更新失败
func DataPermissionUpdateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionUpdateFailed,
		fmt.Errorf("failed to update data permission: %v", err))
}

// DataPermissionDeleteFailed 数据权限删除失败
func DataPermissionDeleteFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionDeleteFailed,
		fmt.Errorf("failed to delete data permission: %v", err))
}

// DataPermissionQueryFailed 数据权限查询失败
func DataPermissionQueryFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDataPermissionQueryFailed,
		fmt.Errorf("failed to query data permission: %v", err))
}
