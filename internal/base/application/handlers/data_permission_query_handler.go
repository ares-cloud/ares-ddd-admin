package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type DataPermissionQueryHandler struct {
	permQuery query.IDataPermissionQuery
}

func NewDataPermissionQueryHandler(permQuery query.IDataPermissionQuery) *DataPermissionQueryHandler {
	return &DataPermissionQueryHandler{
		permQuery: permQuery,
	}
}

// HandleGetByRoleID 处理获取角色数据权限
func (h *DataPermissionQueryHandler) HandleGetByRoleID(ctx context.Context, query queries.GetDataPermissionQuery) (*dto.DataPermissionDto, herrors.Herr) {
	return h.permQuery.GetByRoleID(ctx, query.RoleID)
}
