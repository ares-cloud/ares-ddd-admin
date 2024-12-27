package handlers

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type DataPermissionQueryHandler struct {
	permRepo repository.IDataPermissionRepository
}

func NewDataPermissionQueryHandler(permRepo repository.IDataPermissionRepository) *DataPermissionQueryHandler {
	return &DataPermissionQueryHandler{
		permRepo: permRepo,
	}
}

// HandleGetByRoleID 处理获取角色数据权限
func (h *DataPermissionQueryHandler) HandleGetByRoleID(ctx context.Context, query queries.GetDataPermissionQuery) (*dto.DataPermissionDto, herrors.Herr) {
	// 查询数据权限
	perm, err := h.permRepo.GetByRoleID(ctx, query.RoleID)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil // 没有配置数据权限时返回nil
		}
		return nil, herrors.NewServerHError(err)
	}

	return dto.ToDataPermissionDto(perm), nil
}
