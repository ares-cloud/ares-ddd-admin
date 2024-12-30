package query

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	impl.NewUserQueryService,
	impl.NewTenantQueryService,
	cache.NewCacheEventHandler,
	cache.NewUserQueryCache,
	cache.NewTenantQueryCache,
	wire.Bind(new(UserQueryService), new(*cache.UserQueryCache)),
	wire.Bind(new(TenantQueryService), new(*cache.TenantQueryCache)),
)
