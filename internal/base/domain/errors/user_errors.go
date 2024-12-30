package errors

// 用户领域错误定义
var (

	// 用户操作相关错误
	ErrUserNotFound   = New("user not found")
	ErrUsernameExists = New("username already exists")
	ErrInvalidStatus  = New("invalid user status")
)
