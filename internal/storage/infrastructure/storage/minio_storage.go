package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"

	"net/url"
	"path/filepath"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

const (
	previewServiceURL = "http://preview-service:8080" // 预览服务地址
)

type MinioStorage struct {
	client     *minio.Client
	bucket     string
	region     string
	endpoint   string
	publicURL  string
	previewURL string // 添加预览服务地址
}

func NewMinioStorage(client *minio.Client, bucket, region, endpoint, publicURL string, previewURL string) *MinioStorage {
	return &MinioStorage{
		client:     client,
		bucket:     bucket,
		region:     region,
		endpoint:   endpoint,
		publicURL:  publicURL,
		previewURL: previewURL,
	}
}

func (s *MinioStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	objectName := path.Join(folderPath, filename)

	// 使用断点续传
	info, err := s.client.PutObject(ctx, s.bucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	return &model.File{
		GeneratedName: filename,
		Path:          objectName,
		Size:          info.Size,
		StorageType:   model.StorageTypeMinio,
		URL:           fmt.Sprintf("%s/%s", s.publicURL, objectName),
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}, nil
}

func (s *MinioStorage) Delete(ctx context.Context, file *model.File) error {
	err := s.client.RemoveObject(ctx, s.bucket, file.Path, minio.RemoveObjectOptions{})
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

func (s *MinioStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	// 如果配置了公共访问URL，直接返回
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, file.Path), nil
	}

	// 否则生成临时URL
	url, err := s.client.PresignedGetObject(ctx, s.bucket, file.Path, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// Move 移动文件
func (s *MinioStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	// 1. 检查源文件是否存在
	_, err := s.client.StatObject(ctx, s.bucket, oldPath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return herrors.NewBadReqError("source file does not exist")
		}
		return herrors.NewServerHError(err)
	}

	// 2. 构建源对象路径
	srcOpts := minio.CopySrcOptions{
		Bucket: s.bucket,
		Object: oldPath,
	}

	// 3. 构建目标对象路径
	dstOpts := minio.CopyDestOptions{
		Bucket: s.bucket,
		Object: file.Path,
	}

	// 4. 复制对象
	_, err = s.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return herrors.NewServerHError(fmt.Errorf("copy object error: %v", err))
	}

	// 5. 删除源文件
	err = s.client.RemoveObject(ctx, s.bucket, oldPath, minio.RemoveObjectOptions{})
	if err != nil {
		// 如果删除失败，尝试回滚复制操作
		_ = s.client.RemoveObject(ctx, s.bucket, file.Path, minio.RemoveObjectOptions{})
		return herrors.NewServerHError(fmt.Errorf("remove source object error: %v", err))
	}

	return nil
}

// GetPreviewURL 获取预览URL(根据文件类型生成不同的预览URL)
func (s *MinioStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	// 获取文件类型
	fileType := strings.ToLower(filepath.Ext(file.Name))

	// 图片类型直接返回访问URL
	if isImageType(fileType) {
		return s.GetURL(ctx, file)
	}

	// 文档类型返回预览服务URL
	if isDocumentType(fileType) {
		fileURL, err := s.GetURL(ctx, file)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/preview?url=%s", s.previewURL, url.QueryEscape(fileURL)), nil
	}

	// 其他类型不支持预览
	return "", herrors.NewBadReqError("unsupported preview type")
}

// Download 下载文件
func (s *MinioStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, file.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, herrors.NewErr(err)
	}
	return obj, nil
}

// 判断是否是图片类型
func isImageType(fileType string) bool {
	imageTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	return slices.Contains(imageTypes, fileType)
}

// 判断是否是文档类型
func isDocumentType(fileType string) bool {
	docTypes := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	return slices.Contains(docTypes, fileType)
}
