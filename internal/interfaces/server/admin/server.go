package admin

import (
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/pkg/h_redis"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/i18n"
	psb "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/cors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/sql_injection"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/hertz-contrib/gzip"
	"golang.org/x/text/language"
)

var ProviderSet = wire.NewSet(
	NewServer,
	rest.NewSysRoleController,
	rest.NewSysUserController,
	rest.NewSysTenantController,
	rest.NewSysPermissionsController,
	rest.NewAuthController,
	rest.NewLoginLogController,
	NewCasBinEnforcer,
)

const baseUrl = "/api/admin"

func NewServer(config *configs.Bootstrap, hc *h_redis.RedisClient,
	rc *rest.SysRoleController,
	uc *rest.SysUserController,
	ts *rest.SysTenantController,
	ps *rest.SysPermissionsController,
	as *rest.AuthController,
	lls *rest.LoginLogController,
) *hserver.Serve {
	tk := token.NewRdbToken(hc.GetClient(), config.JWT.Issuer, config.JWT.SigningKey, config.JWT.ExpirationToken, config.JWT.ExpirationRefresh)
	svr := hserver.NewServe(&hserver.ServerConfig{
		Port:               config.Server.Port,
		RateQPS:            config.Server.RateQPS,
		TracerPort:         config.Server.TracerPort,
		Name:               config.Server.Name,
		MaxRequestBodySize: config.Server.MaxRequestBodySize,
	}, hserver.WithTokenizer(tk))
	registerMiddleware(config, svr.GetHertz())
	//创建基础路由
	rg := svr.GetHertz().Group(baseUrl)
	rc.RegisterRouter(rg, tk)
	uc.RegisterRouter(rg, tk)
	ts.RegisterRouter(rg, tk)
	ps.RegisterRouter(rg, tk)
	as.RegisterRouter(rg, tk)
	lls.RegisterRouter(rg, tk)
	return svr
}

func NewCasBinEnforcer(hc *h_redis.RedisClient, pr psb.IPermissionsRepository) (*psb.Enforcer, error) {
	enforcer, err := psb.NewEnforcer(pr, hc, baseUrl)
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}

func registerMiddleware(con *configs.Bootstrap, server *server.Hertz) {
	// Set up cross domain and flow limiting middleware
	server.Use(cors.Handler())
	//Use compression
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	//internationalization
	if con.ConfPath != nil {
		ph := fmt.Sprintf("%s/localize", *con.ConfPath)
		server.Use(i18n.Handler(ph, language.Chinese, language.Chinese, language.English, language.TraditionalChinese))
	}
	// server.Use(ratelimit.RateLimitMiddleware(10))
	// 防止sql注入
	server.Use(sql_injection.PreventSQLInjection())

	// 操作日志
	initOpLog(con.Log)
}

func initOpLog(con *configs.Log) {
	path := con.OutPath
	if path == "" {
		panic(fmt.Errorf("not config log out path"))
	}
	writer := oplog.NewFileWriter(path)
	oplog.Init(writer)
}
