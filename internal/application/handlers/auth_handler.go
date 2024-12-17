package handlers

import (
	"context"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/captcha"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
)

type AuthHandler struct {
	authRepo repository.IAuthRepository
}

func NewAuthHandler(authRepo repository.IAuthRepository) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
	}
}

// HandleLogin 处理登录请求
func (h *AuthHandler) HandleLogin(ctx context.Context, cmd commands.LoginCommand, tk token.IToken) (*dto.AuthDto, herrors.Herr) {
	// 验证验证码
	valid, err := h.authRepo.ValidateCaptcha(ctx, cmd.CaptchaKey, cmd.CaptchaCode)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查找用户认证信息
	auth, err := h.authRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 执行登录
	if err := auth.Login(cmd.Password, valid); err != nil {
		return nil, err
	}

	// 获取角色编码列表
	roleCodes := make([]string, 0, len(auth.User.Roles))
	for _, role := range auth.User.Roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成token
	tokenData, err := tk.GenerateToken(auth.User.ID, &token.AccessToken{
		UserId:   auth.User.ID,
		TenantId: auth.User.TenantID,
		Roles:    roleCodes,
		Platform: cmd.Platform,
	})
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	return dto.ToAuthDto(tokenData), nil
}

// HandleRefreshToken 处理刷新token请求
func (h *AuthHandler) HandleRefreshToken(ctx context.Context, cmd commands.RefreshTokenCommand, tk token.IToken) (*dto.AuthDto, herrors.Herr) {
	// 解析token获取用户信息
	accessToken := token.AccessToken{}
	err := tk.Verify(cmd.Token, &accessToken)
	if err != nil {
		return nil, herrors.NewErr(err)
	}
	// 生成新token
	tokenData, err := tk.GenerateToken(accessToken.UserId, &token.AccessToken{
		UserId:   accessToken.UserId,
		TenantId: accessToken.TenantId,
		Roles:    accessToken.Roles,
		Platform: accessToken.Platform,
	})
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	return dto.ToAuthDto(tokenData), nil
}

// HandleGetCaptcha 处理获取验证码请求
func (h *AuthHandler) HandleGetCaptcha(ctx context.Context, query queries.GetCaptchaQuery) (*dto.CaptchaDto, herrors.Herr) {
	// 生成验证码
	id, image, code, err := captcha.GetMathCaptcha(query.Width, query.Height)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 保存验证码
	if err := h.authRepo.SaveCaptcha(ctx, id, code, 5*time.Minute); err != nil {
		return nil, herrors.NewErr(err)
	}

	return &dto.CaptchaDto{
		Key:   id,
		Image: image,
	}, nil
}
