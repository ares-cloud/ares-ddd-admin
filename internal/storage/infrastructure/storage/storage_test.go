package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Delete(ctx context.Context, file *model.File) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	args := m.Called(ctx, file, filename, size, folderPath)
	return args.Get(0).(*model.File), args.Error(1)
}

func TestCacheStorage_GetURL(t *testing.T) {
	mockStorage := new(MockStorage)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewCacheStorage(mockStorage, rdb, time.Minute)

	file := &model.File{
		Path: "test/file.txt",
	}

	// 测试缓存未命中
	expectedURL := "http://example.com/test/file.txt"
	mockStorage.On("GetURL", mock.Anything, file).Return(expectedURL, nil)

	url, err := cache.GetURL(context.Background(), file)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	// 测试缓存未命中
	url, err = cache.GetURL(context.Background(), file)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	mockStorage.AssertNumberOfCalls(t, "GetURL", 1)
}

func BenchmarkPoolStorage_Upload(b *testing.B) {
	mockStorage := new(MockStorage)
	pool := NewPoolStorage(mockStorage)

	content := strings.Repeat("test", 1024) // 4KB content
	file := strings.NewReader(content)

	mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.File{}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pool.Upload(context.Background(), file, "test.txt", int64(len(content)), "test")
		assert.NoError(b, err)
	}
}

func TestCacheStorage_Delete(t *testing.T) {
	mockStorage := new(MockStorage)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewCacheStorage(mockStorage, rdb, time.Minute)

	file := &model.File{
		Path: "test/file.txt",
	}

	// 先设置缓存
	rdb.Set(context.Background(), "storage:url:test/file.txt", "http://example.com/test.txt", time.Minute)

	// 测试删除
	mockStorage.On("Delete", mock.Anything, file).Return(nil)
	err := cache.Delete(context.Background(), file)
	assert.NoError(t, err)

	// 验证缓存已删除
	_, err = rdb.Get(context.Background(), "storage:url:test/file.txt").Result()
	assert.Equal(t, redis.Nil, err)
}

func TestPoolStorage_Download(t *testing.T) {
	mockStorage := new(MockStorage)
	pool := NewPoolStorage(mockStorage)

	file := &model.File{
		Path: "test/file.txt",
	}

	// 模拟下载
	content := strings.NewReader("test content")
	mockStorage.On("Download", mock.Anything, file).Return(io.NopCloser(content), nil)

	reader, err := pool.Download(context.Background(), file)
	assert.NoError(t, err)

	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "test content", string(data))
}
