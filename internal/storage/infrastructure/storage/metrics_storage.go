package storage

//import (
//	"context"
//	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
//	"io"
//	"time"
//
//	"github.com/prometheus/client_golang/prometheus"
//)
//
//type MetricsStorage struct {
//	storage     Storage
//	storageType string
//}
//
//func (s *MetricsStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
//	startTime := time.Now()
//	fileModel, err := s.storage.Upload(ctx, file, filename, size, folderPath)
//	duration := time.Since(startTime)
//
//	// 记录上传指标
//	prometheus.WithLabelValues(s.storageType).Observe(duration.Seconds())
//	prometheus.WithLabelValues(s.storageType).Add(float64(size))
//	if err != nil {
//		prometheus.WithLabelValues(s.storageType).Inc()
//	}
//
//	return fileModel, err
//}
//
//func NewMetricsStorage(storage Storage, storageType string) *MetricsStorage {
//	return &MetricsStorage{
//		storage:     storage,
//		storageType: storageType,
//	}
//}
