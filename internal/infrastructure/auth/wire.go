package auth

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/auth/casbin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	casbin.NewRepositoryImpl,
)
