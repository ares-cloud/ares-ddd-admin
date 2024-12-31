package infrastructure

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/base"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	base.ProviderSet,
	converter.ProviderSet,
	handlers.ProviderSet,
	persistence.ProviderSet,
	query.ProviderSet,
)
