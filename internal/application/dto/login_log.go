package dto

import (
	"github.com/ares-cloud/ares-ddd-admin/internal/domain/model"
)

type LoginLogDto struct {
	ID        int64  `json:"id"`         // ID
	UserID    string `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	TenantID  string `json:"tenant_id"`  // 租户ID
	IP        string `json:"ip"`         // 登录IP
	Location  string `json:"location"`   // 登录地点
	Device    string `json:"device"`     // 登录设备
	OS        string `json:"os"`         // 操作系统
	Browser   string `json:"browser"`    // 浏览器
	Status    int8   `json:"status"`     // 登录状态
	Message   string `json:"message"`    // 登录消息
	LoginTime int64  `json:"login_time"` // 登录时间
}

func ToLoginLogDto(log *model.LoginLog) *LoginLogDto {
	if log == nil {
		return nil
	}
	return &LoginLogDto{
		ID:        log.ID,
		UserID:    log.UserID,
		Username:  log.Username,
		TenantID:  log.TenantID,
		IP:        log.IP,
		Location:  log.Location,
		Device:    log.Device,
		OS:        log.OS,
		Browser:   log.Browser,
		Status:    log.Status,
		Message:   log.Message,
		LoginTime: log.LoginTime,
	}
}

func ToLoginLogDtoList(logs []*model.LoginLog) []*LoginLogDto {
	if logs == nil {
		return nil
	}
	dtos := make([]*LoginLogDto, 0, len(logs))
	for _, log := range logs {
		if dto := ToLoginLogDto(log); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
