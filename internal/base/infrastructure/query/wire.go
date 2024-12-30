package query

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	impl.NewUserQueryService,
	cache.NewCacheEventHandler,
	cache.NewUserQueryCache,
	wire.Bind(new(UserQueryService), new(*cache.UserQueryCache)),
)
