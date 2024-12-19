package oplog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Record 记录操作日志的中间件
func (l *Logger) Record(opt LogOption) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 获取请求信息
		log := &OperationLog{
			UserID:    actx.GetUserId(ctx),
			Username:  actx.GetUsername(ctx),
			TenantID:  actx.GetTenantId(ctx),
			Method:    string(c.Request.Method()),
			Path:      string(c.Request.URI().Path()),
			Query:     string(c.Request.URI().QueryString()),
			IP:        c.ClientIP(),
			UserAgent: string(c.Request.Header.UserAgent()),
			CreatedAt: startTime,
			Module:    opt.Module,
			Action:    opt.Action,
		}

		// 根据选项记录请求体
		if opt.IncludeBody {
			log.Body = string(c.Request.Body())
		}

		// 处理请求
		c.Next(ctx)

		// 记录响应信息
		log.Duration = time.Since(startTime).Milliseconds()
		log.Status = c.Response.StatusCode()

		// 记录错误信息
		if c.Response.StatusCode() != consts.StatusOK {
			var resp struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Reason  string `json:"reason"`
			}
			if err := json.Unmarshal(c.Response.Body(), &resp); err == nil {
				log.Error = resp.Message
			}
		}

		// 异步写入日志
		go func() {
			if err := l.writer.Write(context.Background(), log); err != nil {
				hlog.Errorf("write operation log error: %v", err)
			}
		}()
	}
}
