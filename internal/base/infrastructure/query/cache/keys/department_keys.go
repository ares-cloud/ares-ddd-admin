package keys

import "fmt"

const (
	deptPrefix = "dept:"
)

// DepartmentKey 部门缓存key
func DepartmentKey(tenantID string, deptID string) string {
	return fmt.Sprintf("%s%s:%s", deptPrefix, tenantID, deptID)
}

// DepartmentTreeKey 部门树缓存key
func DepartmentTreeKey(tenantID string, parentID string) string {
	if parentID == "" {
		return fmt.Sprintf("%s%s:tree", deptPrefix, tenantID)
	}
	return fmt.Sprintf("%s%s:tree:%s", deptPrefix, tenantID, parentID)
}

// UserDepartmentsKey 用户部门缓存key
func UserDepartmentsKey(tenantID string, userID string) string {
	return fmt.Sprintf("user:%s:depts:%s", tenantID, userID)
}

// DepartmentKeys 生成部门相关的所有缓存key
func DepartmentKeys(tenantID string, deptID string) []string {
	return []string{
		DepartmentKey(tenantID, deptID),
		DepartmentTreeKey(tenantID, ""),
	}
}

// DepartmentChildrenKey 部门子节点列表缓存key
func DepartmentChildrenKey(tenantID string, parentID string) string {
	return fmt.Sprintf("%s%s:children:%s", deptPrefix, tenantID, parentID)
}

// DepartmentUsersKey 部门用户列表缓存key
func DepartmentUsersKey(tenantID string, deptID string) string {
	return fmt.Sprintf("%s%s:users:%s", deptPrefix, tenantID, deptID)
}
