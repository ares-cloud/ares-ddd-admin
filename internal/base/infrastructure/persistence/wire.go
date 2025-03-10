package persistence

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/data"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.ProviderSet,
	mapper.ProviderSet,
	repository.ProviderSet,
)
