package query

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/cache/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query/impl"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	impl.NewUserQueryService,
	impl.NewRoleQueryService,
	impl.NewPermissionsQueryService,
	impl.NewTenantQueryService,
	impl.NewDepartmentQueryService,
	impl.NewDataPermissionQueryService,
	impl.NewOperationLogQueryService,
	impl.NewLoginLogQueryService,

	cache.NewUserQueryCache,
	cache.NewRoleQueryCache,
	cache.NewPermissionsQueryCache,
	cache.NewTenantQueryCache,
	cache.NewDepartmentQueryCache,
	cache.NewDataPermissionQueryCache,

	handlers.NewCacheEventHandler,
	// 绑定接口到实现
	wire.Bind(new(IUserQueryService), new(*cache.UserQueryCache)),
	wire.Bind(new(ITenantQueryService), new(*cache.TenantQueryCache)),
	wire.Bind(new(IRoleQueryService), new(*cache.RoleQueryCache)),
	wire.Bind(new(IDepartmentQueryService), new(*cache.DepartmentQueryCache)),
	wire.Bind(new(IPermissionsQuery), new(*cache.PermissionsQueryCache)),
	wire.Bind(new(IDataPermissionQuery), new(*cache.DataPermissionQueryCache)),
	wire.Bind(new(IOperationLogQuery), new(*impl.OperationLogQueryService)),
	wire.Bind(new(ILoginLogQuery), new(*impl.LoginLogQueryService)),
)
