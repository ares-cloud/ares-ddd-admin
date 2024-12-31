package entity

// UserDepartment 用户部门关系表
type UserDepartment struct {
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一ID"` // 唯一ID
	UserID string `gorm:"column:user_id;comment:用户ID"`
	DeptID string `gorm:"column:dept_id;comment:部门ID"`
}

// TableName 表名
func (UserDepartment) TableName() string {
	return "sys_user_dept"
}
