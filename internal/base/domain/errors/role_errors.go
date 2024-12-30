package errors

// 角色相关错误
var (
	ErrRoleNotFound   = New("role not found")
	ErrRoleCodeExists = New("role code already exists")
	ErrRoleInUse      = New("role is in use")
)
