package repeated_submit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
)

type RepeatedSubmitLock interface {
	// AcquireLock 获取锁
	AcquireLock(key string) bool
	// ReleaseLock 释放锁
	ReleaseLock(key string)
}

// Handler 重复提交的处理器
func Handler(rl RepeatedSubmitLock) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userId := actx.GetUserId(ctx)
		path := c.FullPath()
		lockKey := generateLockKey(path, userId)
		if userId != "" && path != "" {
			lock := rl.AcquireLock(lockKey)
			if !lock {
				i18Mag := hertzI18n.MustGetMessage(ctx, constant.PleaseDoNotResubmit)
				c.JSON(http.StatusOK, utils.H{constant.RespCode: http.StatusTooManyRequests, constant.RespMsg: i18Mag, constant.RespReason: constant.PleaseDoNotResubmit, constant.RespData: utils.H{}})
				c.Abort()
				return
			}
			defer rl.ReleaseLock(lockKey)
		}
		c.Next(ctx)
	}
}

func generateLockKey(path, userID string) string {
	return fmt.Sprintf("req_lock:%s:%s", path, userID)
}
