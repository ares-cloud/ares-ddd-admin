package rest

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"

	_ "github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	_ "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/base_info"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type SysRoleController struct {
	cmdHandel   *handlers.RoleCommandHandler
	queryHandel *handlers.RoleQueryHandler
	ef          *casbin.Enforcer
	modeNma     string
}

func NewSysRoleController(cmdHandel *handlers.RoleCommandHandler, queryHandel *handlers.RoleQueryHandler, ef *casbin.Enforcer) *SysRoleController {
	return &SysRoleController{
		cmdHandel:   cmdHandel,
		queryHandel: queryHandel,
		ef:          ef,
		modeNma:     "角色",
	}
}

func (c *SysRoleController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	ur := v1.Group("/sys/role", jwt.Handler(t))
	{
		ur.POST("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "新增",
		}), hserver.NewHandlerFu[commands.CreateRoleCommand](c.AddRole))
		ur.GET("", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListRolesQuery](c.RoleList))
		ur.PUT("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "修改",
		}), casbin.Handler(c.ef), hserver.NewHandlerFu[commands.UpdateRoleCommand](c.UpdateRole))
		ur.DELETE("/:id", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.modeNma,
			Action:      "删除",
		}), casbin.Handler(c.ef), hserver.NewHandlerFu[models.IntIdReq](c.DeleteRole))
		ur.GET("/:id", casbin.Handler(c.ef), hserver.NewHandlerFu[models.IntIdReq](c.GetDetails))
		ur.GET("/enabled", hserver.NewNotParHandlerFu(c.GetAllEnabled))
		ur.GET("/data-permission", hserver.NewNotParHandlerFu(c.GetAllDataPermission))
	}
}

// AddRole 添加角色
// @Summary 添加角色
// @Description 添加角色
// @Tags 系统角色
// @ID AddRole
// @Param req body commands.CreateRoleCommand true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role [post]
func (c *SysRoleController) AddRole(ctx context.Context, params *commands.CreateRoleCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleCreate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// RoleList 获取角色
// @Summary 获取角色
// @Description 获取角色
// @Tags 系统角色
// @ID RoleList
// @Param req query queries.ListRolesQuery true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.RoleDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role [get]
func (c *SysRoleController) RoleList(ctx context.Context, params *queries.ListRolesQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息，包括基本信息和权限关联
// @Tags 系统角色
// @ID UpdateRole
// @Accept json
// @Produce json
// @Param req body commands.UpdateRoleCommand true "角色更新信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role [put]
func (c *SysRoleController) UpdateRole(ctx context.Context, params *commands.UpdateRoleCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定ID的角色
// @Tags 系统角色
// @ID DeleteRole
// @Accept json
// @Produce json
// @Param id path int64 true "角色ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role/{id} [delete]
func (c *SysRoleController) DeleteRole(ctx context.Context, params *models.IntIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleDelete(ctx, commands.DeleteRoleCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetDetails 根据id获取角色
// @Summary 根据id获取角色
// @Description 根据id获取角色
// @Tags 系统角色
// @ID GetRoleDetails
// @Accept json
// @Produce json
// @Param id path int64 true "角色ID"
// @Success 200 {object} base_info.Success{data=dto.RoleDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role/:id [get]
func (c *SysRoleController) GetDetails(ctx context.Context, params *models.IntIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGet(ctx, queries.GetRoleQuery{Id: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetAllEnabled 获取所有启用状态的角色
// @Summary 获取所有启用状态的角色
// @Description 获取系统中所有启用状态的角色列表
// @Tags 系统角色
// @ID GetAllEnabledRoles
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=[]dto.RoleDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/role/enabled [get]
func (c *SysRoleController) GetAllEnabled(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetAllEnabled(ctx)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetAllDataPermission 获取所有数据权限角色
// @Summary 获取所有数据权限角色
// @Description 获取系统中所有数据权限类型的角色列表
// @Tags 系统角色
// @ID GetAllDataPermissionRoles
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=[]dto.RoleDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/role/data-permission [get]
func (c *SysRoleController) GetAllDataPermission(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetAllDataPermission(ctx)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
