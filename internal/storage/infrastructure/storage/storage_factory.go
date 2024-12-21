package storage

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/redis/go-redis/v9"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/model"
	domainstorage "github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/storage"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type StorageFactory struct {
	config *configs.StorageConfig
	cache  *redis.Client
}

func NewStorageFactory(config *configs.StorageConfig, cache *h_redis.RedisClient) domainstorage.StorageFactory {
	return &StorageFactory{
		config: config,
		cache:  cache.GetClient(),
	}
}

// GetCurrentStorage 当前配置的
func (f *StorageFactory) GetCurrentStorage() (domainstorage.Storage, error) {
	return f.GetStorage(model.StorageType(f.config.Type))
}

func (f *StorageFactory) GetStorage(storageType model.StorageType) (domainstorage.Storage, error) {
	// 1. 创建基础存储
	var baseStorage domainstorage.Storage
	var err error

	switch storageType {
	case model.StorageTypeMinio:
		baseStorage, err = f.createMinioStorage()
	case model.StorageTypeAliyun:
		baseStorage, err = f.createAliyunStorage()
	case model.StorageTypeTencent:
		baseStorage, err = f.createTencentStorage()
	default:
		return nil, herrors.NewBadReqError("unsupported storage type")
	}

	if err != nil {
		return nil, err
	}

	// 2. 添加缓存装饰器
	if f.cache != nil {
		baseStorage = NewCacheStorage(baseStorage, f.cache, f.config.CacheTTL)
	}

	// 3. 添加监控装饰器
	baseStorage = NewMetricsStorage(baseStorage, string(storageType))

	return baseStorage, nil
}

func (f *StorageFactory) createMinioStorage() (domainstorage.Storage, error) {
	// 1. 验证配置
	if err := f.validateMinioConfig(); err != nil {
		return nil, err
	}

	// 2. 创建客户端
	client, err := minio.New(f.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(f.config.Minio.AccessKey, f.config.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 3. 创建存储实现
	return NewMinioStorage(
		client,
		f.config.Minio.Bucket,
		f.config.Minio.Region,
		f.config.Minio.Endpoint,
		f.config.Minio.PublicURL,
		f.config.PreviewURL,
	), nil
}

func (f *StorageFactory) validateMinioConfig() error {
	if f.config.Minio.Endpoint == "" {
		return herrors.NewBadReqError("minio endpoint is empty")
	}
	if f.config.Minio.AccessKey == "" {
		return herrors.NewBadReqError("minio access key is empty")
	}
	if f.config.Minio.SecretKey == "" {
		return herrors.NewBadReqError("minio secret key is empty")
	}
	if f.config.Minio.Bucket == "" {
		return herrors.NewBadReqError("minio bucket is empty")
	}
	return nil
}

// createAliyunStorage 创建阿里云存储
func (f *StorageFactory) createAliyunStorage() (domainstorage.Storage, error) {
	// 1. 验证配置
	if err := f.validateAliyunConfig(); err != nil {
		return nil, err
	}

	// 2. 创建阿里云OSS客户端
	client, err := oss.New(
		f.config.Aliyun.Endpoint,
		f.config.Aliyun.AccessKeyID,
		f.config.Aliyun.AccessKeySecret,
	)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 3. 获取Bucket
	bucket, err := client.Bucket(f.config.Aliyun.BucketName)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 4. 创建存储实现
	return NewAliyunStorage(
		client,
		bucket,
		f.config.Aliyun.Region,
		f.config.Aliyun.PublicURL,
		f.config.PreviewURL,
	), nil
}

// createTencentStorage 创建腾讯云存储
func (f *StorageFactory) createTencentStorage() (domainstorage.Storage, error) {
	// 1. 验证配置
	if err := f.validateTencentConfig(); err != nil {
		return nil, err
	}

	// 2. 创建腾讯云COS客户端
	bucketURL, err := cos.NewBucketURL(f.config.Tencent.Bucket, f.config.Tencent.Region, false)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	b := &cos.BaseURL{BucketURL: bucketURL}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  f.config.Tencent.SecretID,
			SecretKey: f.config.Tencent.SecretKey,
		},
	})

	// 3. 创建存储实现
	return NewTencentStorage(
		client,
		f.config.Tencent.Bucket,
		f.config.Tencent.Region,
		f.config.Tencent.PublicURL,
		f.config.PreviewURL,
	), nil
}

// 配置验证方法
func (f *StorageFactory) validateAliyunConfig() error {
	if f.config.Aliyun.Endpoint == "" {
		return herrors.NewBadReqError("aliyun endpoint is empty")
	}
	if f.config.Aliyun.AccessKeyID == "" {
		return herrors.NewBadReqError("aliyun access key id is empty")
	}
	if f.config.Aliyun.AccessKeySecret == "" {
		return herrors.NewBadReqError("aliyun access key secret is empty")
	}
	if f.config.Aliyun.BucketName == "" {
		return herrors.NewBadReqError("aliyun bucket name is empty")
	}
	return nil
}

func (f *StorageFactory) validateTencentConfig() error {
	if f.config.Tencent.SecretID == "" {
		return herrors.NewBadReqError("tencent secret id is empty")
	}
	if f.config.Tencent.SecretKey == "" {
		return herrors.NewBadReqError("tencent secret key is empty")
	}
	if f.config.Tencent.Region == "" {
		return herrors.NewBadReqError("tencent region is empty")
	}
	if f.config.Tencent.Bucket == "" {
		return herrors.NewBadReqError("tencent bucket is empty")
	}
	return nil
}
