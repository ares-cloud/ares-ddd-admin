package keys

import "fmt"

const (
	deptPrefix = "dept:"
)

// DepartmentKey 部门缓存key
func DepartmentKey(deptID string) string {
	return fmt.Sprintf("%s%s", deptPrefix, deptID)
}

// DepartmentTreeKey 部门树缓存key
func DepartmentTreeKey(parentID string) string {
	if parentID == "" {
		return fmt.Sprintf("%stree", deptPrefix)
	}
	return fmt.Sprintf("%stree:%s", deptPrefix, parentID)
}

// UserDepartmentsKey 用户部门缓存key
func UserDepartmentsKey(userID string) string {
	return fmt.Sprintf("user:depts:%s", userID)
}

// DepartmentKeys 生成部门相关的所有缓存key
func DepartmentKeys(deptID string) []string {
	return []string{
		DepartmentKey(deptID),
		DepartmentTreeKey(""),
	}
}
