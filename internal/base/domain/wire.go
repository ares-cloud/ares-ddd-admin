package domain

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAuthService,
	service.NewRoleCommandService,
	service.NewPermissionService,
	service.NewTenantCommandService,
	service.NewDepartmentService,
	service.NewUserCommandService,
	service.NewDataPermissionService,
)
