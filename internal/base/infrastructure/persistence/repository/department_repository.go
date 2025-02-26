package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	drepository "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/entity"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/mapper"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/baserepo"
)

type ISysDepartmentRepo interface {
	baserepo.IBaseRepo[entity.Department, string]
	GetByCode(ctx context.Context, code string) (*entity.Department, error)
	GetByParentID(ctx context.Context, parentID string) ([]*entity.Department, error)
	GetByUserID(ctx context.Context, userID string) ([]*entity.UserDepartment, error)
	GetDeptByUserID(ctx context.Context, userID string) ([]*entity.Department, error)
	FindByIds(ctx context.Context, ids []string) ([]*entity.Department, error)
	AssignUsers(ctx context.Context, deptID string, userIDs []string) error
	RemoveUsers(ctx context.Context, deptID string, userIDs []string) error
	TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error
}

type departmentRepository struct {
	repo   ISysDepartmentRepo
	mapper *mapper.DepartmentMapper
}

func NewDepartmentRepository(repo ISysDepartmentRepo) drepository.IDepartmentRepository {
	return &departmentRepository{
		repo:   repo,
		mapper: &mapper.DepartmentMapper{},
	}
}

func (r *departmentRepository) Create(ctx context.Context, dept *model.Department) error {
	deptEntity := r.mapper.ToEntity(dept)
	deptEntity.ID = r.repo.GenStringId()
	_, err := r.repo.Add(ctx, deptEntity)
	return err
}

func (r *departmentRepository) Update(ctx context.Context, dept *model.Department) error {
	deptEntity := r.mapper.ToEntity(dept)
	return r.repo.EditById(ctx, deptEntity.ID, deptEntity)
}

func (r *departmentRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

// AssignUsers 分配用户到部门
func (r *departmentRepository) AssignUsers(ctx context.Context, deptID string, userIDs []string) error {
	return r.repo.AssignUsers(ctx, deptID, userIDs)
}

// RemoveUsers 从部门移除用户
func (r *departmentRepository) RemoveUsers(ctx context.Context, deptID string, userIDs []string) error {
	return r.repo.RemoveUsers(ctx, deptID, userIDs)
}

// ExistsByCode 检查部门编码是否存在
func (r *departmentRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	dept, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return dept != nil, nil
}

// FindByID 根据ID查询部门
func (r *departmentRepository) FindByID(ctx context.Context, id string) (*model.Department, error) {
	dept, err := r.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(dept), nil
}

// FindByCode 根据编码查询部门
func (r *departmentRepository) FindByCode(ctx context.Context, code string) (*model.Department, error) {
	dept, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(dept), nil
}

// GetTreeByParentID 获取指定父部门下的部门树
func (r *departmentRepository) GetTreeByParentID(ctx context.Context, parentID string) ([]*model.Department, error) {
	// 1. 获取直接子部门
	children, err := r.repo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	domainChildren := r.mapper.ToDomainList(children)

	// 2. 递归获取每个子部门的子部门
	for _, child := range domainChildren {
		subChildren, err := r.GetTreeByParentID(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		for _, subChild := range subChildren {
			child.AddChild(subChild)
		}
	}

	return domainChildren, nil
}

// TransferUser 调动用户部门
func (r *departmentRepository) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	return r.repo.TransferUser(ctx, userID, fromDeptID, toDeptID)
}

// GetByParentID 获取指定父部门下的直接子部门列表
func (r *departmentRepository) GetByParentID(ctx context.Context, parentID string) ([]*model.Department, error) {
	// 获取子部门列表
	children, err := r.repo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	return r.mapper.ToDomainList(children), nil
}
