//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/auth"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence"
	"github.com/ares-cloud/ares-ddd-admin/internal/interfaces/server/admin"
	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configs.Bootstrap, *configs.Data) (*app, func(), error) {
	panic(wire.Build(database.ProviderSet, auth.ProviderSet, domain.ProviderSet, persistence.ProviderSet, handlers.ProviderSet, admin.ProviderSet, newApp))
}
