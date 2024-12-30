package application

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	handlers.NewUserCommandHandler,
	handlers.NewUserQueryHandler,
)
