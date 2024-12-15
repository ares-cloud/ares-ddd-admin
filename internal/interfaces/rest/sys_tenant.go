package rest

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	_ "github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	_ "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/base_info"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type SysTenantController struct {
	cmdHandel   *handlers.TenantCommandHandler
	queryHandel *handlers.TenantQueryHandler
}

func NewSysTenantController(cmdHandel *handlers.TenantCommandHandler, queryHandel *handlers.TenantQueryHandler) *SysTenantController {
	return &SysTenantController{
		cmdHandel:   cmdHandel,
		queryHandel: queryHandel,
	}
}

func (c *SysTenantController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	ur := v1.Group("/sys/tenant")
	{
		ur.POST("", hserver.NewHandlerFu[commands.CreateTenantCommand](c.AddTenant))
		ur.GET("", hserver.NewHandlerFu[queries.ListTenantsQuery](c.TenantList))
		ur.PUT("", hserver.NewHandlerFu[commands.UpdateTenantCommand](c.UpdateTenant))
		ur.DELETE("/:id", hserver.NewHandlerFu[models.StringIdReq](c.DeleteTenant))
		ur.POST("/permissions", hserver.NewHandlerFu[commands.AssignTenantPermissionsCommand](c.AssignPermissions))
		ur.GET("/permissions", hserver.NewHandlerFu[queries.GetTenantPermissionsQuery](c.GetPermissions))
	}
}

// AddTenant 添加租户
// @Summary 添加租户
// @Description 添加租户，包括管理员用户信息
// @Tags 系统租户
// @ID AddTenant
// @Accept json
// @Produce json
// @Param req body commands.CreateTenantCommand true "租户创建信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant [post]
func (c *SysTenantController) AddTenant(ctx context.Context, params *commands.CreateTenantCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleCreate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// TenantList 获取租户列表
// @Summary 获取租户列表
// @Description 获取租户列表，支持分页和条件查询
// @Tags 系统租户
// @ID TenantList
// @Accept json
// @Produce json
// @Param req query queries.ListTenantsQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.TenantDto]}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant [get]
func (c *SysTenantController) TenantList(ctx context.Context, params *queries.ListTenantsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdateTenant 更新租户
// @Summary 更新租户
// @Description 更新租户信息
// @Tags 系统租户
// @ID UpdateTenant
// @Accept json
// @Produce json
// @Param req body commands.UpdateTenantCommand true "租户更新信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant [put]
func (c *SysTenantController) UpdateTenant(ctx context.Context, params *commands.UpdateTenantCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteTenant 删除租户
// @Summary 删除租户
// @Description 删除指定ID的租户
// @Tags 系统租户
// @ID DeleteTenant
// @Accept json
// @Produce json
// @Param id path string true "租户ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant/{id} [delete]
func (c *SysTenantController) DeleteTenant(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleDelete(ctx, commands.DeleteTenantCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// AssignPermissions 分配权限给租户
// @Summary 分配权限给租户
// @Description 为指定租户分配权限
// @Tags 系统租户
// @ID AssignTenantPermissions
// @Accept json
// @Produce json
// @Param req body commands.AssignTenantPermissionsCommand true "权限分配信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant/permissions [post]
func (c *SysTenantController) AssignPermissions(ctx context.Context, params *commands.AssignTenantPermissionsCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleAssignPermissions(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetPermissions 获取租户权限
// @Summary 获取租户权限
// @Description 获取指定租户的权限列表
// @Tags 系统租户
// @ID GetTenantPermissions
// @Accept json
// @Produce json
// @Param tenant_id query string true "租户ID"
// @Success 200 {object} base_info.Success{data=[]dto.PermissionsDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/tenant/permissions [get]
func (c *SysTenantController) GetPermissions(ctx context.Context, params *queries.GetTenantPermissionsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleGetPermissions(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
