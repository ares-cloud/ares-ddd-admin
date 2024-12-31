package herrors

import (
	"net/http"
)

// NewBadRequestHError 创建400错误
func NewBadRequestHError(reason string, err error) Herr {
	return &HError{Code: http.StatusBadRequest, Reason: reason, DefMessage: reason, BusinessError: err}
}

// NewUnauthorizedHError 创建401错误
func NewUnauthorizedHError(reason string, err error) Herr {
	return &HError{Code: http.StatusUnauthorized, Reason: reason, DefMessage: reason, BusinessError: err}
}

// NewForbiddenHError 创建403错误
func NewForbiddenHError(reason string, err error) Herr {
	return &HError{Code: http.StatusForbidden, Reason: reason, DefMessage: reason, BusinessError: err}
}

// NewNotFoundHError 创建404错误
func NewNotFoundHError(reason string, err error) Herr {
	return &HError{Code: http.StatusNotFound, Reason: reason, DefMessage: reason, BusinessError: err}
}

// NewConflictHError 创建409错误
func NewConflictHError(reason string, err error) Herr {
	return &HError{Code: http.StatusConflict, Reason: reason, DefMessage: reason, BusinessError: err}
}
