package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// sysUserRepo ， 用户数据层
type sysUserRepo struct {
	*baserepo.BaseRepo[entity.SysUser, string]
}

// NewSysUserRepo ， 用户数据层工厂方法
// 参数：
//
//	data ： desc
//
// 返回值：
//
//	biz.ISysUserRepo ：desc
func NewSysUserRepo(data database.IDataBase) repository.ISysUserRepo {
	model := new(entity.SysUser)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.SysUserRole{}); err != nil {
		hlog.Fatalf("sync sys user tables to db error: %v", err)
	}
	return &sysUserRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.SysUser, string](data, entity.SysUser{}),
	}
}
func (s sysUserRepo) GetByUsername(ctx context.Context, username string) (*entity.SysUser, error) {
	var result *entity.SysUser
	err := s.Db(ctx).Where("username = ?", username).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s sysUserRepo) DeleteRoleByUserId(ctx context.Context, userId string) error {
	return s.Db(ctx).Where("user_id = ?", userId).Delete(&entity.SysUserRole{}).Error
}
