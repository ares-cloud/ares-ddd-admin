package device

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/cloudwego/hertz/pkg/app"
)

// Handler 设备处理器
func Handler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ipAddress := c.ClientIP()
		userAgent := c.UserAgent()
		deviceId := c.Request.Header.Get("Device")
		deviceName := c.Request.Header.Get("OS")
		ctx = actx.WithIpAddress(ctx, ipAddress)
		ctx = actx.WithDeviceId(ctx, deviceId)
		ctx = actx.WithDeviceName(ctx, deviceName)
		ctx = actx.WithUserAgent(ctx, string(userAgent))
		c.Next(ctx)
	}
}
