package rest

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	_ "github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type DataPermissionController struct {
	cmdHandler   *handlers.DataPermissionCommandHandler
	queryHandler *handlers.DataPermissionQueryHandler
}

func NewDataPermissionController(
	cmdHandler *handlers.DataPermissionCommandHandler,
	queryHandler *handlers.DataPermissionQueryHandler,
) *DataPermissionController {
	return &DataPermissionController{
		cmdHandler:   cmdHandler,
		queryHandler: queryHandler,
	}
}

func (c *DataPermissionController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	dp := g.Group("/v1/data-permission", jwt.Handler(t))
	{
		dp.POST("/assign", hserver.NewHandlerFu[commands.AssignDataPermissionCommand](c.AssignDataPermission))
		dp.POST("/remove", hserver.NewHandlerFu[commands.RemoveDataPermissionCommand](c.RemoveDataPermission))
		dp.GET("/:id", hserver.NewHandlerFu[models.IntIdReq](c.GetByRoleID))
	}
}

// AssignDataPermission 分配数据权限
// @Summary 分配数据权限
// @Description 为角色分配数据权限范围
// @Tags 数据权限
// @Accept json
// @Produce json
// @Param req body commands.AssignDataPermissionCommand true "分配数据权限参数"
// @Success 200 {object} base_info.Success
// @Router /v1/data-permission/assign [post]
func (c *DataPermissionController) AssignDataPermission(ctx context.Context, req *commands.AssignDataPermissionCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleAssign(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}

// RemoveDataPermission 移除数据权限
// @Summary 移除数据权限
// @Description 移除角色的数据权限配置
// @Tags 数据权限
// @Accept json
// @Produce json
// @Param req body commands.RemoveDataPermissionCommand true "移除数据权限参数"
// @Success 200 {object} base_info.Success
// @Router /v1/data-permission/remove [post]
func (c *DataPermissionController) RemoveDataPermission(ctx context.Context, req *commands.RemoveDataPermissionCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleRemove(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}

// GetByRoleID 获取角色的数据权限
// @Summary 获取角色的数据权限
// @Description 获取指定角色ID的数据权限配置
// @Tags 数据权限
// @Accept json
// @Produce json
// @Param roleId path int64 true "角色ID"
// @Success 200 {object} base_info.Success{data=dto.DataPermissionDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/data-permission/{id} [get]
func (c *DataPermissionController) GetByRoleID(ctx context.Context, req *models.IntIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGetByRoleID(ctx, queries.GetDataPermissionQuery{RoleID: req.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
