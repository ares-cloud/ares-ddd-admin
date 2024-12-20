package handlers

import (
	"context"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/internal/application/dto"
	"github.com/ares-cloud/ares-ddd-admin/internal/application/queries"
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/repository"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/models"
)

type LoginLogQueryHandler struct {
	loginLogRepo repository.ILoginLogRepository
}

func NewLoginLogQueryHandler(loginLogRepo repository.ILoginLogRepository) *LoginLogQueryHandler {
	return &LoginLogQueryHandler{
		loginLogRepo: loginLogRepo,
	}
}
func (h *LoginLogQueryHandler) HandleAppList(ctx context.Context, q *queries.ListLoginLogsQuery) (*models.PageRes[dto.LoginLogDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	qb.Where("login_type", query.Eq, model.LoginTypeMember)
	return h.HandleList(ctx, qb, q)
}

func (h *LoginLogQueryHandler) HandleAdminList(ctx context.Context, q *queries.ListLoginLogsQuery) (*models.PageRes[dto.LoginLogDto], herrors.Herr) {
	// 构建查询条件
	qb := query.NewQueryBuilder()
	qb.Where("login_type", query.Eq, model.LoginTypeAdmin)
	return h.HandleList(ctx, qb, q)
}

func (h *LoginLogQueryHandler) HandleList(ctx context.Context, qb *query.QueryBuilder, q *queries.ListLoginLogsQuery) (*models.PageRes[dto.LoginLogDto], herrors.Herr) {
	// 解析查询月份
	month, err := time.Parse("200601", q.Month)
	if err != nil {
		return nil, herrors.NewBadReqError("invalid month format")
	}

	if q.Username != "" {
		qb.Where("username", query.Like, "%"+q.Username+"%")
	}
	if q.IP != "" {
		qb.Where("ip", query.Like, "%"+q.IP+"%")
	}
	if q.Status != 0 {
		qb.Where("status", query.Eq, q.Status)
	}
	if q.StartTime > 0 {
		qb.Where("login_time", query.Gte, time.Unix(q.StartTime, 0))
	}
	if q.EndTime > 0 {
		qb.Where("login_time", query.Lte, time.Unix(q.EndTime, 0))
	}

	// 设置排序
	qb.OrderBy("login_time", false)

	// 设置分页
	qb.WithPage(&q.Page)

	tenant := actx.GetTenantId(ctx)
	// 获取总数
	total, err := h.loginLogRepo.Count(ctx, tenant, month, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	logs, err := h.loginLogRepo.Find(ctx, tenant, month, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 转换为DTO
	dtoList := dto.ToLoginLogDtoList(logs)

	return &models.PageRes[dto.LoginLogDto]{
		List:  dtoList,
		Total: total,
	}, nil
}
