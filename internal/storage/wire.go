package storage

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/data"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/query/impl"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/storage"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/google/wire"

	"github.com/ares-cloud/ares-ddd-admin/internal/storage/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/cleaner"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/storage/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
)

// ProviderSet is storage providers.
var ProviderSet = wire.NewSet(
	data.NewStorageRepo,
	repository.NewStorageRepository,
	service.NewStorageService,
	impl.NewStorageQueryService,
	handlers.NewStorageQueryHandler,
	handlers.NewStorageCommandHandler,
	rest.NewStorageController,
	cleaner.NewRecycleCleaner,
	storage.NewStorageFactory,
	NewServer,
)

// Server represents storage server.
type Server struct {
	controller *rest.StorageController
	cleaner    *cleaner.RecycleCleaner
}

// NewServer creates a new storage server.
func NewServer(controller *rest.StorageController, cleaner *cleaner.RecycleCleaner) (*Server, func(), error) {
	s := &Server{
		controller: controller,
		cleaner:    cleaner,
	}
	cleanup := func() {
		hlog.Info("closing the data resources")
		s.cleaner.Stop()
	}
	return s, cleanup, nil
}

// Init initializes the storage server.
func (s *Server) Init(g *route.RouterGroup, t token.IToken) {
	s.controller.RegisterRouter(g, t)
	s.cleaner.Start() // 启动回收站清理任务
}
