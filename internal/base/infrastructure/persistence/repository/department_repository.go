package repository

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"

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
	FindByIds(ctx context.Context, ids []string) ([]*entity.Department, error)
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

func (r *departmentRepository) GetByID(ctx context.Context, id string) (*model.Department, error) {
	deptEntity, err := r.repo.FindById(ctx, id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.mapper.ToDomain(deptEntity), nil
}

func (r *departmentRepository) GetByCode(ctx context.Context, code string) (*model.Department, error) {
	deptEntity, err := r.repo.GetByCode(ctx, code)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.mapper.ToDomain(deptEntity), nil
}

func (r *departmentRepository) GetByParentID(ctx context.Context, parentID string) ([]*model.Department, error) {
	qb := query.NewQueryBuilder()
	qb.Where("parentId", query.Eq, parentID)
	qb.OrderBy("sort", false)

	deptEntities, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(deptEntities), nil
}

func (r *departmentRepository) List(ctx context.Context, req *drepository.ListDepartmentQuery) ([]*model.Department, error) {
	qb := query.NewQueryBuilder()
	if req.Name != "" {
		qb.Where("name", query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		qb.Where("code", query.Like, "%"+req.Code+"%")
	}
	if req.Status != nil {
		qb.Where("status", query.Eq, *req.Status)
	}
	if req.ParentID != "" {
		qb.Where("parentId", query.Eq, req.ParentID)
	}
	qb.OrderBy("sort", false)

	deptEntities, err := r.repo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(deptEntities), nil
}

func (r *departmentRepository) GetUserDepartments(ctx context.Context, userID string) ([]*model.Department, error) {
	// 查询用户部门关联
	userDepts, err := r.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(userDepts) == 0 {
		return []*model.Department{}, nil
	}

	// 获取部门ID列表
	deptIDs := make([]string, 0, len(userDepts))
	for _, ud := range userDepts {
		deptIDs = append(deptIDs, ud.DeptID)
	}

	// 查询部门信息
	depts, err := r.repo.FindByIds(ctx, deptIDs)
	if err != nil {
		return nil, err
	}

	return r.mapper.ToDomainList(depts), nil
}

func (r *departmentRepository) GetAllUserIDs(ctx context.Context) ([]string, error) {
	var userIDs []string
	err := r.repo.Db(ctx).Model(&entity.UserDepartment{}).
		Distinct("userId").
		Pluck("userId", &userIDs).Error
	return userIDs, err
}

// GetParentChain 获取部门的所有上级部门
func (r *departmentRepository) GetParentChain(ctx context.Context, deptID string) ([]*model.Department, error) {
	var chain []*model.Department
	currentID := deptID

	for currentID != "" {
		dept, err := r.GetByID(ctx, currentID)
		if err != nil {
			return nil, err
		}
		if dept == nil {
			break
		}
		chain = append(chain, dept)
		currentID = dept.ParentID
	}

	return chain, nil
}

// GetChildrenRecursively 递归获取所有子部门
func (r *departmentRepository) GetChildrenRecursively(ctx context.Context, parentID string) ([]*model.Department, error) {
	children, err := r.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	var allChildren []*model.Department
	allChildren = append(allChildren, children...)

	for _, child := range children {
		subChildren, err := r.GetChildrenRecursively(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		allChildren = append(allChildren, subChildren...)
	}

	return allChildren, nil
}
