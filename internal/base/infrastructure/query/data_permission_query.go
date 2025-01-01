package query

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type IDataPermissionQuery interface {
	// GetByRoleID 获取角色的数据权限
	GetByRoleID(ctx context.Context, roleID int64) (*dto.DataPermissionDto, herrors.Herr)
}
