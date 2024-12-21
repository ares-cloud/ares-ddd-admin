package storage

import (
	"context"
	"io"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	ds "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	uploadDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_upload_duration_seconds",
			Help:    "Duration of file uploads in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)

	uploadSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_upload_size_bytes",
			Help:    "Size of uploaded files in bytes",
			Buckets: []float64{1024, 1024 * 1024, 10 * 1024 * 1024, 100 * 1024 * 1024},
		},
		[]string{"storage_type"},
	)

	operationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "storage_operation_errors_total",
			Help: "Total number of storage operation errors",
		},
		[]string{"storage_type", "operation"},
	)

	deleteDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_delete_duration_seconds",
			Help:    "Duration of file deletions in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)

	moveDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_move_duration_seconds",
			Help:    "Duration of file moves in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)

	downloadDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_download_duration_seconds",
			Help:    "Duration of file downloads in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)

	getURLDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_get_url_duration_seconds",
			Help:    "Duration of getting file URLs in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)

	getPreviewURLDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_get_preview_url_duration_seconds",
			Help:    "Duration of getting preview URLs in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"storage_type"},
	)
)

func init() {
	prometheus.MustRegister(uploadDuration)
	prometheus.MustRegister(uploadSize)
	prometheus.MustRegister(operationErrors)
	prometheus.MustRegister(deleteDuration)
	prometheus.MustRegister(moveDuration)
	prometheus.MustRegister(downloadDuration)
	prometheus.MustRegister(getURLDuration)
	prometheus.MustRegister(getPreviewURLDuration)
}

// 包装存储接口，添加监控指标
type metricsStorage struct {
	storage ds.Storage
	sType   string
}

func NewMetricsStorage(storage ds.Storage, storageType string) ds.Storage {
	return &metricsStorage{
		storage: storage,
		sType:   storageType,
	}
}

func (m *metricsStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	start := time.Now()
	f, err := m.storage.Upload(ctx, file, filename, size, folderPath)
	duration := time.Since(start).Seconds()

	uploadDuration.WithLabelValues(m.sType).Observe(duration)
	uploadSize.WithLabelValues(m.sType).Observe(float64(size))

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "upload").Inc()
	}

	return f, err
}

// 添加其他监控方法
func (m *metricsStorage) Delete(ctx context.Context, file *model.File) error {
	start := time.Now()
	err := m.storage.Delete(ctx, file)
	duration := time.Since(start).Seconds()

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "delete").Inc()
	}
	deleteDuration.WithLabelValues(m.sType).Observe(duration)

	return err
}

func (m *metricsStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	start := time.Now()
	err := m.storage.Move(ctx, file, oldPath)
	duration := time.Since(start).Seconds()

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "move").Inc()
	}
	moveDuration.WithLabelValues(m.sType).Observe(duration)

	return err
}

func (m *metricsStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	start := time.Now()
	reader, err := m.storage.Download(ctx, file)
	duration := time.Since(start).Seconds()

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "download").Inc()
	}
	downloadDuration.WithLabelValues(m.sType).Observe(duration)

	return reader, err
}

// GetURL 获取文件访问URL
func (m *metricsStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	start := time.Now()
	url, err := m.storage.GetURL(ctx, file)
	duration := time.Since(start).Seconds()

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "get_url").Inc()
	}

	getURLDuration.WithLabelValues(m.sType).Observe(duration)

	return url, err
}

// GetPreviewURL 获取预览URL
func (m *metricsStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	start := time.Now()
	url, err := m.storage.GetPreviewURL(ctx, file)
	duration := time.Since(start).Seconds()

	if err != nil {
		operationErrors.WithLabelValues(m.sType, "get_preview_url").Inc()
	}

	getPreviewURLDuration.WithLabelValues(m.sType).Observe(duration)

	return url, err
}
