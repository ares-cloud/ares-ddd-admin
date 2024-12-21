package configs

import "time"

// StorageConfig 存储配置
type StorageConfig struct {
	Type          string         `mapstructure:"type"`           // 存储类型
	Local         LocalStorage   `mapstructure:"local"`          // 本地存储配置
	Minio         MinioStorage   `mapstructure:"minio"`          // Minio存储配置
	Aliyun        AliyunStorage  `mapstructure:"aliyun"`         // 阿里云存储配置
	Tencent       TencentStorage `mapstructure:"tencent"`        // 腾讯云存储配置
	Qiniu         QiniuStorage   `mapstructure:"qiniu"`          // 七牛云存储配置
	PreviewURL    string         `mapstructure:"preview_url"`    // 预览服务地址
	CacheTTL      time.Duration  `mapstructure:"cache_ttl"`      // 缓存过期时间
	RetentionDays int            `mapstructure:"retention_days"` // 保留天数
	Interval      time.Duration  `mapstructure:"interval"`       // 清理间隔
}

// LocalStorage 本地存储配置
type LocalStorage struct {
	RootPath   string `mapstructure:"root_path"`   // 存储根路径
	PublicPath string `mapstructure:"public_path"` // 公共访问路径
}

// MinioStorage Minio存储配置
type MinioStorage struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	PublicURL string `mapstructure:"public_url"`
}

// AliyunStorage 阿里云存储配置
type AliyunStorage struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
	PublicURL       string `mapstructure:"public_url"`
}

// TencentStorage 腾讯云存储配置
type TencentStorage struct {
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	PublicURL string `mapstructure:"public_url"`
}

// QiniuStorage 七牛云存储配置
type QiniuStorage struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
	Zone      string `mapstructure:"zone"`
}
