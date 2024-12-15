package admin

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/cors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/sql_injection"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/hertz-contrib/gzip"
)

var ProviderSet = wire.NewSet(
	NewServer,
	rest.NewSysRoleController,
	rest.NewSysUserController,
	rest.NewSysTenantController,
	rest.NewSysPermissionsController,
)

func NewServer(config *configs.Bootstrap, hclient *h_redis.RedisClient,
	rc *rest.SysRoleController,
	uc *rest.SysUserController,
	ts *rest.SysTenantController,
	ps *rest.SysPermissionsController,
) *hserver.Serve {
	tk := token.NewRdbToken(hclient.GetClient(), config.JWT.Issuer, config.JWT.SigningKey, config.JWT.ExpirationToken, config.JWT.ExpirationRefresh)
	svr := hserver.NewServe(&hserver.ServerConfig{
		Port:               config.Server.Port,
		RateQPS:            config.Server.RateQPS,
		TracerPort:         config.Server.TracerPort,
		Name:               config.Server.Name,
		MaxRequestBodySize: config.Server.MaxRequestBodySize,
	}, hserver.WithTokenizer(tk))
	registerMiddleware(config, svr.GetHertz())
	//创建基础路由
	rg := svr.GetHertz().Group("/api/admin")
	rc.RegisterRouter(rg, tk)
	uc.RegisterRouter(rg, tk)
	ts.RegisterRouter(rg, tk)
	ps.RegisterRouter(rg, tk)
	return svr
}

func registerMiddleware(con *configs.Bootstrap, server *server.Hertz) {
	// Set up cross domain and flow limiting middleware
	server.Use(cors.Handler())
	//Use compression
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	//internationalization
	//server.Use(i18n.Handler(con.Server,language.Chinese, language.Chinese, language.English, language.TraditionalChinese))
	// server.Use(ratelimit.RateLimitMiddleware(10))
	// 防止sql注入
	server.Use(sql_injection.PreventSQLInjection())
}
