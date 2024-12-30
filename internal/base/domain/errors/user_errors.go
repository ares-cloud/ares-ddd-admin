package errors

import (
	"fmt"
	"net/http"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

const (
	ReasonUserNotFound      = "USER_NOT_FOUND"
	ReasonUserExists        = "USER_EXISTS"
	ReasonUserInvalid       = "USER_INVALID"
	ReasonUserStatusInvalid = "USER_STATUS_INVALID"
	ReasonUserDisabled      = "USER_DISABLED"
)

// UserNotFound 用户不存在
func UserNotFound(id string) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonUserNotFound,
		fmt.Sprintf("user not found: %s", id))
}

// UserExists 用户已存在
func UserExists(username string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonUserExists,
		fmt.Sprintf("user already exists: %s", username))
}

// UserInvalidField 字段验证错误
func UserInvalidField(field, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonUserInvalid,
		fmt.Sprintf("invalid user %s: %s", field, reason))
}

// UserStatusInvalid 状态无效
func UserStatusInvalid(status int8) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonUserStatusInvalid,
		fmt.Sprintf("invalid user status: %d, must be 1(enabled) or 2(disabled)", status))
}

// UserDisabled 用户已禁用
func UserDisabled(reason string) herrors.Herr {
	return herrors.New(http.StatusForbidden, ReasonUserDisabled,
		fmt.Sprintf("user is disabled: %s", reason))
}

// UserInvalidDepartment 用户部门错误
func UserInvalidDepartment(userID string, deptID string, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonUserInvalid,
		fmt.Sprintf("invalid user[%s] department[%s]: %s", userID, deptID, reason))
}

// UserDepartmentNotFound 用户不属于该部门
func UserDepartmentNotFound(userID string, deptID string) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonUserNotFound,
		fmt.Sprintf("user[%s] does not belong to department[%s]", userID, deptID))
}
