
-- 添加左右值和层级索引
CREATE INDEX idx_dept_tree ON departments (left_value, right_value, level);

-- 添加部门状态索引
CREATE INDEX idx_dept_status ON departments (status);

-- 添加复合索引
CREATE INDEX idx_dept_parent_status ON departments (parent_id, status);