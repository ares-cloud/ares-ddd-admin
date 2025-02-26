package data

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
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
	err := r.Db(ctx).Where("parent_id = ?", parentID).Order("sequence").Find(&depts).Error
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
func (r *sysDepartmentRepo) GetDeptByUserID(ctx context.Context, userID string) ([]*entity.Department, error) {
	var list []*entity.Department
	err := r.Db(ctx).Model(&entity.UserDepartment{}).
		Joins("LEFT JOIN sys_department ON sys_department.id = user_department.dept_id").
		Where("user_department.user_id = ?", userID).
		Find(&list).Error
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

// AssignUsers 分配用户到部门
func (r *sysDepartmentRepo) AssignUsers(ctx context.Context, deptID string, userIDs []string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 批量创建用户部门关联
		userDepts := make([]*entity.UserDepartment, 0, len(userIDs))
		for _, userID := range userIDs {
			userDepts = append(userDepts, &entity.UserDepartment{
				ID:     r.GenInt64Id(),
				UserID: userID,
				DeptID: deptID,
			})
		}
		return r.Db(ctx).Create(&userDepts).Error
	})
}

// RemoveUsers 从部门移除用户
func (r *sysDepartmentRepo) RemoveUsers(ctx context.Context, deptID string, userIDs []string) error {
	return r.Db(ctx).Where("dept_id = ? AND user_id IN ?", deptID, userIDs).
		Delete(&entity.UserDepartment{}).Error
}

// TransferUser 调动用户部门
func (r *sysDepartmentRepo) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 1. 如果有原部门，先移除
		if fromDeptID != "" {
			if err := r.RemoveUsers(ctx, fromDeptID, []string{userID}); err != nil {
				return err
			}
		}

		// 2. 添加到新部门
		if err := r.AssignUsers(ctx, toDeptID, []string{userID}); err != nil {
			return err
		}

		return nil
	})
}
