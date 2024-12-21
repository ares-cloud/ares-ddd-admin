package data

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSysUserRepo,
	NewSysMenuRepo,
	NewSysRoleRepo,
	NewSysTenantRepo,
)
