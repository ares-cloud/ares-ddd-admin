package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestTencentStorage_Upload(t *testing.T) {
	storage := createTestTencentStorage(t)

	// 测试上传文件
	content := "test content"
	reader := strings.NewReader(content)
	file, err := storage.Upload(context.Background(), reader, "test.txt", int64(len(content)), "test")

	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, model.StorageTypeTencent, file.StorageType)
}

func TestTencentStorage_GetPreviewURL(t *testing.T) {
	storage := createTestTencentStorage(t)

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
			contains: "http://test.cos.ap-guangzhou.myqcloud.com",
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
			assert.NoError(t, err)
			assert.Contains(t, url, tt.contains)
		})
	}
}

func createTestTencentStorage(t *testing.T) *TencentStorage {
	//client := createTestTencentClient(t)
	return nil
}

func createTestTencentClient(t *testing.T) *cos.Client {
	// 创建测试用的COS客户端
	//u := "https://test-bucket.cos.ap-guangzhou.myqcloud.com"
	//b := &cos.BaseURL{BucketURL: cos.NewBucketURL(u)}
	//return cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  "test-secret-id",
	//		SecretKey: "test-secret-key",
	//	},
	//})
	return nil
}
