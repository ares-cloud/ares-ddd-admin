package casbin

import (
	"context"
	"fmt"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	psb "github.com/ares-cloud/ares-ddd-admin/pkg/hserver/middleware/casbin"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type RepositoryImpl struct {
	rr repository.ISysRoleRepo
	pr repository.IPermissionsRepo
}

func NewRepositoryImpl(rr repository.ISysRoleRepo, pr repository.IPermissionsRepo) psb.IPermissionsRepository {
	return &RepositoryImpl{
		rr: rr,
		pr: pr,
	}
}

// FindAllEnabled 获取所有启用的角色及其权限
func (r *RepositoryImpl) FindAllEnabled(ctx context.Context) ([]*psb.Role, error) {
	ctx = context.Background()
	// 获取所有启用的角色
	roles, err := r.rr.FindAllEnabled(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "casbin [FindAllEnabled] error: %v", err)
		return nil, err
	}
	if len(roles) == 0 {
		return []*psb.Role{}, nil
	}

	// 获取角色ID列表
	roleIds := make([]int64, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
	}
	roleResourcesMap, err := r.pr.GetResourcesByRolesGrouped(ctx, roleIds)
	// 转换为 casbin 角色格式
	var casbinRoles []*psb.Role
	for _, role := range roles {
		casbinRole := &psb.Role{
			Id:       fmt.Sprintf("%d", role.ID),
			Code:     role.Code,
			TenantID: role.TenantID,
		}
		// 添加权限
		if resources, ok := roleResourcesMap[role.ID]; ok {
			for _, resource := range resources {
				casbinRole.Permissions = append(casbinRole.Permissions, psb.ApiPermissions{
					Id:     fmt.Sprintf("%d", resource.PermissionsID),
					Method: resource.Method,
					Path:   resource.Path,
				})
			}
		}

		casbinRoles = append(casbinRoles, casbinRole)
	}

	return casbinRoles, nil
}
