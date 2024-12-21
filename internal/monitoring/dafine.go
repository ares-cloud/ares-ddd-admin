package monitoring

import (
	monrest "github.com/ares-cloud/ares-ddd-admin/internal/monitoring/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type Server struct {
	metrice *monrest.MetricsController
}

func NewServer(
	metrice *monrest.MetricsController,
) *Server {
	return &Server{
		metrice: metrice,
	}
}

func (s *Server) Init(rg *route.RouterGroup, tk token.IToken) {
	s.metrice.RegisterRouter(rg, tk)
}
