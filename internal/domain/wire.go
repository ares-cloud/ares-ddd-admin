package domain

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewPermissionService,
	service.NewRoleService,
	service.NewUserService,
)
