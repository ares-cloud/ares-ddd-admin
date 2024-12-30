package events

import (
	"github.com/ares-cloud/ares-ddd-admin/pkg/events"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	events.NewEventBus,
)
