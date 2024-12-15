package ratelimit

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"golang.org/x/time/rate"
)

// 定义一个IP限制器的结构体
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// 管理所有IP限制器
var visitors = make(map[string]*ipLimiter)
var mu sync.Mutex

// 清理旧的IP限制器
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// 获取特定IP的限制器
func getVisitor(ip string, r int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if v, exists := visitors[ip]; exists {
		v.lastSeen = time.Now()
		return v.limiter
	}

	limiter := rate.NewLimiter(1, r) // 每秒允许1个请求，最多积累3个请求
	visitors[ip] = &ipLimiter{limiter, time.Now()}

	return limiter
}

// RateLimitMiddleware 针对IP的速率限制中间件
func RateLimitMiddleware(r int) app.HandlerFunc {
	go cleanupVisitors()
	if r == 0 {
		r = 1
	}
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		limiter := getVisitor(ip, r)

		if !limiter.Allow() {
			c.JSON(http.StatusForbidden, utils.H{"code": http.StatusTooManyRequests, "msg": hertzI18n.MustGetMessage(ctx, "ServerBusy"), "reason": "ServerBusy"})
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
