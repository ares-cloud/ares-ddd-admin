package dto

import "time"

// SystemMetricsDto 系统指标DTO
type SystemMetricsDto struct {
	CPUUsage    float64   `json:"cpu_usage"`    // CPU使用率
	MemoryUsage float64   `json:"memory_usage"` // 内存使用率
	DiskUsage   float64   `json:"disk_usage"`   // 磁盘使用率
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
}

// RuntimeMetricsDto 运行时指标DTO
type RuntimeMetricsDto struct {
	Goroutines    int     `json:"goroutines"`      // Goroutine数量
	HeapAlloc     uint64  `json:"heap_alloc"`      // 已分配堆内存
	HeapSys       uint64  `json:"heap_sys"`        // 系统预留堆内存
	HeapObjects   uint64  `json:"heap_objects"`    // 堆对象数量
	StackInUse    uint64  `json:"stack_in_use"`    // 正在使用的栈内存
	StackSys      uint64  `json:"stack_sys"`       // 系统预留栈内存
	MSpanInUse    uint64  `json:"mspan_in_use"`    // 正在使用的MSpan内存
	MSpanSys      uint64  `json:"mspan_sys"`       // 系统预留MSpan内存
	MCacheInUse   uint64  `json:"mcache_in_use"`   // 正在使用的MCache内存
	MCacheSys     uint64  `json:"mcache_sys"`      // 系统预留MCache内存
	GCPauseNs     uint64  `json:"gc_pause_ns"`     // 最后一次GC暂停时间(纳秒)
	LastGC        uint64  `json:"last_gc"`         // 上次GC时间
	NumGC         uint32  `json:"num_gc"`          // GC次数
	GCCPUFraction float64 `json:"gc_cpu_fraction"` // GC占用CPU时间比例
}

// DatabaseMetricsDto 数据库指标DTO
type DatabaseMetricsDto struct {
	Connections  int64     `json:"connections"`   // 连接数
	SlowQueries  int64     `json:"slow_queries"`  // 慢查询数
	QPS          float64   `json:"qps"`           // 每秒查询数
	ResponseTime float64   `json:"response_time"` // 响应时间
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
}

// RedisMetricsDto Redis指标DTO
type RedisMetricsDto struct {
	Memory      int64     `json:"memory"`      // 内存使用
	Connections int64     `json:"connections"` // 连接数
	Commands    int64     `json:"commands"`    // 命令执行数
	KeyCount    int64     `json:"key_count"`   // 键数量
	HitRate     float64   `json:"hit_rate"`    // 命中率
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
}

// APIMetricsDto API指标DTO
type APIMetricsDto struct {
	Path         string    `json:"path"`          // 接口路径
	Method       string    `json:"method"`        // 请求方法
	Count        int64     `json:"count"`         // 请求次数
	ErrorCount   int64     `json:"error_count"`   // 错误次数
	ResponseTime float64   `json:"response_time"` // 平均响应时间
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
}
