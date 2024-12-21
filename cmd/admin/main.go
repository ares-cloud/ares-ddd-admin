package main

import (
	"flag"
	_ "github.com/ares-cloud/ares-ddd-admin/docs/admin"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	// flagconf is the config flag.
	flagConf string
	// flagconf is the config flag.
	env string
	// flagLocalize is the Localize flag.
	flagLocalize string
	// flagLog is the login flag.
	flagLog string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&env, "env", "dev", "Operating environment, eg: -env dev")
	flag.StringVar(&env, "lc", "", "localize config path, eg: -lc localize")
	flag.StringVar(&flagLog, "log", "", "Operating environment, eg: -log app.log")
}

type app struct {
	server *hserver.Serve
}

func newApp(server *hserver.Serve) *app {
	return &app{
		server: server,
	}
}

// @title ares-ddd-admin
// @version 1.0
// @description This is a demo using go-server-template-admin.

// @contact.name go-server-template-admin

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /api/admin
// @schemes http
func main() {
	flag.Parse()
	//load config
	bootstrap, err := configs.Load(flagConf, flagLocalize, env, flagLog)
	if err != nil {
		panic(err)
	}
	bootstrap.ConfPath = &flagConf
	application, cleanup, err := wireApp(bootstrap, bootstrap.Data, bootstrap.Storage)
	if err != nil {
		panic(err)
	}
	url := swagger.URL("/swagger/doc.json") // The url pointing to API definition
	application.server.GetHertz().GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	defer cleanup()
	application.server.Run()
}
