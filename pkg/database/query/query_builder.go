package query

import (
	"fmt"
	"strings"
)

// 数据权限范围枚举
const (
	DataScopeAll      int8 = 1 // 全部数据
	DataScopeTenant   int8 = 2 // 本租户数据
	DataScopeDept     int8 = 3 // 本部门数据
	DataScopeDeptTree int8 = 4 // 本部门及以下数据
	DataScopeSelf     int8 = 5 // 仅本人数据
	DataScopeCustom   int8 = 6 // 自定义数据
)

// Operator 查询操作符
type Operator string

const (
	Eq        Operator = "="
	Neq       Operator = "!="
	Gt        Operator = ">"
	Gte       Operator = ">="
	Lt        Operator = "<"
	Lte       Operator = "<="
	Like      Operator = "LIKE"
	In        Operator = "IN"
	NotIn     Operator = "NOT IN"
	IsNull    Operator = "IS NULL"
	IsNotNull Operator = "IS NOT NULL"
)

// Condition 查询条件
type Condition struct {
	Field    string      // 字段名
	Operator Operator    // 操作符
	Value    interface{} // 值
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	conditions []Condition
	orderBy    []string
	page       *Page
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		conditions: make([]Condition, 0),
		orderBy:    make([]string, 0),
	}
}

// Where 添加查询条件
func (qb *QueryBuilder) Where(field string, operator Operator, value interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return qb
}

// OrderBy 添加排序
func (qb *QueryBuilder) OrderBy(field string, asc bool) *QueryBuilder {
	direction := "DESC"
	if asc {
		direction = "ASC"
	}
	qb.orderBy = append(qb.orderBy, fmt.Sprintf("%s %s", field, direction))
	return qb
}

// WithPage 设置分页
func (qb *QueryBuilder) WithPage(page *Page) *QueryBuilder {
	qb.page = page
	return qb
}

// BuildWhere 构建WHERE子句
func (qb *QueryBuilder) BuildWhere() (string, []interface{}) {
	if len(qb.conditions) == 0 {
		return "", nil
	}

	var (
		where  strings.Builder
		values []interface{}
	)

	where.WriteString("WHERE ")
	for i, cond := range qb.conditions {
		if i > 0 {
			where.WriteString(" AND ")
		}

		switch cond.Operator {
		case IsNull, IsNotNull:
			where.WriteString(fmt.Sprintf("%s %s", cond.Field, cond.Operator))
		case In, NotIn:
			where.WriteString(fmt.Sprintf("%s %s (?)", cond.Field, cond.Operator))
			values = append(values, cond.Value)
		default:
			where.WriteString(fmt.Sprintf("%s %s ?", cond.Field, cond.Operator))
			values = append(values, cond.Value)
		}
	}

	return where.String(), values
}

// BuildOrderBy 构建ORDER BY子句
func (qb *QueryBuilder) BuildOrderBy() string {
	if len(qb.orderBy) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(qb.orderBy, ", ")
}

// BuildLimit 构建LIMIT子句
func (qb *QueryBuilder) BuildLimit() (string, []int) {
	if qb.page == nil {
		return "", nil
	}
	qb.page.Fix()
	return "LIMIT ?, ?", []int{qb.page.Offset(), qb.page.Limit()}
}
