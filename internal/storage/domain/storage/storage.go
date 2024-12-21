package storage

import (
	"context"
	"io"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
)

// Storage 存储接口
// Storage 定义了文件存储的基本操作，包括上传、下载、删除、移动等功能。
// 该接口可以有多种实现，如本地存储、Minio、阿里云OSS、腾讯云COS等。
type Storage interface {
	// Upload 上传文件
	// 将文件上传到指定的存储路径。
	//
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - file: 文件内容的读取器
	//   - filename: 文件名，包含扩展名
	//   - size: 文件大小，单位字节
	//   - folderPath: 存储路径，不包含文件名
	//
	// 返回:
	//   - *model.File: 文件信息，包含ID、路径、URL等
	//   - error: 错误信息，如果上传成功则为nil
	Upload(ctx context.Context, file io.Reader, filename string, size int64, folderPath string) (*model.File, error)

	// Delete 删除文件
	// 从存储中删除指定的文件。
	//
	// 参数:
	//   - ctx: 上下文
	//   - file: 要删除的文件信息
	//
	// 返回:
	//   - error: 错误信息，如果删除成功则为nil
	Delete(ctx context.Context, file *model.File) error

	// Move 移动文件
	// 将文件从一个路径移动到另一个路径。
	//
	// 参数:
	//   - ctx: 上下文
	//   - file: 文件信息，包含新的路径
	//   - oldPath: 原始路径
	//
	// 返回:
	//   - error: 错误信息，如果移动成功则为nil
	Move(ctx context.Context, file *model.File, oldPath string) error

	// GetURL 获取文件访问URL
	// 获取文件的访问地址，可以是公共URL或带签名的临时URL。
	//
	// 参数:
	//   - ctx: 上下文
	//   - file: 文件信息
	//
	// 返回:
	//   - string: 文件访问URL
	//   - error: 错误信息
	GetURL(ctx context.Context, file *model.File) (string, error)

	// GetPreviewURL 获取预览URL
	// 获取文件的预览地址，支持图片直接预览和文档在线预览。
	//
	// 参数:
	//   - ctx: 上下文
	//   - file: 文件信息
	//
	// 返回:
	//   - string: 预览URL
	//   - error: 错误信息，如果文件类型不支持预览则返回错误
	GetPreviewURL(ctx context.Context, file *model.File) (string, error)

	// Download 下载文件
	// 获取文件的读取流，用于下载文件内容。
	//
	// 参数:
	//   - ctx: 上下文
	//   - file: 文件信息
	//
	// 返回:
	//   - io.ReadCloser: 文件内容读取器，使用完毕后需要调用Close()
	//   - error: 错误信息
	Download(ctx context.Context, file *model.File) (io.ReadCloser, error)
}

// StorageFactory 存储工厂接口
// StorageFactory 用于创建不同类型的存储实现。
// 工厂根据配置创建相应的存储实例，并添加缓存、监控等装饰器。
type StorageFactory interface {
	// GetStorage 获取存储实现
	// 根据存储类型创建对应的存储实例。
	//
	// 参数:
	//   - storageType: 存储类型，如minio、aliyun、tencent等
	//
	// 返回:
	//   - Storage: 存储接口实现
	//   - error: 错误信息，如果创建失败则返回错误
	GetStorage(storageType model.StorageType) (Storage, error)
	GetCurrentStorage() (Storage, error)
}
