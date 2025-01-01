package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	domainstorage "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentStorage struct {
	client     *cos.Client
	bucket     string
	region     string
	publicURL  string
	previewURL string
}

func NewTencentStorage(client *cos.Client, bucket, region, publicURL, previewURL string) domainstorage.Storage {
	return &TencentStorage{
		client:     client,
		bucket:     bucket,
		region:     region,
		publicURL:  publicURL,
		previewURL: previewURL,
	}
}

func (s *TencentStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	objectName := path.Join(folderPath, filename)

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: size,
		},
	}

	_, err := s.client.Object.Put(ctx, objectName, file, opt)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	return &model.File{
		GeneratedName: filename,
		Path:          objectName,
		Size:          size,
		StorageType:   model.StorageTypeTencent,
		URL:           fmt.Sprintf("%s/%s", s.publicURL, objectName),
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}, nil
}

// Delete 删除文件
func (s *TencentStorage) Delete(ctx context.Context, file *model.File) error {
	// 使用生成的文件名构建路径
	_, err := s.client.Object.Delete(ctx, file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// Move 移动文件
func (s *TencentStorage) Move(ctx context.Context, file *model.File, toPath string) error {
	// 1. 构建源对象URL
	sourceURL := fmt.Sprintf("%s-%s.cos.%s.myqcloud.com/%s",
		s.bucket, s.client.BaseURL.BucketURL.Host, s.region, file.Path)

	// 2. 复制对象(使用生成的文件名)
	_, _, err := s.client.Object.Copy(ctx, toPath, sourceURL, nil)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	// 3. 删除源文件
	_, err = s.client.Object.Delete(ctx, file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}

// GetURL 获取文件访问URL
func (s *TencentStorage) GetURL(ctx context.Context, file *model.File) (string, error) {
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s", s.publicURL, file.Path), nil
	}

	// 生成临时URL
	presignedURL, err := s.client.Object.GetPresignedURL(ctx,
		http.MethodGet, file.Path,
		s.client.GetCredential().SecretID,
		s.client.GetCredential().SecretKey,
		time.Hour*24, nil)
	if err != nil {
		return "", herrors.NewErr(err)
	}

	return presignedURL.String(), nil
}

// GetPreviewURL 获取预览URL
func (s *TencentStorage) GetPreviewURL(ctx context.Context, file *model.File) (string, error) {
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
		return fmt.Sprintf("%s/preview?url=%s", s.previewURL, url.QueryEscape(fileURL)), nil
	}

	return "", herrors.NewBadReqError("unsupported preview type")
}

// Download 下载文件
func (s *TencentStorage) Download(ctx context.Context, file *model.File) (io.ReadCloser, error) {
	resp, err := s.client.Object.Get(ctx, file.Path, nil)
	if err != nil {
		return nil, herrors.NewErr(err)
	}
	return resp.Body, nil
}

// ... 实现其他接口方法
