package cache

import "errors"

var (
	// ErrNotFound 记录未找到
	ErrNotFound = errors.New("cache: record not found")
	// ErrInvalidKey 无效的缓存键
	ErrInvalidKey = errors.New("cache: invalid key")
	// ErrInvalidValue 无效的缓存值
	ErrInvalidValue = errors.New("cache: invalid value")
	// ErrCacheMiss 缓存未命中
	ErrCacheMiss = errors.New("cache: cache miss")
)
