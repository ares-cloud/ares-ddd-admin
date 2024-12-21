package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type MetricsQueryHandler struct {
	service *service.MetricsService
}

func NewMetricsQueryHandler(service *service.MetricsService) *MetricsQueryHandler {
	return &MetricsQueryHandler{
		service: service,
	}
}

// HandleGetSystemMetrics 处理获取系统指标
func (h *MetricsQueryHandler) HandleGetSystemMetrics(ctx context.Context, q *queries.GetSystemMetricsQuery) (*dto.SystemMetricsDto, herrors.Herr) {
	metrics, err := h.service.GetSystemMetrics(ctx)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	return &dto.SystemMetricsDto{
		CPUUsage:    metrics.CPUUsage,
		MemoryUsage: metrics.MemoryUsage,
		DiskUsage:   metrics.DiskUsage,
		CreatedAt:   metrics.CreatedAt,
	}, nil
}

// HandleGetRuntimeMetrics 处理获取运行时指标
func (h *MetricsQueryHandler) HandleGetRuntimeMetrics(ctx context.Context) (*dto.RuntimeMetricsDto, herrors.Herr) {
	metrics := h.service.GetRuntimeMetrics()

	return &dto.RuntimeMetricsDto{
		Goroutines:    metrics.Goroutines,
		HeapAlloc:     metrics.HeapAlloc,
		HeapSys:       metrics.HeapSys,
		HeapObjects:   metrics.HeapObjects,
		StackInUse:    metrics.StackInUse,
		StackSys:      metrics.StackSys,
		MSpanInUse:    metrics.MSpanInUse,
		MSpanSys:      metrics.MSpanSys,
		MCacheInUse:   metrics.MCacheInUse,
		MCacheSys:     metrics.MCacheSys,
		GCPauseNs:     metrics.GCPauseNs,
		LastGC:        metrics.LastGC,
		NumGC:         metrics.NumGC,
		GCCPUFraction: metrics.GCCPUFraction,
	}, nil
}
