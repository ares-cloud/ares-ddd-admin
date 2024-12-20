package repository

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewRoleRepository,
	NewPermissionsRepository,
	NewTenantRepository,
	NewAuthRepository,
	NewLoginLogRepository,
)
