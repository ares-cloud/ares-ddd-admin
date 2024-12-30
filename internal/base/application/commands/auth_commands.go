package commands

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/validator"
)

// LoginType 登录类型
type LoginType int8

const (
	LoginTypeAdmin  LoginType = 1 // 管理端登录
	LoginTypeMember LoginType = 2 // 前台用户登录
)

// LoginCommand 登录命令
type LoginCommand struct {
	Username    string    `json:"username" validate:"required" label:"用户名"`
	Password    string    `json:"password" validate:"required" label:"密码"`
	CaptchaKey  string    `json:"captchaKey" validate:"required" label:"验证码Key"`
	CaptchaCode string    `json:"captchaCode" validate:"required" label:"验证码"`
	Platform    string    `json:"platform" validate:"required" label:"登录平台"`
	LoginType   LoginType `json:"login_type" validate:"required" label:"登录类型"`
	IP          string    `json:"ip" label:"登录IP"`
	Location    string    `json:"location" label:"登录地点"`
	UserAgent   string    `json:"user_agent" label:"User-Agent"`
}

func (c *LoginCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	Token string `json:"token" validate:"required" label:"令牌"`
}

func (c *RefreshTokenCommand) Validate() herrors.Herr {
	return validator.Validate(c)
}
