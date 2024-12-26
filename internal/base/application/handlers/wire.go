package handlers

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserCommandHandler,
	NewUserQueryHandler,
	NewRoleCommandHandler,
	NewRoleQueryHandler,
	NewPermissionsCommandHandler,
	NewPermissionsQueryHandler,
	NewTenantCommandHandler,
	NewTenantQueryHandler,
	NewAuthHandler,
	NewLoginLogQueryHandler,
	NewOperationLogQueryHandler,
	NewDepartmentCommandHandler,
	NewDepartmentQueryHandler,
)
