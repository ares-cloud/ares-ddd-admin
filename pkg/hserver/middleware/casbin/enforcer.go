package casbin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

const (
	// Redis key前缀
	policyCacheKey = "casbin:policies:"
	// 权限更新channel
	policyUpdateChannel = "casbin:policy:update"
	// 缓存过期时间
	cacheExpiration = 24 * time.Hour
)

type Enforcer struct {
	enforcer *casbin.Enforcer
	permRepo IPermissionsRepository
	rdb      *redis.Client
	mutex    sync.RWMutex
	basePath string
}

// NewEnforcer 创建一个新的enforcer
func NewEnforcer(permRepo IPermissionsRepository, rdb *h_redis.RedisClient, basePath string) (*Enforcer, error) {
	// 从embed读取模型配置
	modelBytes, err := modelConf.ReadFile("model.conf")
	if err != nil {
		return nil, fmt.Errorf("read model config: %w", err)
	}

	// 加载模型
	m, err := model.NewModelFromString(string(modelBytes))
	if err != nil {
		return nil, fmt.Errorf("new model: %w", err)
	}

	// 创建适配器
	adapter := NewCasbinAdapter(permRepo)

	// 创建执行器
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("new enforcer: %w", err)
	}

	enforcer := &Enforcer{
		enforcer: e,
		permRepo: permRepo,
		rdb:      rdb.GetClient(),
		basePath: basePath,
	}

	// 启动订阅监听
	go enforcer.subscribeToUpdates()

	// 初始加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load policy: %w", err)
	}

	return enforcer, nil
}

// Enforce 执行权限检查 roleCode, tenantID, method, path
func (e *Enforcer) Enforce(role string, tenantID string, method string, path string) (bool, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	if e.basePath != "" {
		path = strings.TrimPrefix(path, e.basePath)
	}
	return e.enforcer.Enforce(role, tenantID, method, path)
}

// LoadPolicy 加载策略(从数据库加载并缓存)
func (e *Enforcer) LoadPolicy() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// 1. 从数据库加载
	if err := e.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("load from db: %w", err)
	}

	// 2. 将策略缓存到Redis
	ctx := context.Background()
	if err := e.cachePolicies(ctx); err != nil {
		hlog.Errorf("cache policies error: %v", err)
	}

	return nil
}

// ReloadPolicy 重新加载策略(始终从数据库加载)
func (e *Enforcer) ReloadPolicy() error {
	return e.LoadPolicy() // 直接调用LoadPolicy从数据库加载
}

// loadPoliciesFromCache 从缓存加载策略
func (e *Enforcer) loadPoliciesFromCache(ctx context.Context) ([][]string, error) {
	data, err := e.rdb.Get(ctx, policyCacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	var policies [][]string
	if err := json.Unmarshal(data, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

// cachePolicies 缓存策略到Redis
func (e *Enforcer) cachePolicies(ctx context.Context) error {
	policies, err := e.enforcer.GetPolicy()
	if err != nil {
		return fmt.Errorf("get policy: %w", err)
	}
	data, err := json.Marshal(policies)
	if err != nil {
		return err
	}

	return e.rdb.Set(ctx, policyCacheKey, data, cacheExpiration).Err()
}

// PublishUpdate 发布策略更新消息
func (e *Enforcer) PublishUpdate(ctx context.Context) error {
	return e.rdb.Publish(ctx, policyUpdateChannel, "update").Err()
}

// subscribeToUpdates 订阅策略更新消息
func (e *Enforcer) subscribeToUpdates() {
	ctx := context.Background()
	pubsub := e.rdb.Subscribe(ctx, policyUpdateChannel)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			hlog.Errorf("close pubsub error: %v", err)
		}
	}(pubsub)

	ch := pubsub.Channel()
	for range ch {
		if err := e.LoadPolicy(); err != nil { // 直接从数据库加载
			hlog.Errorf("reload policy error: %v", err)
		}
	}
}
