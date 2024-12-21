package base

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/base/casbin"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/base/oplog"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	casbin.NewRepositoryImpl,
	oplog.NewDbOperationLogWriter,
)
