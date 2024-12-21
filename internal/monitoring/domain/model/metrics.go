package model

import "time"

// SystemMetrics 系统指标
type SystemMetrics struct {
	ID          int64     // ID
	TenantID    string    // 租户ID
	CPUUsage    float64   // CPU使用率
	MemoryUsage float64   // 内存使用率
	DiskUsage   float64   // 磁盘使用率
	CreatedAt   time.Time // 创建时间
}

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	ID           int64     // ID
	TenantID     string    // 租户ID
	Connections  int64     // 连接数
	SlowQueries  int64     // 慢查询数
	QPS          float64   // 每秒查询数
	ResponseTime float64   // 响应时间
	CreatedAt    time.Time // 创建时间
}

// RedisMetrics Redis指标
type RedisMetrics struct {
	ID          int64     // ID
	TenantID    string    // 租户ID
	Memory      int64     // 内存使用
	Connections int64     // 连接数
	Commands    int64     // 命令执行数
	KeyCount    int64     // 键数量
	HitRate     float64   // 命中率
	CreatedAt   time.Time // 创建时间
}

// APIMetrics API指标
type APIMetrics struct {
	ID           int64     // ID
	TenantID     string    // 租户ID
	Path         string    // 接口路径
	Method       string    // 请求方法
	Count        int64     // 请求次数
	ErrorCount   int64     // 错误次数
	ResponseTime float64   // 平均响应时间
	CreatedAt    time.Time // 创建时间
}
