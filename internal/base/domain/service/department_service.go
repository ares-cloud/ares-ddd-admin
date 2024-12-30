package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/events"

	domanevent "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
)

type DepartmentService struct {
	repo     repository.IDepartmentRepository
	eventBus *events.EventBus
}

func NewDepartmentService(repo repository.IDepartmentRepository, eventBus *events.EventBus) *DepartmentService {
	return &DepartmentService{repo: repo, eventBus: eventBus}
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

	if err := s.repo.Update(ctx, dept); err != nil {
		return err
	}

	// 获取部门下的用户
	//users, err := s.repo.GetDepartmentUsers(ctx, dept.ID)
	//if err != nil {
	//	return err
	//}

	// 发布部门更新事件
	//userIDs := make([]string, len(users))
	//for i, user := range users {
	//	userIDs[i] = user.ID
	//}

	event := domanevent.NewDepartmentEvent(dept.TenantID, dept.ID, domanevent.DepartmentUpdated)
	return s.eventBus.Publish(ctx, event)
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
	return nil, nil //s.repo.GetDepartmentTree(ctx)
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

// TransferUser 调动用户部门
func (s *DepartmentService) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) error {
	// 1. 检查用户是否存在
	//user, err := s.repo.FindByID(ctx, userID)
	//if err != nil {
	//	return err
	//}
	//if user == nil {
	//	return fmt.Errorf("用户不存在: %s", userID)
	//}

	// 2. 检查原部门是否存在
	fromDept, err := s.repo.GetByID(ctx, fromDeptID)
	if err != nil {
		return err
	}
	if fromDept == nil {
		return fmt.Errorf("原部门不存在: %s", fromDeptID)
	}

	// 3. 检查目标部门是否存在
	toDept, err := s.repo.GetByID(ctx, toDeptID)
	if err != nil {
		return err
	}
	if toDept == nil {
		return fmt.Errorf("目标部门不存在: %s", toDeptID)
	}

	//// 4. 执行部门调动
	//if err := s.repo.TransferUser(ctx, userID, fromDeptID, toDeptID); err != nil {
	//	return err
	//}
	//
	//// 5. 发布用户部门变更事件
	//event := domanevent.NewUserDeptEvent(user.TenantID, userID, fromDeptID, toDeptID)
	//return s.eventBus.Publish(ctx, event)
	return nil
}

func (s *DepartmentService) GetDepartmentUsers(ctx context.Context, deptID string) ([]*model.User, error) {
	//return s.repo.GetDepartmentUsers(ctx, deptID)
	return nil, nil
}
