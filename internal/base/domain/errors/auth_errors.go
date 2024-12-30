package errors

// 认证相关错误
var (
	ErrInvalidCredentials = New("invalid username or password")
	ErrUserDisabled       = New("user is disabled")
	ErrPasswordMismatch   = New("password mismatch")
	ErrTokenInvalid       = New("invalid token")
	ErrTokenExpired       = New("token expired")
)
