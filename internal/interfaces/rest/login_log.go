package rest

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type LoginLogController struct {
	queryHandler *handlers.LoginLogQueryHandler
	ef           *casbin.Enforcer
}

func NewLoginLogController(queryHandler *handlers.LoginLogQueryHandler, ef *casbin.Enforcer) *LoginLogController {
	return &LoginLogController{
		queryHandler: queryHandler,
		ef:           ef,
	}
}

func (c *LoginLogController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	lg := v1.Group("/sys/login-log", jwt.Handler(t))
	{
		lg.GET("admin", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListLoginLogsQuery](c.AdminLogList))
		lg.GET("app", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListLoginLogsQuery](c.AppList))
	}
}

// AdminLogList 查询管理端登录日志列表
// @Summary 查询管理端登录日志列表
// @Description 查询管理端登录日志列表
// @Tags 登录日志
// @ID ListAdminLoginLogs
// @Accept json
// @Produce json
// @Param month query string true "查询月份(格式:202403)"
// @Param current query int false "页码"
// @Param size query int false "每页大小"
// @Param username query string false "用户名"
// @Param ip query string false "登录IP"
// @Param status query int false "登录状态"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} base_info.Success{data=[]dto.LoginLogDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/sys/login-log/admin [get]
func (c *LoginLogController) AdminLogList(ctx context.Context, params *queries.ListLoginLogsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleAdminList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// AppList 查询app登录日志列表
// @Summary 查询app登录日志列表
// @Description 查询app登录日志列表
// @Tags 登录日志
// @ID ListAppLoginLogs
// @Accept json
// @Produce json
// @Param month query string true "查询月份(格式:202403)"
// @Param page query int false "页码"
// @Param page_size query int false "每页大小"
// @Param username query string false "用户名"
// @Param ip query string false "登录IP"
// @Param status query int false "登录状态"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} base_info.Success{data=[]dto.LoginLogDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/sys/login-log/app [get]
func (c *LoginLogController) AppList(ctx context.Context, params *queries.ListLoginLogsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleAppList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
