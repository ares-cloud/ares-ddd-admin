package rest

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	_ "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/base_info"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type DepartmentController struct {
	cmdHandler   *handlers.DepartmentCommandHandler
	queryHandler *handlers.DepartmentQueryHandler
	ef           *casbin.Enforcer
	moduleName   string
}

func NewDepartmentController(cmdHandler *handlers.DepartmentCommandHandler, queryHandler *handlers.DepartmentQueryHandler, ef *casbin.Enforcer) *DepartmentController {
	return &DepartmentController{
		cmdHandler:   cmdHandler,
		queryHandler: queryHandler,
		ef:           ef,
		moduleName:   "部门",
	}
}

func (c *DepartmentController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	dept := v1.Group("/sys/dept", jwt.Handler(t))
	{
		dept.POST("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "新增",
		}), hserver.NewHandlerFu[commands.CreateDepartmentCommand](c.AddDepartment))

		dept.GET("", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListDepartmentsQuery](c.DepartmentList))

		dept.PUT("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "修改",
		}), hserver.NewHandlerFu[commands.UpdateDepartmentCommand](c.UpdateDepartment))

		dept.DELETE("/:id", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](c.DeleteDepartment))

		dept.GET("/:id", casbin.Handler(c.ef), hserver.NewHandlerFu[models.StringIdReq](c.GetDetails))

		dept.GET("/tree", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.GetDepartmentTreeQuery](c.GetDepartmentTree))

		dept.POST("/move", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "移动",
		}), hserver.NewHandlerFu[commands.MoveDepartmentCommand](c.MoveDepartment))
	}
}

// AddDepartment 添加部门
// @Summary 添加部门
// @Description 添加部门
// @Tags 系统部门
// @ID AddDepartment
// @Accept json
// @Produce json
// @Param req body commands.CreateDepartmentCommand true "部门信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [post]
func (c *DepartmentController) AddDepartment(ctx context.Context, params *commands.CreateDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleCreate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DepartmentList 获取部门列表
// @Summary 获取部门列表
// @Description 获取部门列表
// @Tags 系统部门
// @ID DepartmentList
// @Accept json
// @Produce json
// @Param req query queries.ListDepartmentsQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=[]model.Department}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [get]
func (c *DepartmentController) DepartmentList(ctx context.Context, params *queries.ListDepartmentsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdateDepartment 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 系统部门
// @ID UpdateDepartment
// @Accept json
// @Produce json
// @Param req body commands.UpdateDepartmentCommand true "部门信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [put]
func (c *DepartmentController) UpdateDepartment(ctx context.Context, params *commands.UpdateDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleUpdate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Description 删除部门
// @Tags 系统部门
// @ID DeleteDepartment
// @Accept json
// @Produce json
// @Param id path string true "部门ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/{id} [delete]
func (c *DepartmentController) DeleteDepartment(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleDelete(ctx, &commands.DeleteDepartmentCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetDetails 获取部门详情
// @Summary 获取部门详情
// @Description 获取部门详情
// @Tags 系统部门
// @ID GetDepartmentDetails
// @Accept json
// @Produce json
// @Param id path string true "部门ID"
// @Success 200 {object} base_info.Success{data=model.Department}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/{id} [get]
func (c *DepartmentController) GetDetails(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGet(ctx, &queries.GetDepartmentQuery{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetDepartmentTree 获取部门树
// @Summary 获取部门树
// @Description 获取部门树形结构
// @Tags 系统部门
// @ID GetDepartmentTree
// @Accept json
// @Produce json
// @Param req query queries.GetDepartmentTreeQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=[]model.Department}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/tree [get]
func (c *DepartmentController) GetDepartmentTree(ctx context.Context, params *queries.GetDepartmentTreeQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGetTree(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// MoveDepartment 移动部门
// @Summary 移��部门
// @Description 移动部门位置
// @Tags 系统部门
// @ID MoveDepartment
// @Accept json
// @Produce json
// @Param req body commands.MoveDepartmentCommand true "移动信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/move [post]
func (c *DepartmentController) MoveDepartment(ctx context.Context, params *commands.MoveDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleMove(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}
