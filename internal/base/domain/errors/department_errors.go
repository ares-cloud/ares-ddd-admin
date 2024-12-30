package errors

// 部门相关错误
var (
	ErrDepartmentNotFound = New("department not found")
	ErrDepartmentExists   = New("department already exists")
	ErrInvalidParentDept  = New("invalid parent department")
	ErrHasChildDept       = New("department has child departments")
	ErrHasUsers           = New("department has users")
)
