package cache

//
//import (
//	"context"
//	"github.com/bytedance/gopkg/util/logger"
//	"github.com/cloudwego/hertz/pkg/common/hlog"
//	"sync"
//	"time"
//
//	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
//	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
//	"github.com/ares-cloud/ares-ddd-admin/internal/base/metrics"
//)
//
//type DepartmentCacheLoader struct {
//	repo  repository.IDepartmentRepository
//	cache *DepartmentCache
//}
//
//func NewDepartmentCacheLoader(repo repository.IDepartmentRepository, cache *DepartmentCache) *DepartmentCacheLoader {
//	return &DepartmentCacheLoader{
//		repo:  repo,
//		cache: cache,
//	}
//}
//
//// PreloadAll 预加载所有缓存
//func (l *DepartmentCacheLoader) PreloadAll(ctx context.Context) error {
//	start := time.Now()
//	hlog.CtxInfof(ctx, "start preloading department cache...")
//
//	var wg sync.WaitGroup
//	errChan := make(chan error, 3)
//
//	// 1. 预加载部门树
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := l.preloadDeptTree(ctx); err != nil {
//			errChan <- err
//		}
//	}()
//
//	// 2. 预加载部门列表
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := l.preloadDeptList(ctx); err != nil {
//			errChan <- err
//		}
//	}()
//
//	// 3. 预加载用户部门关系
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		if err := l.preloadUserDepts(ctx); err != nil {
//			errChan <- err
//		}
//	}()
//
//	// 等待所有预加载完成
//	wg.Wait()
//	close(errChan)
//
//	// 检查错误
//	for err := range errChan {
//		if err != nil {
//			metrics.CacheErrorCounter.WithLabelValues("department", "preload").Inc()
//			return err
//		}
//	}
//
//	// 记录预加载耗时
//	duration := time.Since(start)
//	metrics.CacheLatencyHistogram.WithLabelValues("department", "preload").Observe(duration.Seconds())
//	logger.Infof("department cache preload completed, took %v", duration)
//
//	return nil
//}
//
//// StartPreloadJob 启动定时预加载任务
//func (l *DepartmentCacheLoader) StartPreloadJob(ctx context.Context) {
//	go func() {
//		// 启动时先执行一次预加载
//		if err := l.PreloadAll(ctx); err != nil {
//			logger.Errorf("initial department cache preload failed: %v", err)
//		}
//
//		// 定时预加载
//		ticker := time.NewTicker(12 * time.Hour)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case <-ticker.C:
//				if err := l.PreloadAll(ctx); err != nil {
//					logger.Errorf("department cache preload failed: %v", err)
//				}
//			}
//		}
//	}()
//}
//
//// preloadDeptTree 预加载部���树
//func (l *DepartmentCacheLoader) preloadDeptTree(ctx context.Context) error {
//	tree, err := l.repo.GetDepartmentTree(ctx, "")
//	if err != nil {
//		return err
//	}
//
//	// 设置较长的过期时间
//	return l.cache.SetDepartmentTreeWithExpire(ctx, "", tree, 24*time.Hour)
//}
//
//// preloadDeptList 预加载部门列表
//func (l *DepartmentCacheLoader) preloadDeptList(ctx context.Context) error {
//	depts, err := l.repo.List(ctx, "", "", nil)
//	if err != nil {
//		return err
//	}
//
//	// 并发写入缓存
//	errChan := make(chan error, len(depts))
//	var wg sync.WaitGroup
//
//	for _, dept := range depts {
//		wg.Add(1)
//		go func(d *model.Department) {
//			defer wg.Done()
//			if err := l.cache.SetDepartmentWithExpire(ctx, d, 24*time.Hour); err != nil {
//				errChan <- err
//			}
//		}(dept)
//	}
//
//	wg.Wait()
//	close(errChan)
//
//	// 检查错误
//	for err := range errChan {
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// preloadUserDepts 预加载用户部门关系
//func (l *DepartmentCacheLoader) preloadUserDepts(ctx context.Context) error {
//	// 获取所有用户ID
//	userIDs, err := l.repo.GetAllUserIDs(ctx)
//	if err != nil {
//		return err
//	}
//
//	// 并发预加载每个用户的部门关系
//	errChan := make(chan error, len(userIDs))
//	var wg sync.WaitGroup
//
//	for _, userID := range userIDs {
//		wg.Add(1)
//		go func(uid string) {
//			defer wg.Done()
//			depts, err := l.repo.GetUserDepartments(ctx, uid)
//			if err != nil {
//				errChan <- err
//				return
//			}
//			if err := l.cache.SetUserDepartments(ctx, uid, depts); err != nil {
//				errChan <- err
//			}
//		}(userID)
//	}
//
//	wg.Wait()
//	close(errChan)
//
//	// 检查错误
//	for err := range errChan {
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
