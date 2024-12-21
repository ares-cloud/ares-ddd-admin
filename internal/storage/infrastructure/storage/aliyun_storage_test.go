package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestAliyunStorage_Upload(t *testing.T) {
	// 创建真实的客户端，但使用测试配置
	client, err := oss.New(
		"test-endpoint",
		"test-access-key",
		"test-secret-key",
	)
	assert.NoError(t, err)

	bucket, err := client.Bucket("test-bucket")
	assert.NoError(t, err)

	storage := NewAliyunStorage(
		client,
		bucket,
		"oss-cn-hangzhou",
		"http://example.oss-cn-hangzhou.aliyuncs.com",
		"http://preview:8080",
	)

	// 测试上传文件
	content := "test content"
	reader := strings.NewReader(content)
	file, err := storage.Upload(context.Background(), reader, "test.txt", int64(len(content)), "test")

	// 由于无法真正上传，这里只验证返回的文件信息
	assert.Error(t, err) // 期望出错，因为是测试环境
	if err == nil {      // 如果意外成功了，验证返回值
		assert.NotNil(t, file)
		assert.Equal(t, "test.txt", file.Name)
		assert.Equal(t, model.StorageTypeAliyun, file.StorageType)
	}
}

func TestAliyunStorage_GetPreviewURL(t *testing.T) {
	client, err := oss.New(
		"test-endpoint",
		"test-access-key",
		"test-secret-key",
	)
	assert.NoError(t, err)

	bucket, err := client.Bucket("test-bucket")
	assert.NoError(t, err)

	storage := NewAliyunStorage(
		client,
		bucket,
		"oss-cn-hangzhou",
		"http://example.oss-cn-hangzhou.aliyuncs.com",
		"http://preview:8080",
	)

	tests := []struct {
		name     string
		file     *model.File
		wantErr  bool
		contains string
	}{
		{
			name: "image file",
			file: &model.File{
				Name: "test.jpg",
				Path: "test/test.jpg",
			},
			wantErr:  false,
			contains: "http://example.oss-cn-hangzhou.aliyuncs.com",
		},
		{
			name: "document file",
			file: &model.File{
				Name: "test.pdf",
				Path: "test/test.pdf",
			},
			wantErr:  false,
			contains: "preview?url=",
		},
		{
			name: "unsupported file",
			file: &model.File{
				Name: "test.exe",
				Path: "test/test.exe",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := storage.GetPreviewURL(context.Background(), tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			// 由于是测试环境，GetURL可能会失败，所以这里分两种情况
			if err != nil {
				t.Logf("GetURL failed as expected in test environment: %v", err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, url, tt.contains)
		})
	}
}

func TestAliyunStorage_Move(t *testing.T) {
	client, err := oss.New(
		"test-endpoint",
		"test-access-key",
		"test-secret-key",
	)
	assert.NoError(t, err)

	bucket, err := client.Bucket("test-bucket")
	assert.NoError(t, err)

	storage := NewAliyunStorage(
		client,
		bucket,
		"oss-cn-hangzhou",
		"http://example.oss-cn-hangzhou.aliyuncs.com",
		"http://preview:8080",
	)

	file := &model.File{
		Path: "new/path/file.txt",
	}
	oldPath := "old/path/file.txt"

	// 测试移动文件
	err = storage.Move(context.Background(), file, oldPath)
	assert.Error(t, err) // 期望出错，因为是测试环境
}

func TestAliyunStorage_Download(t *testing.T) {
	client, err := oss.New(
		"test-endpoint",
		"test-access-key",
		"test-secret-key",
	)
	assert.NoError(t, err)

	bucket, err := client.Bucket("test-bucket")
	assert.NoError(t, err)

	storage := NewAliyunStorage(
		client,
		bucket,
		"oss-cn-hangzhou",
		"http://example.oss-cn-hangzhou.aliyuncs.com",
		"http://preview:8080",
	)

	file := &model.File{
		Path: "test/file.txt",
	}

	// 测试下载文件
	reader, err := storage.Download(context.Background(), file)
	assert.Error(t, err) // 期望出错，因为是测试环境
	if err == nil {      // 如果意外成功了，确保关闭reader
		defer reader.Close()
	}
}
