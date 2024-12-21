package storage

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"io"
	"os"
	"path"
	"time"
)

type LocalStorage struct {
	rootPath   string
	publicPath string
	baseURL    string
}

func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error) {
	// 构建完整的存储路径
	fullPath := path.Join(s.rootPath, folderPath, filename)
	dirPath := path.Dir(fullPath)

	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, herrors.NewServerHError(err)
	}

	// 创建文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	defer dst.Close()

	// 复制文件内容
	written, err := io.Copy(dst, file)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}

	// 构建访问URL
	publicPath := path.Join(s.publicPath, folderPath, filename)
	publicURL := fmt.Sprintf("%s/%s", s.baseURL, publicPath)

	return &model.File{
		GeneratedName: filename, // 使用���入的生成文件名
		Path:          fullPath,
		Size:          written,
		StorageType:   model.StorageTypeLocal,
		URL:           publicURL,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}, nil
}

func (s *LocalStorage) Delete(ctx context.Context, file *model.File) error {
	// 使用生成的文件名构建路径
	err := os.Remove(file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

func (s *LocalStorage) Move(ctx context.Context, file *model.File, oldPath string) error {
	// 1. 确保目标目录存在
	targetDir := path.Dir(file.Path)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return herrors.NewServerHError(err)
	}

	// 2. 移动文件(使用生成的文件名)
	err := os.Rename(oldPath, file.Path)
	if err != nil {
		return herrors.NewServerHError(err)
	}

	return nil
}
