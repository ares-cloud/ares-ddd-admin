package rest

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	_ "github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	_ "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/base_info"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type SysPermissionsController struct {
	cmdHandel   *handlers.PermissionsCommandHandler
	queryHandel *handlers.PermissionsQueryHandler
	ef          *casbin.Enforcer
	modeNma     string
}

func NewSysPermissionsController(cmdHandel *handlers.PermissionsCommandHandler, queryHandel *handlers.PermissionsQueryHandler, ef *casbin.Enforcer) *SysPermissionsController {
	return &SysPermissionsController{
		cmdHandel:   cmdHandel,
		queryHandel: queryHandel,
		ef:          ef,
		modeNma:     "系统权限",
	}
}

func (c *SysPermissionsController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	ur := v1.Group("/sys/permissions", jwt.Handler(t))
	{
		ur.POST("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "新增",
		}), hserver.NewHandlerFu[commands.CreatePermissionsCommand](c.AddPermissions))
		ur.GET("", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListPermissionsQuery](c.PermissionsList))
		ur.PUT("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "修改",
		}), hserver.NewHandlerFu[commands.UpdatePermissionsCommand](c.UpdatePermissions))
		ur.DELETE("/:id", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "删除",
		}), casbin.Handler(c.ef), hserver.NewHandlerFu[models.IntIdReq](c.DeletePermissions))
		ur.GET("/tree", hserver.NewHandlerFu[queries.GetPermissionsTreeQuery](c.GetPermissionsTree))
		ur.GET("/simple/tree", hserver.NewNotParHandlerFu(c.GetPermissionsSimpleTree))
		ur.GET("/enabled", hserver.NewNotParHandlerFu(c.GetAllEnabled))
	}
}

// AddPermissions 添加权限
// @Summary 添加权限
// @Description 添加权限
// @Tags 系统权限
// @ID AddPermissions
// @Param req body commands.CreatePermissionsCommand true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions [post]
func (c *SysPermissionsController) AddPermissions(ctx context.Context, params *commands.CreatePermissionsCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleCreate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// PermissionsList 获取权限列表
// @Summary 获取权限列表
// @Description 获取权限列表
// @Tags 系统权限
// @ID PermissionsList
// @Param req query queries.ListPermissionsQuery true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.PermissionsDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions [get]
func (c *SysPermissionsController) PermissionsList(ctx context.Context, params *queries.ListPermissionsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdatePermissions 更新权限
// @Summary 更新权限
// @Description 更新权限信息，包括基本信息和资源关联
// @Tags 系统权限
// @ID UpdatePermissions
// @Accept json
// @Produce json
// @Param req body commands.UpdatePermissionsCommand true "权限更新信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions [put]
func (c *SysPermissionsController) UpdatePermissions(ctx context.Context, params *commands.UpdatePermissionsCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeletePermissions 删除权限
// @Summary 删除权限
// @Description 删除指定ID的权限
// @Tags 系统权限
// @ID DeletePermissions
// @Accept json
// @Produce json
// @Param id path int64 true "权限ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions/{id} [delete]
func (c *SysPermissionsController) DeletePermissions(ctx context.Context, params *models.IntIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleDelete(ctx, params.Id)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetPermissionsTree 获取权限树
// @Summary 获取权限树
// @Description 获取权限树形结构
// @Tags 系统权限
// @ID GetPermissionsTree
// @Accept json
// @Produce json
// @Param req query queries.GetPermissionsTreeQuery true "权限树查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.PermissionsDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions/tree [get]
func (c *SysPermissionsController) GetPermissionsTree(ctx context.Context, params *queries.GetPermissionsTreeQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetTree(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetPermissionsSimpleTree 获取简化的权限树结构
// @Summary 获取简化的权限树结构
// @Description 获取简化的权限树结构
// @Tags 系统权限
// @ID GetPermissionsSimpleTree
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=dto.PermissionsTreeResult}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions/simple/tree [get]
func (c *SysPermissionsController) GetPermissionsSimpleTree(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetPermissionsTree(ctx)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetAllEnabled 获取所有启用状态的权限
// @Summary 获取所有启用状态的权限
// @Description 获取系统中所有启用状态的权限列表
// @Tags 系统权限
// @ID GetAllEnabledPermissions
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=[]dto.PermissionsDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/permissions/enabled [get]
func (c *SysPermissionsController) GetAllEnabled(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetAllEnabled(ctx)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
