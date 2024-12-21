package rest

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/device"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type AuthController struct {
	authHandler *handlers.AuthHandler
	t           token.IToken
}

func NewAuthController(authHandler *handlers.AuthHandler) *AuthController {
	return &AuthController{
		authHandler: authHandler,
	}
}

func (c *AuthController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	c.t = t
	v1 := g.Group("/v1")
	auth := v1.Group("/auth")
	{
		auth.POST("/login", device.Handler(), hserver.NewHandlerFu[commands.LoginCommand](c.Login))
		auth.POST("/refresh", hserver.NewHandlerFu[commands.RefreshTokenCommand](c.RefreshToken))
		auth.GET("/captcha", hserver.NewHandlerFu[queries.GetCaptchaQuery](c.GetCaptcha))
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @ID Login
// @Accept json
// @Produce json
// @Param data body commands.LoginCommand true "登录参数"
// @Success 200 {object} base_info.Success{data=dto.AuthDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "认证失败"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/auth/login [post]
func (c *AuthController) Login(ctx context.Context, req *commands.LoginCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.authHandler.HandleLogin(ctx, *req, c.t)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用旧令牌获取新的访问令牌
// @Tags 认证
// @ID RefreshToken
// @Accept json
// @Produce json
// @Param req body commands.RefreshTokenCommand true "刷新令牌请求"
// @Success 200 {object} base_info.Success{data=dto.AuthDto}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "认证失败"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx context.Context, params *commands.RefreshTokenCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.authHandler.HandleRefreshToken(ctx, *params, c.t)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetCaptcha 获取验证码
// @Summary 获取验证码
// @Description 获取图形验证码
// @Tags 认证
// @ID GetCaptcha
// @Accept json
// @Produce json
// @Param width query int false "验证码宽度"
// @Param height query int false "验证码高度"
// @Success 200 {object} base_info.Success{data=dto.CaptchaDto}
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/auth/captcha [get]
func (c *AuthController) GetCaptcha(ctx context.Context, params *queries.GetCaptchaQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.authHandler.HandleGetCaptcha(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}
