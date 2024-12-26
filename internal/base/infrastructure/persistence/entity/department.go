package entity

// Department 部门数据库实体
type Department struct {
	ID          string `gorm:"column:id;primary_key"`
	TenantID    string `gorm:"column:tenant_id"`
	ParentID    string `gorm:"column:parent_id"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	Sort        int    `gorm:"column:sort"`
	Leader      string `gorm:"column:leader"`
	Phone       string `gorm:"column:phone"`
	Email       string `gorm:"column:email"`
	Status      int8   `gorm:"column:status"`
	Description string `gorm:"column:description"`
	CreatedAt   int64  `gorm:"column:created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at"`
}

// TableName 表名
func (Department) TableName() string {
	return "sys_department"
}
func (Department) GetPrimaryKey() string {
	return "id"
}

// UserDepartment 用户部门关系表
type UserDepartment struct {
	ID     string `gorm:"column:id;primary_key;comment:ID"`
	UserID string `gorm:"column:user_id;comment:用户ID"`
	DeptID string `gorm:"column:dept_id;comment:部门ID"`
}

// TableName 表名
func (UserDepartment) TableName() string {
	return "sys_user_dept"
}
