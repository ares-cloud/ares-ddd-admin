package impl

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"

	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/converter"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/infrastructure/persistence/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/db_query"
)

type PermissionsQueryService struct {
	permRepo             repository.IPermissionsRepo
	permissionsConverter *converter.PermissionsConverter
}

func NewPermissionsQueryService(
	permRepo repository.IPermissionsRepo,
	permissionsConverter *converter.PermissionsConverter,
) *PermissionsQueryService {
	return &PermissionsQueryService{
		permRepo:             permRepo,
		permissionsConverter: permissionsConverter,
	}
}

// FindByID 根据ID查询权限
func (s *PermissionsQueryService) FindByID(ctx context.Context, id int64) (*dto.PermissionsDto, error) {
	perm, err := s.permRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, nil
	}
	return s.permissionsConverter.ToDTO(perm, nil), nil
}

// Find 查询权限列表
func (s *PermissionsQueryService) Find(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.PermissionsDto, int64, herrors.Herr) {
	perms, err := s.permRepo.Find(ctx, qb)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}

	total, err := s.permRepo.Count(ctx, qb)
	if err != nil {
		return nil, 0, herrors.QueryFail(err)
	}

	return s.permissionsConverter.ToDTOList(perms), total, nil
}

// FindTreeByType 查询权限树
func (s *PermissionsQueryService) FindTreeByType(ctx context.Context, permType int8) ([]*dto.PermissionsDto, error) {
	perms, _, err := s.permRepo.GetTreeByType(ctx, permType)
	if err != nil {
		return nil, err
	}
	return s.permissionsConverter.ToTreeDTOList(perms), nil
}

// FindAllEnabled 查询所有启用的权限
func (s *PermissionsQueryService) FindAllEnabled(ctx context.Context) ([]*dto.PermissionsDto, error) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	qb.Where("status", db_query.Eq, 1)
	qb.OrderBy("sequence", true)

	// 查询数据
	perms, err := s.permRepo.Find(ctx, qb)
	if err != nil {
		return nil, err
	}

	return s.permissionsConverter.ToDTOList(perms), nil
}

// GetSimplePermissionsTree 获取简化的权限树
func (s *PermissionsQueryService) GetSimplePermissionsTree(ctx context.Context) ([]*dto.PermissionsTreeDto, error) {
	// 1. 查询所有权限
	perms, _, err := s.permRepo.GetTreeByType(ctx, 1) // 1表示菜单类型
	if err != nil {
		return nil, err
	}

	// 2. 转换为树形结构
	return s.permissionsConverter.ToSimpleTreeDTOList(perms), nil
}
