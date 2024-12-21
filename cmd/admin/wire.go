//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/ares-cloud/ares-ddd-admin/cmd/admin/server"
	"github.com/ares-cloud/ares-ddd-admin/internal/base"
	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database"
	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configs.Bootstrap, *configs.Data, *configs.StorageConfig) (*app, func(), error) {
	panic(wire.Build(database.ProviderSet, base.ProviderSet, monitoring.ProviderSet, storage.ProviderSet, server.ProviderSet, newApp))
}
