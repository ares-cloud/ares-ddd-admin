package model

// RuntimeMetrics Go运行时指标
type RuntimeMetrics struct {
	Goroutines    int     // Goroutine数量
	HeapAlloc     uint64  // 已分配堆内存
	HeapSys       uint64  // 系统预留堆内存
	HeapObjects   uint64  // 堆对象数量
	StackInUse    uint64  // 正在使用的栈内存
	StackSys      uint64  // 系统预留栈内存
	MSpanInUse    uint64  // 正在使用的MSpan内存
	MSpanSys      uint64  // 系统预留MSpan内存
	MCacheInUse   uint64  // 正在使用的MCache内存
	MCacheSys     uint64  // 系统预留MCache内存
	GCPauseNs     uint64  // 最后一次GC暂停时间(纳秒)
	LastGC        uint64  // 上次GC时间
	NumGC         uint32  // GC次数
	GCCPUFraction float64 // GC占用CPU时间比例
}

// NewRuntimeMetrics 创建运行时指标
func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{}
}

// SetGoroutines 设置Goroutine数量
func (m *RuntimeMetrics) SetGoroutines(n int) {
	m.Goroutines = n
}

// SetHeapStats 设置堆内存统计
func (m *RuntimeMetrics) SetHeapStats(alloc, sys, objects uint64) {
	m.HeapAlloc = alloc
	m.HeapSys = sys
	m.HeapObjects = objects
}

// SetStackStats 设置栈内存统计
func (m *RuntimeMetrics) SetStackStats(inuse, sys uint64) {
	m.StackInUse = inuse
	m.StackSys = sys
}

// SetMSpanStats 设置MSpan内存统计
func (m *RuntimeMetrics) SetMSpanStats(inuse, sys uint64) {
	m.MSpanInUse = inuse
	m.MSpanSys = sys
}

// SetMCacheStats 设置MCache内存统计
func (m *RuntimeMetrics) SetMCacheStats(inuse, sys uint64) {
	m.MCacheInUse = inuse
	m.MCacheSys = sys
}

// SetGCStats 设置GC统计
func (m *RuntimeMetrics) SetGCStats(pauseNs, lastGC uint64, numGC uint32, cpuFraction float64) {
	m.GCPauseNs = pauseNs
	m.LastGC = lastGC
	m.NumGC = numGC
	m.GCCPUFraction = cpuFraction
}
