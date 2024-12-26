package cache

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
)

type DepartmentCacheWarmer struct {
	repo  repository.IDepartmentRepository
	cache *DepartmentCache
}

func NewDepartmentCacheWarmer(repo repository.IDepartmentRepository, cache *DepartmentCache) *DepartmentCacheWarmer {
	return &DepartmentCacheWarmer{
		repo:  repo,
		cache: cache,
	}
}

// WarmUp 缓存预热
func (w *DepartmentCacheWarmer) WarmUp(ctx context.Context) error {
	// 1. 预热部门树
	//if err := w.warmupDeptTree(ctx); err != nil {
	//	return err
	//}
	//
	//// 2. 预热部门列表
	//if err := w.warmupDeptList(ctx); err != nil {
	//	return err
	//}

	return nil
}

// StartWarmupJob 启动定时预热任务
func (w *DepartmentCacheWarmer) StartWarmupJob(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := w.WarmUp(ctx); err != nil {
					hlog.CtxErrorf(ctx, "department cache warmup failed: %v", err)
				}
			}
		}
	}()
}

//func (w *DepartmentCacheWarmer) warmupDeptTree(ctx context.Context) error {
//	tree, err := w.repo.GetDepartmentTree(ctx, "")
//	if err != nil {
//		return err
//	}
//	return w.cache.SetDepartmentTree(ctx, "", tree)
//}
//
//func (w *DepartmentCacheWarmer) warmupDeptList(ctx context.Context) error {
//	depts, err := w.repo.List(ctx, &repository.ListDepartmentQuery{})
//	if err != nil {
//		return err
//	}
//	for _, dept := range depts {
//		if err := w.cache.SetDepartment(ctx, dept); err != nil {
//			return err
//		}
//	}
//	return nil
//}
