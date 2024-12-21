package service

import (
	"context"
	"runtime"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/domain/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type MetricsService struct{}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// GetSystemMetrics 获取系统指标
func (s *MetricsService) GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	// CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 磁盘使用率
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &model.SystemMetrics{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memInfo.UsedPercent,
		DiskUsage:   diskInfo.UsedPercent,
		CreatedAt:   time.Now(),
	}, nil
}

// GetRuntimeMetrics 获取运行时指标
func (s *MetricsService) GetRuntimeMetrics() *model.RuntimeMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &model.RuntimeMetrics{
		Goroutines:    runtime.NumGoroutine(),
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapObjects:   m.HeapObjects,
		StackInUse:    m.StackInuse,
		StackSys:      m.StackSys,
		MSpanInUse:    m.MSpanInuse,
		MSpanSys:      m.MSpanSys,
		MCacheInUse:   m.MCacheInuse,
		MCacheSys:     m.MCacheSys,
		GCPauseNs:     m.PauseNs[(m.NumGC+255)%256],
		LastGC:        m.LastGC,
		NumGC:         m.NumGC,
		GCCPUFraction: m.GCCPUFraction,
	}
}
