package query

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	impl.NewUserQueryService,
	impl.NewTenantQueryService,
	impl.NewRoleQueryService,
	impl.NewDepartmentQueryService,
	cache.NewCacheEventHandler,
	cache.NewUserQueryCache,
	cache.NewTenantQueryCache,
	cache.NewRoleQueryCache,
	cache.NewDepartmentQueryCache,
	wire.Bind(new(IUserQueryService), new(*cache.UserQueryCache)),
	wire.Bind(new(ITenantQueryService), new(*cache.TenantQueryCache)),
	wire.Bind(new(IRoleQueryService), new(*cache.RoleQueryCache)),
	wire.Bind(new(IDepartmentQueryService), new(*cache.DepartmentQueryCache)),
)
