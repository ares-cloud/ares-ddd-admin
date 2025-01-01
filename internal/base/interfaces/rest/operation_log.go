package rest

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	_ "github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type OperationLogController struct {
	queryHandler *handlers.OperationLogQueryHandler
	ef           *casbin.Enforcer
}

func NewOperationLogController(queryHandler *handlers.OperationLogQueryHandler, ef *casbin.Enforcer) *OperationLogController {
	return &OperationLogController{
		queryHandler: queryHandler,
		ef:           ef,
	}
}

func (c *OperationLogController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	oplog := v1.Group("/oplog", jwt.Handler(t))
	{
		oplog.GET("/list", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListOperationLogQuery](c.List))
	}
}

// List 查询操作日志列表
// @Summary 查询操作日志列表
// @Description 查询操作日志列表
// @Tags 操作日志
// @Accept json
// @Produce json
// @Param tenant_id query string true "租户ID"
// @Param username query string false "用户名"
// @Param module query string false "模块"
// @Param action query string false "操作类型"
// @Param start_time query int64 false "开始时间"
// @Param end_time query int64 false "结束时间"
// @Param current query int false "页码"
// @Param size query int false "每页大小"
// @Success 200 {object} base_info.Success{data=[]dto.OperationLogDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/oplog/list [get]
func (c *OperationLogController) List(ctx context.Context, q *queries.ListOperationLogQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleList(ctx, q)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
