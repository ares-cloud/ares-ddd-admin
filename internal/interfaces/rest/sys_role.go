package rest

import (
	"context"
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

type SysRoleController struct {
	cmdHandel   *handlers.RoleCommandHandler
	queryHandel *handlers.RoleQueryHandler
}

func NewSysRoleController(cmdHandel *handlers.RoleCommandHandler, queryHandel *handlers.RoleQueryHandler) *SysRoleController {
	return &SysRoleController{
		cmdHandel:   cmdHandel,
		queryHandel: queryHandel,
	}
}

func (c *SysRoleController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	ur := v1.Group("/sys/role")
	{
		ur.POST("", hserver.NewHandlerFu[commands.CreateRoleCommand](c.AddRole))
		ur.GET("", hserver.NewHandlerFu[queries.ListRolesQuery](c.RoleList))
		ur.PUT("", hserver.NewHandlerFu[commands.UpdateRoleCommand](c.UpdateRole))
		ur.DELETE("/:id", hserver.NewHandlerFu[models.IntIdReq](c.DeleteRole))
		ur.GET("/:id", hserver.NewHandlerFu[models.IntIdReq](c.GetDetails))
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
// @ID GetDetails
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
