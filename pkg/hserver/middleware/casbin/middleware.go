package casbin

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"net/http"
)

func Handler(enforcer *Enforcer) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文获取用户信息
		roles := actx.GetRoles(ctx)
		tenantID := actx.GetTenantId(ctx)

		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		//对每个角色进行权限检查
		hasPermission := false
		if actx.IsSuperAdmin(ctx) {
			hasPermission = true
		} else {
			for _, role := range roles {
				allowed, err := enforcer.Enforce(role, tenantID, method, path)
				if err != nil {
					hlog.CtxErrorf(ctx, "casbin enforce error: %v", err)
					continue
				}
				if allowed {
					hasPermission = true
					break
				}
			}
		}
		if !hasPermission {
			hlog.CtxInfof(ctx, "permission denied for user %s, path: %s, method: %s", actx.GetUserId(ctx), path, method)
			i18Mag := hertzI18n.MustGetMessage(ctx, constant.ReasonNoAccess)
			c.JSON(http.StatusOK, utils.H{
				constant.RespCode:   http.StatusForbidden,
				constant.RespMsg:    i18Mag,
				constant.RespReason: constant.ReasonNoAccess,
				constant.RespData:   utils.H{},
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
