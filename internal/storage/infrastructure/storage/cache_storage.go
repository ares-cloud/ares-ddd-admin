package storage

import (
	"context"
	"fmt"
	ds "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"github.com/redis/go-redis/v9"
	"io"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
)

type CacheStorage struct {
	storage ds.Storage
	cache   *redis.Client
	ttl     time.Duration
}

func NewCacheStorage(storage ds.Storage, cache *redis.Client, ttl time.Duration) ds.Storage {
	return &CacheStorage{
		storage: storage,
		cache:   cache,
		ttl:     ttl,
	}
}

func (c *CacheStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	// 尝试从缓存获取
	key := fmt.Sprintf("storage:url:%s", file.Path)
	if url, err := c.cache.Get(ctx, key).Result(); err == nil {
		return url, nil
	}

	// 从存储获取
	url, err := c.storage.GetURL(ctx, file)
	if err != nil {
		return "", err
	}

	// 写入缓存
	c.cache.Set(ctx, key, url, c.ttl)
	return url, nil
}

func (c *CacheStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	// 尝试从缓存获取
	key := fmt.Sprintf("storage:preview:%s", file.Path)
	if url, err := c.cache.Get(ctx, key).Result(); err == nil {
		return url, nil
	}

	// 从存储获取
	url, err := c.storage.GetPreviewURL(ctx, file)
	if err != nil {
		return "", err
	}

	// 写入缓存
	c.cache.Set(ctx, key, url, c.ttl)
	return url, nil
}

// Upload 上传文件(不缓存)
func (c *CacheStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	return c.storage.Upload(ctx, file, filename, size, folderPath)
}

// Delete 删除文件(同时删除缓存)
func (c *CacheStorage) Delete(ctx context.Context, file *model.File) error {
	err := c.storage.Delete(ctx, file)
	if err != nil {
		return err
	}

	// 删除缓存
	c.cache.Del(ctx, fmt.Sprintf("storage:url:%s", file.Path))
	c.cache.Del(ctx, fmt.Sprintf("storage:preview:%s", file.Path))
	return nil
}

// Move 移动文件(同时更新缓存)
func (c *CacheStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	err := c.storage.Move(ctx, file, oldPath)
	if err != nil {
		return err
	}

	// 删除旧缓存
	c.cache.Del(ctx, fmt.Sprintf("storage:url:%s", oldPath))
	c.cache.Del(ctx, fmt.Sprintf("storage:preview:%s", oldPath))
	return nil
}

// Download 下载文件(不缓存)
func (c *CacheStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	return c.storage.Download(ctx, file)
}

// ... 实现其他接口方法
