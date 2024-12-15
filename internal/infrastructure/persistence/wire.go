package persistence

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/data"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.ProviderSet,
	repository.ProviderSet,
)
