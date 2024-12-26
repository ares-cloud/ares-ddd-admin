package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// CacheHitCounter 缓存命中计数器
	CacheHitCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Cache hit counter.",
		},
		[]string{"type", "operation"},
	)

	// CacheMissCounter 缓存未命中计数器
	CacheMissCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Cache miss counter.",
		},
		[]string{"type", "operation"},
	)

	// CacheErrorCounter 缓存错误计数器
	CacheErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_error_total",
			Help: "Cache error counter.",
		},
		[]string{"type", "operation"},
	)

	// CacheLatencyHistogram 缓存延迟直方图
	CacheLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_latency_seconds",
			Help:    "Cache operation latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type", "operation"},
	)
)

func init() {
	// 注册监控指标
	prometheus.MustRegister(CacheHitCounter)
	prometheus.MustRegister(CacheMissCounter)
	prometheus.MustRegister(CacheErrorCounter)
	prometheus.MustRegister(CacheLatencyHistogram)
}
