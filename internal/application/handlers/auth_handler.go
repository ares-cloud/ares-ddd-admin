package handlers

import (
	"context"
	"errors"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/ipcity"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/captcha"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
)

type AuthHandler struct {
	authRepo repository.IAuthRepository
	uds      *service.UserService
	llr      repository.ILoginLogRepository
}

func NewAuthHandler(authRepo repository.IAuthRepository, uds *service.UserService, llr repository.ILoginLogRepository) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
		uds:      uds,
		llr:      llr,
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

	// 验证登录类型
	if cmd.LoginType == commands.LoginTypeAdmin && !auth.User.IsAdmin() {
		err = errors.New("非管理员用户不能登录管理端")
		go h.recordLoginLog(ctx, auth.User, cmd, err)
		return nil, herrors.NewBadReqError(err.Error())
	}

	// 执行登录
	if err1 := auth.Login(cmd.Password, valid); herrors.HaveError(err1) {
		return nil, err1
	}
	ctx = actx.WithTenantId(ctx, auth.User.TenantID)
	roles, e := h.uds.GetUserRoles(ctx, auth.User)
	if e != nil {
		hlog.CtxErrorf(ctx, "get user roles failed: %v", e)
		return nil, herrors.QueryFail(e)
	}

	// 生成token
	tokenData, err := tk.GenerateToken(auth.User.ID, &token.AccessToken{
		UserId:   auth.User.ID,
		TenantId: auth.User.TenantID,
		Roles:    roles,
		Platform: cmd.Platform,
	})
	if err != nil {
		return nil, herrors.NewErr(err)
	}
	// 记录登录失败日志
	go h.recordLoginLog(ctx, auth.User, cmd, nil)
	return dto.ToAuthDto(tokenData), nil
}

// recordLoginLog 记录登录日志
func (h *AuthHandler) recordLoginLog(ctx context.Context, user *model.User, cmd commands.LoginCommand, loginErr error) {
	var loginLog *model.LoginLog
	if user != nil {
		loginLog = model.NewLoginLog(user.ID, user.Username, user.TenantID, model.LoginType(cmd.LoginType))
	} else {
		loginLog = model.NewLoginLog("", cmd.Username, "", model.LoginType(cmd.LoginType))
	}
	address := actx.GetIpAddress(ctx)

	//获取登录地
	locationBaiDu, loginErr := ipcity.GetGetLocationBaiDu(address)
	if loginErr != nil {
		hlog.CtxErrorf(ctx, "get location bai du failed: %v", loginErr)
	}
	if loginErr != nil {
		loginLog.SetLoginStatus(2, loginErr.Error())
	}
	dev := actx.GetDeviceId(ctx)
	os := actx.GetDeviceName(ctx)
	bro := actx.GetUserAgent(ctx)
	loginLog.SetLoginInfo(address, locationBaiDu, dev, os, bro)
	loginErr = h.llr.Create(ctx, loginLog)
	if loginErr != nil {
		hlog.CtxErrorf(ctx, "create login failed: %v", loginErr)
	}
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
	id, image, code, err := captcha.GetDigitCaptcha(query.Width, query.Height, 3)
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
