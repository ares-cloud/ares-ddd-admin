package converter

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserConverter,
	NewDepartmentConverter,
	NewDataPermissionConverter,
	NewPermissionsConverter,
	NewRoleConverter,
	NewTenantConverter,
)
