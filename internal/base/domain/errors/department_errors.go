package errors

import (
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 部门错误码定义
const (
	ReasonDepartmentNotFound      = "DEPARTMENT_NOT_FOUND"
	ReasonDepartmentExists        = "DEPARTMENT_EXISTS"
	ReasonDepartmentInvalid       = "DEPARTMENT_INVALID"
	ReasonDepartmentStatusInvalid = "DEPARTMENT_STATUS_INVALID"
	ReasonDepartmentDisabled      = "DEPARTMENT_DISABLED"
	ReasonParentDepartmentInvalid = "PARENT_DEPARTMENT_INVALID"
	ReasonHasChildDepartment      = "HAS_CHILD_DEPARTMENT"
	ReasonDepartmentInvalidOp     = "DEPARTMENT_INVALID_OPERATION"
	ReasonDepartmentCreateFailed  = "DEPARTMENT_CREATE_FAILED"
	ReasonDepartmentUpdateFailed  = "DEPARTMENT_UPDATE_FAILED"
	ReasonDepartmentDeleteFailed  = "DEPARTMENT_DELETE_FAILED"
	ReasonDepartmentQueryFailed   = "DEPARTMENT_QUERY_FAILED"
	ReasonUserAssignFailed        = "USER_ASSIGN_FAILED"
	ReasonUserRemoveFailed        = "USER_REMOVE_FAILED"
	ReasonUserTransferFailed      = "USER_TRANSFER_FAILED"
)

// DepartmentNotFound 部门不存在
func DepartmentNotFound(id string) herrors.Herr {
	return herrors.NewNotFoundHError(ReasonDepartmentNotFound,
		fmt.Errorf("department not found: %s", id))
}

// DepartmentExists 部门已存在
func DepartmentExists(code string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentExists,
		fmt.Errorf("department already exists: %s", code))
}

// DepartmentInvalidField 部门字段无效
func DepartmentInvalidField(field, reason string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentInvalid,
		fmt.Errorf("invalid department field %s: %s", field, reason))
}

// DepartmentStatusInvalid 部门状态无效
func DepartmentStatusInvalid(status int8) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentStatusInvalid,
		fmt.Errorf("invalid department status: %d", status))
}

// DepartmentDisabled 部门已禁用
func DepartmentDisabled(id string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentDisabled,
		fmt.Errorf("department is disabled: %s", id))
}

// ParentDepartmentNotFound 父部门不存在
func ParentDepartmentNotFound(id string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonParentDepartmentInvalid,
		fmt.Errorf("parent department not found: %s", id))
}

// ParentDepartmentDisabled 父部门已禁用
func ParentDepartmentDisabled(id string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonParentDepartmentInvalid,
		fmt.Errorf("parent department is disabled: %s", id))
}

// HasChildDepartment 存在子部门
func HasChildDepartment(id string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonHasChildDepartment,
		fmt.Errorf("department has children: %s", id))
}

// DepartmentInvalidOperation 部门操作无效
func DepartmentInvalidOperation(reason string) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentInvalidOp,
		fmt.Errorf("invalid department operation: %s", reason))
}

// DepartmentCreateFailed 部门创建失败
func DepartmentCreateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentCreateFailed,
		fmt.Errorf("failed to create department: %v", err))
}

// DepartmentUpdateFailed 部门更新失败
func DepartmentUpdateFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentUpdateFailed,
		fmt.Errorf("failed to update department: %v", err))
}

// DepartmentDeleteFailed 部门删除失败
func DepartmentDeleteFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentDeleteFailed,
		fmt.Errorf("failed to delete department: %v", err))
}

// DepartmentQueryFailed 部门查询失败
func DepartmentQueryFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonDepartmentQueryFailed,
		fmt.Errorf("failed to query department: %v", err))
}

// UserAssignFailed 用户分配失败
func UserAssignFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonUserAssignFailed,
		fmt.Errorf("failed to assign users: %v", err))
}

// UserRemoveFailed 用户移除失败
func UserRemoveFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonUserRemoveFailed,
		fmt.Errorf("failed to remove users: %v", err))
}

// UserTransferFailed 用户调动失败
func UserTransferFailed(err error) herrors.Herr {
	return herrors.NewBadRequestHError(ReasonUserTransferFailed,
		fmt.Errorf("failed to transfer user: %v", err))
}
