package model

import "time"

// LoginType 登录类型
type LoginType int8

const (
	LoginTypeAdmin  LoginType = 1 // 管理端登录
	LoginTypeMember LoginType = 2 // 前台用户登录
)

// LoginLog 登录日志领域模型
type LoginLog struct {
	ID        int64     // ID
	UserID    string    // 用户ID
	Username  string    // 用户名
	TenantID  string    // 租户ID
	LoginType LoginType // 登录类型
	IP        string    // 登录IP
	Location  string    // 登录地点
	Device    string    // 登录设备
	OS        string    // 操作系统
	Browser   string    // 浏览器
	Status    int8      // 登录状态(1:成功 2:失败)
	Message   string    // 登录消息
	LoginTime int64     // 登录时间
	CreatedAt int64     // 创建时间
	UpdatedAt int64     // 更新时间
}

// NewLoginLog 创建登录日志
func NewLoginLog(userID, username, tenantID string, loginType LoginType) *LoginLog {
	now := time.Now().Unix()
	return &LoginLog{
		UserID:    userID,
		Username:  username,
		TenantID:  tenantID,
		LoginType: loginType,
		Status:    1,
		LoginTime: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetLoginInfo 设置登录信息
func (l *LoginLog) SetLoginInfo(ip, location, device, os, browser string) {
	l.IP = ip
	l.Location = location
	l.Device = device
	l.OS = os
	l.Browser = browser
	l.UpdatedAt = time.Now().Unix()
}

// SetLoginStatus 设置登录状态
func (l *LoginLog) SetLoginStatus(status int8, message string) {
	l.Status = status
	l.Message = message
	l.UpdatedAt = time.Now().Unix()
}
