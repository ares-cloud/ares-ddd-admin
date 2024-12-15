package tenant

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/cloudwego/hertz/pkg/app"
)

// IgnoreTenantHandler 租户处理
func IgnoreTenantHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ctx = actx.BuildIgnoreTenantCtx(ctx)
		c.Next(ctx)
	}
}
