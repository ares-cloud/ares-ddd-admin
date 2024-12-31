package dto

import "github.com/ares-cloud/ares-ddd-admin/internal/base/domain/model"

// DepartmentDto 部门DTO
type DepartmentDto struct {
	ID          string           `json:"id"`          // 部门ID
	ParentID    string           `json:"parentId"`    // 父部门ID
	Name        string           `json:"name"`        // 部门名称
	Code        string           `json:"code"`        // 部门编码
	Sort        int32            `json:"sort"`        // 排序
	AdminID     string           `json:"adminId"`     //管理员ID
	Leader      string           `json:"leader"`      // 负责人
	Phone       string           `json:"phone"`       // 联系电话
	Email       string           `json:"email"`       // 邮箱
	Status      int8             `json:"status"`      // 部门状态
	Description string           `json:"description"` // 描述
	Children    []*DepartmentDto `json:"children"`    // 子部门
}

// DepartmentTreeDto 部门树DTO
type DepartmentTreeDto struct {
	ID       string               `json:"id"`       // 部门ID
	ParentID string               `json:"parentId"` // 父部门ID
	Name     string               `json:"name"`     // 部门名称
	Children []*DepartmentTreeDto `json:"children"` // 子部门
}

// ToDepartmentDto 将领域模型转换为DTO
func ToDepartmentDto(dept *model.Department) *DepartmentDto {
	if dept == nil {
		return nil
	}
	return &DepartmentDto{
		ID:          dept.ID,
		ParentID:    dept.ParentID,
		Name:        dept.Name,
		Code:        dept.Code,
		Sort:        dept.Sequence,
		AdminID:     dept.AdminID,
		Leader:      dept.Leader,
		Phone:       dept.Phone,
		Email:       dept.Email,
		Status:      dept.Status,
		Description: dept.Description,
		Children:    ToDepartmentDtoList(dept.Children),
	}
}

// ToDepartmentDtoList 将领域模型列表转换为DTO列表
func ToDepartmentDtoList(depts []*model.Department) []*DepartmentDto {
	if depts == nil {
		return nil
	}
	dtos := make([]*DepartmentDto, 0, len(depts))
	for _, dept := range depts {
		if dto := ToDepartmentDto(dept); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ToDepartmentTreeDto 将领域模型转换为树形DTO
func ToDepartmentTreeDto(dept *model.Department) *DepartmentTreeDto {
	if dept == nil {
		return nil
	}
	return &DepartmentTreeDto{
		ID:       dept.ID,
		ParentID: dept.ParentID,
		Name:     dept.Name,
		Children: ToDepartmentTreeDtoList(dept.Children),
	}
}

// ToDepartmentTreeDtoList 将领域模型列表转换为树形DTO列表
func ToDepartmentTreeDtoList(depts []*model.Department) []*DepartmentTreeDto {
	if depts == nil {
		return nil
	}
	dtos := make([]*DepartmentTreeDto, 0, len(depts))
	for _, dept := range depts {
		if dto := ToDepartmentTreeDto(dept); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
