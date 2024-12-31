package rest

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/commands"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/handlers"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver"
	_ "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/base_info"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/jwt"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/oplog"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"
	"github.com/cloudwego/hertz/pkg/route"
)

type DepartmentController struct {
	cmdHandler   *handlers.DepartmentCommandHandler
	queryHandler *handlers.DepartmentQueryHandler
	ef           *casbin.Enforcer
	moduleName   string
}

func NewDepartmentController(cmdHandler *handlers.DepartmentCommandHandler, queryHandler *handlers.DepartmentQueryHandler, ef *casbin.Enforcer) *DepartmentController {
	return &DepartmentController{
		cmdHandler:   cmdHandler,
		queryHandler: queryHandler,
		ef:           ef,
		moduleName:   "部门",
	}
}

func (c *DepartmentController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	dept := v1.Group("/sys/dept", jwt.Handler(t))
	{
		dept.POST("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "新增",
		}), hserver.NewHandlerFu[commands.CreateDepartmentCommand](c.AddDepartment))

		dept.GET("", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.ListDepartmentsQuery](c.DepartmentList))

		dept.PUT("", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "修改",
		}), hserver.NewHandlerFu[commands.UpdateDepartmentCommand](c.UpdateDepartment))

		dept.DELETE("/:id", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](c.DeleteDepartment))

		dept.GET("/:id", casbin.Handler(c.ef), hserver.NewHandlerFu[models.StringIdReq](c.GetDetails))

		dept.GET("/tree", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.GetDepartmentTreeQuery](c.GetDepartmentTree))

		dept.POST("/move", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "移动",
		}), hserver.NewHandlerFu[commands.MoveDepartmentCommand](c.MoveDepartment))
		dept.POST("/admin", hserver.NewHandlerFu[commands.SetDepartmentAdminCommand](c.SetAdmin))
		dept.GET("/:id/users", casbin.Handler(c.ef), hserver.NewHandlerFu[queries.GetDepartmentUsersQuery](c.GetDepartmentUsers))
		dept.POST("/users", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "分配用户",
		}), hserver.NewHandlerFu[commands.AssignUsersToDepartmentCommand](c.AssignUsers))
		dept.PUT("/users", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "移除用户",
		}), hserver.NewHandlerFu[commands.RemoveUsersFromDepartmentCommand](c.RemoveUsers))
		dept.GET("/unassigned-users", casbin.Handler(c.ef),
			hserver.NewHandlerFu[queries.GetUnassignedUsersQuery](c.GetUnassignedUsers))
		dept.POST("/transfer", casbin.Handler(c.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      c.moduleName,
			Action:      "人员调动",
		}), hserver.NewHandlerFu[commands.TransferUserCommand](c.TransferUser))
	}
}

// AddDepartment 添加部门
// @Summary 添加部门
// @Description 添加部门
// @Tags 系统部门
// @ID AddDepartment
// @Accept json
// @Produce json
// @Param req body commands.CreateDepartmentCommand true "部门信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [post]
func (c *DepartmentController) AddDepartment(ctx context.Context, params *commands.CreateDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleCreate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DepartmentList 获取部门列表
// @Summary 获取部门列表
// @Description 获取部门列表
// @Tags 系统部门
// @Accept json
// @Produce json
// @Param req query queries.ListDepartmentsQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=dto.PageResult[dto.DepartmentDto]}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [get]
func (c *DepartmentController) DepartmentList(ctx context.Context, params *queries.ListDepartmentsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleList(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// UpdateDepartment 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 系统部门
// @ID UpdateDepartment
// @Accept json
// @Produce json
// @Param req body commands.UpdateDepartmentCommand true "部门信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept [put]
func (c *DepartmentController) UpdateDepartment(ctx context.Context, params *commands.UpdateDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleUpdate(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Description 删除部门
// @Tags 系统部门
// @ID DeleteDepartment
// @Accept json
// @Produce json
// @Param id path string true "部门ID"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/{id} [delete]
func (c *DepartmentController) DeleteDepartment(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleDelete(ctx, &commands.DeleteDepartmentCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetDetails 获取部门详情
// @Summary 获取部门详情
// @Description 获取部门详情
// @Tags 系统部门
// @ID GetDepartmentDetails
// @Accept json
// @Produce json
// @Param id path string true "部门ID"
// @Success 200 {object} base_info.Success{data=model.Department}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/{id} [get]
func (c *DepartmentController) GetDetails(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGet(ctx, &queries.GetDepartmentQuery{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetDepartmentTree 获取部门树
// @Summary 获取部门树
// @Description 获取部门树形结构
// @Tags 系统部门
// @ID GetDepartmentTree
// @Accept json
// @Produce json
// @Param req query queries.GetDepartmentTreeQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=[]model.Department}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/tree [get]
func (c *DepartmentController) GetDepartmentTree(ctx context.Context, params *queries.GetDepartmentTreeQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGetTree(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// MoveDepartment 移动部门
// @Summary 移动部门
// @Description 移动部门位置
// @Tags 系统部门
// @ID MoveDepartment
// @Accept json
// @Produce json
// @Param req body commands.MoveDepartmentCommand true "移动信息"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/move [post]
func (c *DepartmentController) MoveDepartment(ctx context.Context, params *commands.MoveDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := c.cmdHandler.HandleMove(ctx, params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// SetAdmin 设置部门管理员
// @Summary 设置部门管理员
// @Description 设置部门管理员,同时更新部门负责人信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param req body commands.SetDepartmentAdminCommand true "设置部门管理员参数"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/admin [post]
func (c *DepartmentController) SetAdmin(ctx context.Context, req *commands.SetDepartmentAdminCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleSetAdmin(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}

// GetDepartmentUsers 获取部门用户列表
// @Summary 获取部门用户列表
// @Description 获取部门下的用户列表(分页)
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path string true "部门ID"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页大小"
// @Param username query string false "用户名"
// @Param name query string false "姓名"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.UserDto]}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/{id}/users [get]
func (c *DepartmentController) GetDepartmentUsers(ctx context.Context, req *queries.GetDepartmentUsersQuery) *hserver.ResponseResult {
	return hserver.DefaultResponseResult()
	//data, err := c.queryHandler.HandleGetUsers(ctx, req)
	//if err != nil {
	//	return result.WithError(err)
	//}
	//return result.WithData(data)
}

// AssignUsers 分配用户到部门
// @Summary 分配用户到部门
// @Description 分配用户到部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param req body commands.AssignUsersToDepartmentCommand true "分配用户参数"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/users [post]
func (c *DepartmentController) AssignUsers(ctx context.Context, req *commands.AssignUsersToDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleAssignUsers(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}

// RemoveUsers 从部门移除用户
// @Summary 从部门移除用户
// @Description 从部门移除用户
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param req body commands.RemoveUsersFromDepartmentCommand true "移除用户参数"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/users [delete]
func (c *DepartmentController) RemoveUsers(ctx context.Context, req *commands.RemoveUsersFromDepartmentCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleRemoveUsers(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}

// GetUnassignedUsers 获取未分配部门的用户列表
// @Summary 获取未分配部门的用户列表
// @Description 获取未分配部门的用户列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param req query queries.GetUnassignedUsersQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.UserDto]}
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/unassigned-users [get]
func (c *DepartmentController) GetUnassignedUsers(ctx context.Context, req *queries.GetUnassignedUsersQuery) *hserver.ResponseResult {
	return hserver.DefaultResponseResult()
	//result := hserver.DefaultResponseResult()
	//data, err := c.queryHandler.HandleGetUnassignedUsers(ctx, req)
	//if err != nil {
	//	return result.WithError(err)
	//}
	//return result.WithData(data)
}

// TransferUser 人员部门调动
// @Summary 人员部门调动
// @Description 将用户从一个部门调动到另一个部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param req body commands.TransferUserCommand true "调动参数"
// @Success 200 {object} base_info.Success
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "内部错误"
// @Router /v1/sys/dept/transfer [post]
func (c *DepartmentController) TransferUser(ctx context.Context, req *commands.TransferUserCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	if err := c.cmdHandler.HandleTransferUser(ctx, req); err != nil {
		return result.WithError(err)
	}
	return result
}
