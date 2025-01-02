package service

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/errors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/events"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	pkgEvent "github.com/ares-cloud/ares-ddd-admin/pkg/events"
)

type DepartmentService struct {
	deptRepo repository.IDepartmentRepository
	userRepo repository.IUserRepository
	eventBus pkgEvent.IEventBus
}

func NewDepartmentService(
	deptRepo repository.IDepartmentRepository,
	userRepo repository.IUserRepository,
	eventBus pkgEvent.IEventBus,
) *DepartmentService {
	return &DepartmentService{
		deptRepo: deptRepo,
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(ctx context.Context, dept *model.Department) herrors.Herr {
	// 1. 验证部门信息
	if err := dept.Validate(); herrors.HaveError(err) {
		return err
	}

	// 2. 检查部门编码是否存在
	exists, err := s.deptRepo.ExistsByCode(ctx, dept.Code)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if exists {
		return errors.DepartmentExists(dept.Code)
	}

	// 3. 如果有父部门,检查父部门是否存在
	if dept.HasParent() {
		parent, err := s.deptRepo.FindByID(ctx, dept.ParentID)
		if err != nil {
			return herrors.NewServerHError(err)
		}
		if parent == nil {
			return errors.ParentDepartmentNotFound(dept.ParentID)
		}
		if !parent.IsEnabled() {
			return errors.ParentDepartmentDisabled(dept.ParentID)
		}
	}

	// 4. 创建部门
	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return errors.DepartmentCreateFailed(err)
	}

	// 5. 发布部门创建事件
	if err := s.eventBus.Publish(ctx, events.NewDepartmentEvent(dept.TenantID, dept.ID, events.DepartmentCreated)); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(ctx context.Context, dept *model.Department) herrors.Herr {
	// 1. 验证部门信息
	if err := dept.Validate(); herrors.HaveError(err) {
		return err
	}

	// 2. 检查部门是否存在
	oldDept, err := s.deptRepo.FindByID(ctx, dept.ID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if oldDept == nil {
		return errors.DepartmentNotFound(dept.ID)
	}

	// 3. 如果修改了编码,检查新编码是否存在
	if oldDept.Code != dept.Code {
		exists, err := s.deptRepo.ExistsByCode(ctx, dept.Code)
		if err != nil {
			return herrors.NewServerHError(err)
		}
		if exists {
			return errors.DepartmentExists(dept.Code)
		}
	}

	// 4. 如果修改了父部门,检查父部门是否存在且有效
	if oldDept.ParentID != dept.ParentID && dept.HasParent() {
		parent, err := s.deptRepo.FindByID(ctx, dept.ParentID)
		if err != nil {
			return herrors.NewServerHError(err)
		}
		if parent == nil {
			return errors.ParentDepartmentNotFound(dept.ParentID)
		}
		if !parent.IsEnabled() {
			return errors.ParentDepartmentDisabled(dept.ParentID)
		}
	}

	// 5. 更新部门
	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return errors.DepartmentUpdateFailed(err)
	}

	// 6. 发布部门更新事件
	if err := s.eventBus.Publish(ctx, events.NewDepartmentEvent(dept.TenantID, dept.ID, events.DepartmentUpdated)); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(ctx context.Context, id string) herrors.Herr {
	// 1. 检查部门是否存在
	dept, err := s.deptRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DepartmentQueryFailed(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(id)
	}

	// 2. 检查是否有子部门
	children, err := s.deptRepo.GetTreeByParentID(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if len(children) > 0 {
		return errors.HasChildDepartment(id)
	}

	// 3. 删除部门
	if err := s.deptRepo.Delete(ctx, id); err != nil {
		return errors.DepartmentDeleteFailed(err)
	}

	// 4. 发布部门删除事件
	if err := s.eventBus.Publish(ctx, events.NewDepartmentEvent(dept.TenantID, dept.ID, events.DepartmentDeleted)); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// GetDepartmentTree 获取部门树
func (s *DepartmentService) GetDepartmentTree(ctx context.Context, parentID string) ([]*model.Department, herrors.Herr) {
	depts, err := s.deptRepo.GetTreeByParentID(ctx, parentID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return depts, nil
}

// AssignUsers 分配用户到部门
func (s *DepartmentService) AssignUsers(ctx context.Context, deptID string, userIDs []string) herrors.Herr {
	// 1. 检查部门是否存在且有效
	dept, err := s.deptRepo.FindByID(ctx, deptID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(deptID)
	}
	if !dept.IsEnabled() {
		return errors.DepartmentDisabled(deptID)
	}

	// 2. 分配用户
	if err := s.deptRepo.AssignUsers(ctx, deptID, userIDs); err != nil {
		return errors.UserAssignFailed(err)
	}

	// 3. 发布用户分配事件
	event := events.NewUserAssignedEvent(actx.GetTenantId(ctx), deptID, userIDs)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// RemoveUsers 从部门移除用户
func (s *DepartmentService) RemoveUsers(ctx context.Context, deptID string, userIDs []string) herrors.Herr {
	// 1. 检查部门是否存在
	dept, err := s.deptRepo.FindByID(ctx, deptID)
	if err != nil {
		return errors.DepartmentQueryFailed(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(deptID)
	}

	// 2. 移除用户
	if err := s.deptRepo.RemoveUsers(ctx, deptID, userIDs); err != nil {
		return errors.UserRemoveFailed(err)
	}

	// 3. 发布用户移除事件
	event := events.NewUserRemovedEvent(actx.GetTenantId(ctx), deptID, userIDs)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// TransferUser 调动用户部门
func (s *DepartmentService) TransferUser(ctx context.Context, userID string, fromDeptID string, toDeptID string) herrors.Herr {
	// 1. 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if user == nil {
		return errors.UserNotFound(userID)
	}

	// 2. 检查目标部门是否存在且有效
	toDept, err := s.deptRepo.FindByID(ctx, toDeptID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if toDept == nil {
		return errors.DepartmentNotFound(toDeptID)
	}
	if !toDept.IsEnabled() {
		return errors.DepartmentDisabled(toDeptID)
	}

	// 3. 执行调动
	if err = s.deptRepo.TransferUser(ctx, userID, fromDeptID, toDeptID); err != nil {
		return errors.UserTransferFailed(err)
	}

	// 4. 发布用户调动事件
	event := events.NewUserTransferredEvent(actx.GetTenantId(ctx), userID, fromDeptID, toDeptID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// GetByID 获取部门
func (s *DepartmentService) GetByID(ctx context.Context, id string) (*model.Department, herrors.Herr) {
	dept, err := s.deptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DepartmentQueryFailed(err)
	}
	if dept == nil {
		return nil, errors.DepartmentNotFound(id)
	}
	return dept, nil
}

// MoveDepartment 移动部门
func (s *DepartmentService) MoveDepartment(ctx context.Context, id string, targetParentID string) herrors.Herr {
	// 1. 检查部门是否存在
	dept, err := s.deptRepo.FindByID(ctx, id)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(id)
	}

	// 获取原父部门ID
	oldParentID := dept.ParentID

	// 2. 更新父部门
	dept.UpdateParent(targetParentID)
	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return herrors.NewServerHError(err)
	}

	// 3. 发布部门移动事件
	event := events.NewDepartmentMovedEvent(actx.GetTenantId(ctx), id, oldParentID, targetParentID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// SetDepartmentAdmin 设置部门管理员
func (s *DepartmentService) SetDepartmentAdmin(ctx context.Context, deptID string, adminID string) herrors.Herr {
	// 1. 检查部门是否存在
	dept, err := s.deptRepo.FindByID(ctx, deptID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if dept == nil {
		return errors.DepartmentNotFound(deptID)
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, adminID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if user == nil {
		return errors.UserNotFound(adminID)
	}

	// 3. 检查用户是否属于该部门
	belongs, err := s.userRepo.BelongsToDepartment(ctx, adminID, deptID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	if !belongs {
		return errors.UserDepartmentNotFound(adminID, deptID)
	}

	// 4. 设置管理员
	dept.SetAdmin(adminID)
	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return herrors.NewServerHError(err)
	}

	// 5. 发布管理员设置事件
	if err := s.eventBus.Publish(ctx, events.NewDepartmentEvent(dept.TenantID, dept.ID, events.DepartmentUpdated)); err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}
