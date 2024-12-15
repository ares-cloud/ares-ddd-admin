package hserver

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

var (
	service *Serve
	once    sync.Once
)

type Serve struct {
	Env       string
	routers   []Router
	handlers  []app.HandlerFunc
	Tokenizer token.IToken
	config    *ServerConfig
	hertz     *server.Hertz
}

// NewServe 创建服务
func NewServe(config *ServerConfig, opts ...Option) *Serve {
	once.Do(func() {
		port := config.Port
		if port <= 0 {
			port = 8848
		}
		tracePort := config.TracerPort
		if port <= 0 {
			tracePort = 8849
		}
		addr := fmt.Sprintf(":%d", port)
		bodyMaxSize := 4
		if config.MaxRequestBodySize != 0 {
			bodyMaxSize = config.MaxRequestBodySize
		}
		h := server.Default(server.WithHostPorts(addr), server.WithMaxRequestBodySize(bodyMaxSize*1024*1024), server.WithTracer(prometheus.NewServerTracer(fmt.Sprintf(":%d", tracePort), "/hertz")))
		service = &Serve{
			hertz:  h,
			config: config,
		}
		for _, opt := range opts {
			opt(service)
		}
	})
	return service
}

// RegisterRouters 注册路由
func (s *Serve) RegisterRouters(routers ...Router) {
	s.routers = append(s.routers, routers...)
}

// Use 使用中间件
func (s *Serve) Use(handlers ...app.HandlerFunc) {
	s.handlers = append(s.handlers, handlers...)
}
func (s *Serve) GetHertz() *server.Hertz {
	return s.hertz
}

// Run 运行服务
func (s *Serve) Run() {
	//Register custom processors
	if len(s.handlers) > 0 {
		s.hertz.Use(s.handlers...)
	}
	// Set the routing of each module in sequence
	for _, r := range s.routers {
		r.ConfigRoutes(s.hertz, s.Tokenizer)
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		s.hertz.Spin()
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	hlog.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.hertz.Shutdown(ctx); err != nil {
		hlog.Fatal("Server forced to shutdown:", err)
	}
	hlog.Fatal("Server exiting")
}
