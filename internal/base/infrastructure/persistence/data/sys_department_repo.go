package data

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type sysDepartmentRepo struct {
	*baserepo.BaseRepo[entity.Department, string]
}

func NewSysDepartmentRepo(data database.IDataBase) repository.ISysDepartmentRepo {
	model := new(entity.Department)
	// 同步表
	if err := data.DB(context.Background()).AutoMigrate(model, &entity.UserDepartment{}); err != nil {
		hlog.Fatalf("sync sys department tables to db error: %v", err)
	}
	return &sysDepartmentRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.Department, string](data, entity.Department{}),
	}
}

// GetByCode 根据编码获取部门
func (r *sysDepartmentRepo) GetByCode(ctx context.Context, code string) (*entity.Department, error) {
	var dept entity.Department
	err := r.Db(ctx).Where("code = ?", code).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// GetByParentID 获取子部门
func (r *sysDepartmentRepo) GetByParentID(ctx context.Context, parentID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.Db(ctx).Where("parent_id = ?", parentID).Order("sort").Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

// GetByUserID 获取用户部门关联
func (r *sysDepartmentRepo) GetByUserID(ctx context.Context, userID string) ([]*entity.UserDepartment, error) {
	var list []*entity.UserDepartment
	err := r.Db(ctx).Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// FindByIds 根据ID列表查询部门
func (r *sysDepartmentRepo) FindByIds(ctx context.Context, ids []string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.Db(ctx).Where("id IN ?", ids).Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

// DelById 删除部门（包括关联关系）
func (r *sysDepartmentRepo) DelById(ctx context.Context, id string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 删除用户部门关联
		if err := r.Db(ctx).Where("dept_id = ?", id).Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		// 删除部门
		return r.Db(ctx).Delete(&entity.Department{}, "id = ?", id).Error
	})
}

// DelByIdUnScoped 硬删除部门（包括关联关系）
func (r *sysDepartmentRepo) DelByIdUnScoped(ctx context.Context, id string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 删除用户部门关联
		if err := r.Db(ctx).Unscoped().Where("dept_id = ?", id).Delete(&entity.UserDepartment{}).Error; err != nil {
			return err
		}

		// 删除部门
		return r.Db(ctx).Unscoped().Delete(&entity.Department{}, "id = ?", id).Error
	})
}

// FindAllEnabled 查询所有启用的部门
func (r *sysDepartmentRepo) FindAllEnabled(ctx context.Context) ([]*entity.Department, error) {
	var list []*entity.Department
	err := r.Db(ctx).Where("status = 1").Order("sort").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
