package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
)

type DepartmentService struct {
	repo repository.IDepartmentRepository
}

func NewDepartmentService(repo repository.IDepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(ctx context.Context, dept *model.Department) error {
	// 1. 检查编码是否重复
	exists, err := s.repo.GetByCode(ctx, dept.Code)
	if err != nil {
		return err
	}
	if exists != nil {
		return errors.New("部门编码已存在")
	}

	// 2. 检查父部门是否存在
	if dept.ParentID != "" {
		parent, err := s.repo.GetByID(ctx, dept.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("父部门不存在")
		}
	}

	// 4. 创建部门
	return s.repo.Create(ctx, dept)
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(ctx context.Context, dept *model.Department) error {
	// 1. 检查部门是否存在
	exists, err := s.repo.GetByID(ctx, dept.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return errors.New("部门不存在")
	}

	// 2. 检查编码是否重复
	if exists.Code != dept.Code {
		codeExists, err := s.repo.GetByCode(ctx, dept.Code)
		if err != nil {
			return err
		}
		if codeExists != nil && codeExists.ID != dept.ID {
			return errors.New("部门编码已存在")
		}
	}

	// 3. 检查父部门是否存在
	if dept.ParentID != "" && dept.ParentID != exists.ParentID {
		parent, err := s.repo.GetByID(ctx, dept.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("父部门不存在")
		}
	}

	return s.repo.Update(ctx, dept)
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(ctx context.Context, id string) error {
	// 1. 检查是否有子部门
	children, err := s.repo.GetByParentID(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("存在子部门,不能删除")
	}

	return s.repo.Delete(ctx, id)
}

// GetDepartmentList 获取部门列表
func (s *DepartmentService) GetDepartmentList(ctx context.Context, query *repository.ListDepartmentQuery) ([]*model.Department, error) {
	return s.repo.List(ctx, query)
}

// GetDepartmentTree 获取部门树
func (s *DepartmentService) GetDepartmentTree(ctx context.Context, parentID string) ([]*model.Department, error) {
	// 1. 获取所有部门
	depts, err := s.repo.List(ctx, &repository.ListDepartmentQuery{})
	if err != nil {
		return nil, err
	}

	// 2. 构建部门映射
	deptMap := make(map[string]*model.Department)
	for _, dept := range depts {
		deptMap[dept.ID] = dept
		dept.Children = make([]*model.Department, 0) // 初始化子部门切片
	}

	// 3. 构建树形结构
	var root []*model.Department
	for _, dept := range depts {
		if dept.ParentID == parentID {
			// 根节点
			root = append(root, dept)
		} else {
			// 将当前部门添加到父部门的子部门列表中
			if parent, ok := deptMap[dept.ParentID]; ok {
				parent.Children = append(parent.Children, dept)
			}
		}
	}

	return root, nil
}

// MoveDepartment 移动部门
func (s *DepartmentService) MoveDepartment(ctx context.Context, deptID, targetParentID string) error {
	// 1. 检查部门是否存在
	dept, err := s.repo.GetByID(ctx, deptID)
	if err != nil {
		return err
	}
	if dept == nil {
		return fmt.Errorf("department not found: %s", deptID)
	}

	// 2. 检查目标父部门是否存在
	if targetParentID != "" {
		parent, err := s.repo.GetByID(ctx, targetParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return fmt.Errorf("target parent department not found: %s", targetParentID)
		}
	}

	// 3. 检查是否形成循环引用
	if err := s.checkCircularReference(ctx, deptID, targetParentID); err != nil {
		return err
	}

	// 4. 更新父部门
	dept.UpdateParent(targetParentID)

	// 5. 保存更新
	return s.repo.Update(ctx, dept)
}

// checkCircularReference 检查是否形成循环引用
func (s *DepartmentService) checkCircularReference(ctx context.Context, deptID, targetParentID string) error {
	if targetParentID == "" {
		return nil
	}

	// 获取目标父部门的所有上级部门
	var parentIDs []string
	currentID := targetParentID
	for currentID != "" {
		parent, err := s.repo.GetByID(ctx, currentID)
		if err != nil {
			return err
		}
		if parent == nil {
			break
		}
		// 检查是否形成循环
		if parent.ID == deptID {
			return fmt.Errorf("circular reference detected")
		}
		parentIDs = append(parentIDs, parent.ID)
		currentID = parent.ParentID
	}

	return nil
}
