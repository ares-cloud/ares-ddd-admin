//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/ares-cloud/ares-ddd-admin/cmd/admin/server"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/base"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database"
	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configs.Bootstrap, *configs.Data) (*app, func(), error) {
	panic(wire.Build(database.ProviderSet, base.ProviderSet, domain.ProviderSet, persistence.ProviderSet, handlers.ProviderSet, server.ProviderSet, newApp))
}
