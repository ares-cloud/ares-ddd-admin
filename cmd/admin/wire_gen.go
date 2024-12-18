// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/service"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/auth/casbin"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/configs"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/database"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/data"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/interfaces/rest"
	"github.com/ares-cloud/ares-ddd-admin/internal/interfaces/server/admin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/snowflake_id"
)

import (
	_ "github.com/ares-cloud/ares-ddd-admin/docs/admin"
)

// Injectors from wire.go:

// wireApp init application.
func wireApp(bootstrap *configs.Bootstrap, configsData *configs.Data) (*app, func(), error) {
	redisClient, err := database.NewHdbClient(configsData)
	if err != nil {
		return nil, nil, err
	}
	iIdGenerate := snowflake_id.NewSnowIdGen()
	iDataBase, cleanup, err := database.NewDataBase(iIdGenerate, configsData)
	if err != nil {
		return nil, nil, err
	}
	iSysRoleRepo := data.NewSysRoleRepo(iDataBase)
	iPermissionsRepo := data.NewSysMenuRepo(iDataBase)
	iRoleRepository := repository.NewRoleRepository(iSysRoleRepo, iPermissionsRepo)
	iPermissionsRepository := repository.NewPermissionsRepository(iPermissionsRepo)
	roleCommandHandler := handlers.NewRoleCommandHandler(iRoleRepository, iPermissionsRepository)
	roleQueryHandler := handlers.NewRoleQueryHandler(iRoleRepository)
	casbinIPermissionsRepository := casbin.NewRepositoryImpl(iSysRoleRepo, iPermissionsRepo)
	enforcer, err := admin.NewCasBinEnforcer(redisClient, casbinIPermissionsRepository)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sysRoleController := rest.NewSysRoleController(roleCommandHandler, roleQueryHandler, enforcer)
	iSysUserRepo := data.NewSysUserRepo(iDataBase)
	iUserRepository := repository.NewUserRepository(iSysUserRepo, iSysRoleRepo)
	userCommandHandler := handlers.NewUserCommandHandler(iUserRepository, iRoleRepository)
	iSysTenantRepo := data.NewSysTenantRepo(iDataBase)
	iTenantRepository := repository.NewTenantRepository(iSysTenantRepo, iSysUserRepo)
	userService := service.NewUserService(iUserRepository, iPermissionsRepository, iRoleRepository, iTenantRepository)
	userQueryHandler := handlers.NewUserQueryHandler(iUserRepository, userService, iPermissionsRepository)
	sysUserController := rest.NewSysUserController(userCommandHandler, userQueryHandler, enforcer)
	tenantCommandHandler := handlers.NewTenantCommandHandler(iTenantRepository)
	tenantQueryHandler := handlers.NewTenantQueryHandler(iTenantRepository, iPermissionsRepository)
	sysTenantController := rest.NewSysTenantController(tenantCommandHandler, tenantQueryHandler, enforcer)
	permissionsCommandHandler := handlers.NewPermissionsCommandHandler(iPermissionsRepository, enforcer)
	permissionsQueryHandler := handlers.NewPermissionsQueryHandler(iPermissionsRepository)
	sysPermissionsController := rest.NewSysPermissionsController(permissionsCommandHandler, permissionsQueryHandler, enforcer)
	iAuthRepository := repository.NewAuthRepository(iUserRepository, redisClient)
	authHandler := handlers.NewAuthHandler(iAuthRepository, userService)
	authController := rest.NewAuthController(authHandler)
	serve := admin.NewServer(bootstrap, redisClient, sysRoleController, sysUserController, sysTenantController, sysPermissionsController, authController)
	mainApp := newApp(serve)
	return mainApp, func() {
		cleanup()
	}, nil
}
