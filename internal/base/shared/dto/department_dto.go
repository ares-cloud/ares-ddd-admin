package dto

// DepartmentDto 部门DTO
type DepartmentDto struct {
	ID          string           `json:"id"`          // 部门ID
	ParentID    string           `json:"parentId"`    // 父部门ID
	Name        string           `json:"name"`        // 部门名称
	Code        string           `json:"code"`        // 部门编码
	Sort        int32            `json:"sort"`        // 排序
	Leader      string           `json:"leader"`      // 负责人
	Phone       string           `json:"phone"`       // 联系电话
	Email       string           `json:"email"`       // 邮箱
	Status      int8             `json:"status"`      // 部门状态
	Description string           `json:"description"` // 描述
	Children    []*DepartmentDto `json:"children"`    // 子部门
}

// DepartmentTreeDto 部门树DTO
type DepartmentTreeDto struct {
	ID       string               `json:"id"`       // 部门ID
	ParentID string               `json:"parentId"` // 父部门ID
	Name     string               `json:"name"`     // 部门名称
	Children []*DepartmentTreeDto `json:"children"` // 子部门
}
