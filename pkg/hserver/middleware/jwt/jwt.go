package jwt

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"

	"net/http"
	"strings"
)

// Handler 校验的处理器
func Handler(tokenizer token.IToken) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenEmpty)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenEmpty, constant.RespData: utils.H{}})
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenEmpty)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenEmpty, constant.RespData: utils.H{}})
			c.Abort()
			return
		}

		var accessToken token.AccessToken
		if err := tokenizer.Verify(parts[1], &accessToken); err != nil {
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonTokenVerifyFail)
			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.ReasonTokenVerifyFail, constant.RespData: utils.H{}})
			c.Abort()
			return
		}
		accessToken.AccessToken = parts[1]
		ctx = actx.Store(ctx, accessToken)
		// 将身份信息缓存到Context
		c.Set(constant.KeyAccessToken, accessToken)
		c.Next(ctx)
	}
}

// AdminChickHandler 管理校验
//func AdminChickHandler() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		role := actx.GetRole(ctx)
//		if role == "" {
//			i18Mag := hertzI18n.MustGetMessage(ctx, constant.IsNotAdminAccount)
//			c.JSON(http.StatusOK, utils.H{constant.RespCode: 401, constant.RespMsg: i18Mag, constant.RespReason: constant.IsNotAdminAccount, constant.RespData: utils.H{}})
//			c.Abort()
//			return
//		}
//		if role == constant.RoleSuperAdmin {
//			ctx = actx.BuildIgnoreTenantCtx(ctx)
//		}
//		c.Next(ctx)
//	}
//}
