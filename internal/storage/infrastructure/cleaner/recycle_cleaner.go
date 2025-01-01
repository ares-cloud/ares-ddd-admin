package cleaner

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/repository"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/service"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type RecycleCleaner struct {
	repo    repository.IStorageRepos
	service *service.StorageService
	// 保留天数
	retentionDays int
	// 清理间隔
	interval time.Duration
	// 停止信号
	stopChan chan struct{}
}

func NewRecycleCleaner(repo repository.IStorageRepos, service *service.StorageService, conf *configs.StorageConfig) *RecycleCleaner {
	return &RecycleCleaner{
		repo:          repo,
		service:       service,
		retentionDays: conf.RetentionDays,
		interval:      conf.Interval,
		stopChan:      make(chan struct{}),
	}
}

// Start 启动清理任务
func (c *RecycleCleaner) Start() {
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.clean()
			case <-c.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止清理任务
func (c *RecycleCleaner) Stop() {
	close(c.stopChan)
}

// clean 执行清理
func (c *RecycleCleaner) clean() {
	ctx := context.Background()

	// 计算过期时间
	expireTime := time.Now().AddDate(0, 0, -c.retentionDays)

	// 获取过期文件
	files, err := c.repo.GetExpiredRecycleFiles(ctx, expireTime)
	if err != nil {
		hlog.Errorf("get expired recycle files error: %v", err)
		return
	}

	// 删除文件
	for _, file := range files {
		if err := c.service.DeleteFile(ctx, file.ID); err != nil {
			hlog.Errorf("delete expired file error: %v", err)
			continue
		}
	}
}
