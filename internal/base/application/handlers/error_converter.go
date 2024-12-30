package handlers

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

// 错误转换函数
func convertDomainError(err error) herrors.Herr {
	if err == nil {
		return nil
	}

	switch err {
	// 认证错误
	case errors.ErrInvalidCredentials:
		return herrors.NewBadReqError("invalid username or password")
	case errors.ErrUserDisabled:
		return herrors.NewBadReqError("user is disabled")
	case errors.ErrPasswordMismatch:
		return herrors.NewBadReqError("password mismatch")
	case errors.ErrTokenInvalid, errors.ErrTokenExpired:
		return herrors.BaseTokenVerifyFail

	// 用户错误
	case errors.ErrUserNotFound:
		return herrors.ErrRecordNotFount
	case errors.ErrUsernameExists:
		return herrors.DataIsExist

	// 租户错误
	case errors.ErrTenantNotFound:
		return herrors.ErrRecordNotFount
	case errors.ErrTenantCodeExists, errors.ErrTenantDomainExists:
		return herrors.DataIsExist

	// 其他错误
	default:
		return herrors.NewServerHError(err)
	}
}
