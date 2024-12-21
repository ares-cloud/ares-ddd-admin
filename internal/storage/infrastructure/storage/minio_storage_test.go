package storage

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strings"
	"testing"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestMinioStorage_Upload(t *testing.T) {
	storage := createTestMinioStorage(t)

	// 测试上传文件
	content := "test content"
	reader := strings.NewReader(content)
	file, err := storage.Upload(context.Background(), reader, "test.txt", int64(len(content)), "test")

	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, model.StorageTypeMinio, file.StorageType)
}

func TestMinioStorage_Delete(t *testing.T) {
	storage := createTestMinioStorage(t)

	// 先上传文件
	content := "test content"
	reader := strings.NewReader(content)
	file, err := storage.Upload(context.Background(), reader, "test.txt", int64(len(content)), "test")
	assert.NoError(t, err)

	// 测试删除文件
	err = storage.Delete(context.Background(), file)
	assert.NoError(t, err)
}

func TestMinioStorage_GetURL(t *testing.T) {
	storage := createTestMinioStorage(t)

	// 测试公共URL
	file := &model.File{
		Path: "test/file.txt",
	}
	url, err := storage.GetURL(context.Background(), file)
	assert.NoError(t, err)
	assert.Contains(t, url, "http://localhost:9000")

	// 测试临时URL
	//storage = createTestMinioStorageWithoutPublicURL(t)
	//url, err = storage.GetURL(context.Background(), file)
	//assert.NoError(t, err)
	//assert.Contains(t, url, "X-Amz-Signature")
}

func TestMinioStorage_GetPreviewURL(t *testing.T) {
	storage := createTestMinioStorage(t)

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
			contains: "http://localhost:9000",
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

func createTestMinioStorage(t *testing.T) *MinioStorage {
	cof := configs.MinioStorage{
		Endpoint:  "localhost:10005",
		Bucket:    "go-dev",
		AccessKey: "root",
		SecretKey: "opeIM123",
		Region:    "",
		PublicURL: "http://localhost:9000",
	}

	// 2. 创建客户端
	client, err := minio.New(cof.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cof.AccessKey, cof.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 3. 创建存储实现
	return NewMinioStorage(
		client,
		cof.Bucket,
		cof.Region,
		cof.Endpoint,
		cof.PublicURL,
		cof.PublicURL,
	)
}
func Test_createTestMinioStorage(t *testing.T) {
	storage := createTestMinioStorage(t)
	t.Log(storage)
}

// ... 其他测试用例
