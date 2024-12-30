package mapper

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserMapper,
	NewRoleMapper,
	NewPermissionsMapper,
)
