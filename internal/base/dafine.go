package base

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/handlers"
	baserest "github.com/ares-cloud/ares-ddd-admin/internal/base/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type BaseServer struct {
	rc           *baserest.SysRoleController
	uc           *baserest.SysUserController
	ts           *baserest.SysTenantController
	ps           *baserest.SysPermissionsController
	as           *baserest.AuthController
	lls          *baserest.LoginLogController
	ols          *baserest.OperationLogController
	des          *baserest.DepartmentController
	dps          *baserest.DataPermissionController
	handlerEvent *handlers.HandlerEvent
}

func NewBaseServer(
	rc *baserest.SysRoleController,
	uc *baserest.SysUserController,
	ts *baserest.SysTenantController,
	ps *baserest.SysPermissionsController,
	as *baserest.AuthController,
	lls *baserest.LoginLogController,
	ols *baserest.OperationLogController,
	des *baserest.DepartmentController,
	dps *baserest.DataPermissionController,
	handlerEvent *handlers.HandlerEvent,
) *BaseServer {
	return &BaseServer{
		rc:           rc,
		uc:           uc,
		ts:           ts,
		ps:           ps,
		as:           as,
		lls:          lls,
		ols:          ols,
		des:          des,
		dps:          dps,
		handlerEvent: handlerEvent,
	}
}

func (s *BaseServer) Init(rg *route.RouterGroup, tk token.IToken) {
	s.rc.RegisterRouter(rg, tk)
	s.uc.RegisterRouter(rg, tk)
	s.ts.RegisterRouter(rg, tk)
	s.ps.RegisterRouter(rg, tk)
	s.as.RegisterRouter(rg, tk)
	s.lls.RegisterRouter(rg, tk)
	s.ols.RegisterRouter(rg, tk)
	s.des.RegisterRouter(rg, tk)
	s.dps.RegisterRouter(rg, tk)
	s.handlerEvent.Register()
}
