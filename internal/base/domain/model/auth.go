package model

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/password"
)

var (
	ErrInvalidPassword = herrors.NewBadReqError("PasswordError")
	ErrInvalidCaptcha  = herrors.NewBadReqError("IncorrectVerificationCode")
)

// Auth 认证领域模型
type Auth struct {
	User           *User
	AccessToken    string
	ExpiresIn      int64
	RefreshToken   string
	RefreshExpires int64
	Platform       string // 登录平台
}

// NewAuth 创建认证实例
func NewAuth(user *User, platform string) *Auth {
	return &Auth{
		User:     user,
		Platform: platform,
	}
}

// Login 用户登录
func (a *Auth) Login(plainPwd string, captchaValid bool) herrors.Herr {
	if !captchaValid {
		return ErrInvalidCaptcha
	}
	if !password.CheckPasswordHash(plainPwd, a.User.Password) {
		return ErrInvalidPassword
	}
	return nil
}

// SetTokens 设置访问令牌和刷新令牌
func (a *Auth) SetTokens(accessToken string, expiresIn int64, refreshToken string, refreshExpires int64) {
	a.AccessToken = accessToken
	a.ExpiresIn = expiresIn
	a.RefreshToken = refreshToken
	a.RefreshExpires = refreshExpires
}

// IsTokenExpired 检查令牌是否过期
func (a *Auth) IsTokenExpired() bool {
	return time.Now().Unix() >= a.ExpiresIn
}
