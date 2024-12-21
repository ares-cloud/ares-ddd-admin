package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	domainstorage "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
)

type AliyunStorage struct {
	client     *oss.Client
	bucket     *oss.Bucket
	region     string
	publicURL  string
	previewURL string
}

func NewAliyunStorage(client *oss.Client, bucket *oss.Bucket, region, publicURL, previewURL string) domainstorage.Storage {
	return &AliyunStorage{
		client:     client,
		bucket:     bucket,
		region:     region,
		publicURL:  publicURL,
		previewURL: previewURL,
	}
}

func (s *AliyunStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	objectName := path.Join(folderPath, filename)

	// 上传文件
	err := s.bucket.PutObject(objectName, file)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	return &model.File{
		GeneratedName: filename,
		Path:          objectName,
		Size:          size,
		StorageType:   model.StorageTypeAliyun,
		URL:           fmt.Sprintf("%s/%s", s.publicURL, objectName),
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}, nil
}

func (s *AliyunStorage) Delete(ctx context.Context, file *model.File) error {
	// 使用生成的文件名构建路径
	err := s.bucket.DeleteObject(file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

func (s *AliyunStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	// 1. 复制文件(使用生成的文件名)
	_, err := s.bucket.CopyObject(oldPath, file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 2. 删除源文件
	err = s.bucket.DeleteObject(oldPath)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}

func (s *AliyunStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s", s.publicURL, file.Path), nil
	}

	// 生成签名URL
	signedURL, err := s.bucket.SignURL(file.Path, oss.HTTPGet, 3600)
	if err != nil {
		return "", herrors.NewServerHError(err)
	}

	return signedURL, nil
}

func (s *AliyunStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
	fileType := strings.ToLower(path.Ext(file.Name))

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
		return fmt.Sprintf("%s/preview?url=%s", s.previewURL, fileURL), nil
	}

	return "", herrors.NewServerHError(fmt.Errorf("file type %s not supported", fileType))
}

func (s *AliyunStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	body, err := s.bucket.GetObject(file.Path)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return body, nil
}
