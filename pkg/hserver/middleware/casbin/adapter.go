package casbin

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2/persist"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/casbin/casbin/v2/model"
)

type CasbinAdapter struct {
	permRepo IPermissionsRepository
}

func NewCasbinAdapter(permRepo IPermissionsRepository) *CasbinAdapter {
	return &CasbinAdapter{
		permRepo: permRepo,
	}
}

// LoadPolicy 从数据库加载策略
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	roles, err := a.permRepo.FindAllEnabled(ctx)
	if err != nil {
		return err
	}

	for _, r := range roles {
		if len(r.Permissions) > 0 {
			for _, perm := range r.Permissions {
				// 添加策略: p, roleCode, tenantID, method, path
				line := fmt.Sprintf("p, %s, %s, %s, %s", r.Code, r.TenantID, perm.Method, perm.Path)
				hlog.Debug("Loading policy:", line)
				err := persist.LoadPolicyArray([]string{"p", r.Code, r.TenantID, perm.Method, perm.Path}, model)
				if err != nil {
					hlog.Errorf("load policy error: %v", err)
					return err
				}
			}
		}
	}

	return nil
}

// SavePolicy 保存策略到数据库
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	return nil // 只读模式
}

// AddPolicy 添加策略
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy 移除策略
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy 移除过滤后的策略
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
