package actx

import (
	"context"
	"fmt"
	"strings"

	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
)

const (
	keyAccessToken = "access_token"
	KeyUserId      = "userId"
	KeyPlatform    = "platform"
	KeyToken       = "token"
	KeyRole        = "role"
	KeyTenantId    = "tenant_id"
	DeviceId       = "deviceId"
	DeviceName     = "deviceName"
	IpAddress      = "ipAddress"
	UserAgent      = "UserAgent"
	IgnoreTenantId = "ignore_tenant_Id"
)

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, KeyUserId, userId)
}

func GetUserId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyUserId))
}
func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, KeyPlatform, platform)
}

func GetPlatform(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyPlatform))
}
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, KeyToken, token)
}

func GetToken(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyToken))
}

func WithRole(ctx context.Context, role []string) context.Context {
	return context.WithValue(ctx, KeyRole, strings.Join(role, ","))
}

func GetRoles(ctx context.Context) []string {
	value := fmt.Sprintf("%v", ctx.Value(KeyRole))
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}
func WithTenantId(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, KeyTenantId, tenantId)
}

func GetTenantId(ctx context.Context) string {
	tenId := fmt.Sprintf("%v", ctx.Value(KeyTenantId))
	if tenId == IgnoreTenantId {
		return ""
	}
	return tenId
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	return context.WithValue(ctx, DeviceId, deviceId)
}

func GetDeviceId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceId))
}
func WithDeviceName(ctx context.Context, deviceName string) context.Context {
	return context.WithValue(ctx, DeviceName, deviceName)
}

func GetDeviceName(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceName))
}

func WithIpAddress(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, IpAddress, addr)
}

func GetIpAddress(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(IpAddress))
}

func WithUserAgent(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, UserAgent, addr)
}

func GetUserAgent(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(UserAgent))
}

func WithIgnoreTenantId(ctx context.Context) context.Context {
	return context.WithValue(ctx, IgnoreTenantId, IgnoreTenantId)
}
func GetIgnoreTenantId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(IgnoreTenantId))
}
func IsIgnoreTenantId(ctx context.Context) bool {
	return GetIgnoreTenantId(ctx) == IgnoreTenantId
}

// BuildIgnoreTenantCtx 构建忽略租户的ctx
func BuildIgnoreTenantCtx(ctx context.Context) context.Context {
	return WithIgnoreTenantId(ctx)
}

func Store(ctx context.Context, accessToken token.AccessToken) context.Context {
	ctx = WithUserId(ctx, accessToken.UserId)
	ctx = WithPlatform(ctx, accessToken.Platform)
	ctx = WithToken(ctx, accessToken.AccessToken)
	ctx = WithRole(ctx, accessToken.Roles)
	ctx = WithTenantId(ctx, accessToken.TenantId)

	return ctx
}
