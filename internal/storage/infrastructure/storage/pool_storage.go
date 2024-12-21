package storage

import (
	"context"
	"io"
	"sync"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	ds "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
)

// 对象池大小
const poolSize = 1024

type PoolStorage struct {
	storage ds.Storage
	pool    sync.Pool
}

func NewPoolStorage(storage ds.Storage) ds.Storage {
	return &PoolStorage{
		storage: storage,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 32*1024) // 32KB buffer
			},
		},
	}
}

func (p *PoolStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	// 从对象池获取缓冲区
	buf := p.pool.Get().([]byte)
	defer p.pool.Put(buf)

	// 使用缓冲区复制数据
	r := io.TeeReader(file, &bytesBuffer{buf: buf})
	return p.storage.Upload(ctx, r, filename, size, folderPath)
}

func (p *PoolStorage) Delete(ctx context.Context, file *model.File) error {
	return p.storage.Delete(ctx, file)
}

func (p *PoolStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	return p.storage.Move(ctx, file, oldPath)
}

func (p *PoolStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	return p.storage.GetURL(ctx, file)
}

func (p *PoolStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	return p.storage.GetPreviewURL(ctx, file)
}

func (p *PoolStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	reader, err := p.storage.Download(ctx, file)
	if err != nil {
		return nil, err
	}

	// 从对象池获取缓冲区
	buf := p.pool.Get().([]byte)

	// 创建带缓冲的读取器
	bufferedReader := &poolReader{
		reader: reader,
		buffer: buf,
		pool:   &p.pool,
	}

	return bufferedReader, nil
}

type bytesBuffer struct {
	buf []byte
	pos int
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	n = copy(b.buf[b.pos:], p)
	b.pos += n
	return n, nil
}

type poolReader struct {
	reader io.ReadCloser
	buffer []byte
	pool   *sync.Pool
}

func (r *poolReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *poolReader) Close() error {
	// 将缓冲区放回对象池
	r.pool.Put(r.buffer)
	return r.reader.Close()
}
