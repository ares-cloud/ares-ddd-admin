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

type SysUserController struct {
	cmdHandel   *handlers.UserCommandHandler
	queryHandel *handlers.UserQueryHandler
}

func NewSysUserController(cmdHandel *handlers.UserCommandHandler, queryHandel *handlers.UserQueryHandler) *SysUserController {
	return &SysUserController{
		cmdHandel:   cmdHandel,
		queryHandel: queryHandel,
	}
}

func (c *SysUserController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	ur := v1.Group("/sys/user")
	{
		ur.POST("", hserver.NewHandlerFu[commands.CreateUserCommand](c.AddUser))
		ur.GET("", hserver.NewHandlerFu[queries.ListUsersQuery](c.UserList))
		ur.PUT("", hserver.NewHandlerFu[commands.UpdateUserCommand](c.UpdateUser))
		ur.DELETE("/:id", hserver.NewHandlerFu[models.StringIdReq](c.DeleteUser))
		ur.PUT("/status", hserver.NewHandlerFu[commands.UpdateUserStatusCommand](c.UpdateUserStatus))
	}
}

// AddUser 添加用户
// @Summary 添加用户
// @Description 添加用户
// @Tags 系统用户
// @ID AddUser
// @Param req body commands.CreateUserCommand true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/user [post]
func (c *SysUserController) AddUser(ctx context.Context, params *commands.CreateUserCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleCreate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// UserList 获取用户
// @Summary 获取用户
// @Description 获取用户
// @Tags 系统用户
// @ID UserList
// @Param req query queries.ListUsersQuery true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.UserDto}
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /v1/sys/user [get]
func (c *SysUserController) UserList(ctx context.Context, params *queries.ListUsersQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandel.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息，包括基本信息和角色关联
// @Tags 系统用户
// @ID UpdateUser
// @Accept json
// @Produce json
// @Param req body commands.UpdateUserCommand true "用户更新信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/sys/user [put]
func (c *SysUserController) UpdateUser(ctx context.Context, params *commands.UpdateUserCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定ID的用户
// @Tags 系统用户
// @ID DeleteUser
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/sys/user/{id} [delete]
func (c *SysUserController) DeleteUser(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleDelete(ctx, commands.DeleteUserCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户的启用/禁用状态
// @Tags 系统用户
// @ID UpdateUserStatus
// @Accept json
// @Produce json
// @Param req body commands.UpdateUserStatusCommand true "用户状态更新信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/sys/user/status [put]
func (c *SysUserController) UpdateUserStatus(ctx context.Context, params *commands.UpdateUserStatusCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandel.HandleUpdateStatus(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}
