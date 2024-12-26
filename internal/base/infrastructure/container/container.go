package container

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/command"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/query"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/service"
)

func (c *Container) GetDataPermissionCommandHandler() *command.DataPermissionCommandHandler {
	return &command.DataPermissionCommandHandler{
		tx:              c.GetTransaction(),
		permissionRepo:  c.GetPermissionRepository(),
		roleRepo:        c.GetRoleRepository(),
		deptService:     c.GetDepartmentService(),
		userDeptService: c.GetUserDeptService(),
	}
}

func (c *Container) GetDataPermissionQueryHandler() *query.DataPermissionQueryHandler {
	return &query.DataPermissionQueryHandler{
		permissionService: c.GetDataPermissionService(),
		deptService:       c.GetDepartmentService(),
	}
}

func (c *Container) GetDataPermissionService() *service.DataPermissionService {
	return service.NewDataPermissionService(
		c.GetPermissionRepository(),
		c.GetRoleRepository(),
		c.GetDepartmentRepository(),
		c.GetPermissionCache(),
	)
}

func (c *Container) GetUserDeptService() *service.UserDeptService {
	return service.NewUserDeptService(
		c.GetUserRepository(),
		c.GetDepartmentRepository(),
		c.GetPermissionCache(),
	)
}

func (c *Container) GetDepartmentCommandHandler() *handlers.DepartmentCommandHandler {
	return &handlers.DepartmentCommandHandler{
		tx:          c.GetTransaction(),
		deptRepo:    c.GetDepartmentRepository(),
		deptService: c.GetDepartmentService(),
	}
}

func (c *Container) GetDepartmentQueryHandler() *handlers.DepartmentQueryHandler {
	return &handlers.DepartmentQueryHandler{
		deptService: c.GetDepartmentService(),
	}
}

func (c *Container) GetDepartmentService() *service.DepartmentService {
	return service.NewDepartmentService(
		c.GetDepartmentRepository(),
		c.GetPermissionCache(),
	)
}
