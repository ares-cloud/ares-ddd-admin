package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/internal/base/shared/dto"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

type OperationLogQueryHandler struct {
	repo repository.IOperationLogRepository
}

func NewOperationLogQueryHandler(repo repository.IOperationLogRepository) *OperationLogQueryHandler {
	return &OperationLogQueryHandler{
		repo: repo,
	}
}

// HandleList 处理查询操作日志列表
func (h *OperationLogQueryHandler) HandleList(ctx context.Context, q *queries.ListOperationLogQuery) (*models.PageRes[dto.OperationLogDto], herrors.Herr) {
	tm := time.Now()
	if q.Month != "" {
		// 解析查询月份
		month, err := time.Parse("200601", q.Month)
		if err != nil {
			return nil, herrors.NewBadReqError("invalid month format")
		}
		tm = month
	}

	// 构建查询条件
	qb := query.NewQueryBuilder()
	if q.Username != "" {
		qb.Where("username", query.Like, "%"+q.Username+"%")
	}
	if q.Module != "" {
		qb.Where("ip", query.Like, "%"+q.Module+"%")
	}
	if q.Action != "" {
		qb.Where("status", query.Eq, q.Action)
	}
	if q.StartTime > 0 {
		qb.Where("login_time", query.Gte, time.Unix(q.StartTime, 0))
	}
	if q.EndTime > 0 {
		qb.Where("login_time", query.Lte, time.Unix(q.EndTime, 0))
	}

	// 设置排序
	qb.OrderBy("created_at", false)

	// 设置分页
	qb.WithPage(&q.Page)

	tenant := actx.GetTenantId(ctx)
	// 获取总数
	total, err := h.repo.Count(ctx, tenant, tm, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	logs, err := h.repo.Find(ctx, tenant, tm, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToOperationLogDtoList(logs)

	return &models.PageRes[dto.OperationLogDto]{
		List:  dtoList,
		Total: total,
	}, nil
}
