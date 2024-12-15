package cors

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
)

func Handler() app.HandlerFunc {
	return cors.New(cors.Config{
		// 允许跨源访问的 origin 列表
		AllowOrigins: []string{"http://*", "https://*"},
		// 允许客户端跨源访问所使用的 HTTP 方法列表
		AllowMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		// 允许使用的头信息字段列表
		AllowHeaders: []string{"*"},
		// 允许暴露给客户端的响应头列表
		ExposeHeaders: []string{"*"},
		// 允许客户端请求携带用户凭证
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	})
}
