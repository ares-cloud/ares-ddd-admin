package queries

// GetMetricsQuery 获取监控指标查询基类
type GetMetricsQuery struct{}

// GetSystemMetricsQuery 获取系统指标查询
type GetSystemMetricsQuery struct {
	GetMetricsQuery
}

// GetRuntimeMetricsQuery 获取运行时指标查询
type GetRuntimeMetricsQuery struct {
	GetMetricsQuery
}

// GetDatabaseMetricsQuery 获取数据库指标查询
type GetDatabaseMetricsQuery struct {
	GetMetricsQuery
}

// GetRedisMetricsQuery 获取Redis指标查询
type GetRedisMetricsQuery struct {
	GetMetricsQuery
}

// GetAPIMetricsQuery 获取API指标查询
type GetAPIMetricsQuery struct {
	GetMetricsQuery
	Path   string `query:"path"`   // 接口路径
	Method string `query:"method"` // 请求方法
}
