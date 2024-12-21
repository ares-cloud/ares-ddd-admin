package monitoring

import (
	"github.com/google/wire"

	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/monitoring/interfaces/rest"
)

// ProviderSet is monitoring providers.
var ProviderSet = wire.NewSet(
	service.NewMetricsService,
	handlers.NewMetricsQueryHandler,
	rest.NewMetricsController,
	NewServer,
)
