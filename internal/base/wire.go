package base

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/base"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/interfaces/rest"
	"github.com/google/wire"
)

// ProviderSet is monitoring providers.
var ProviderSet = wire.NewSet(
	domain.ProviderSet,
	handlers.ProviderSet,
	persistence.ProviderSet,
	base.ProviderSet,
	rest.NewSysRoleController,
	rest.NewSysUserController,
	rest.NewSysTenantController,
	rest.NewSysPermissionsController,
	rest.NewAuthController,
	rest.NewLoginLogController,
	rest.NewOperationLogController,
	rest.NewDepartmentController,
	rest.NewDataPermissionController,
	NewBaseServer,
)
