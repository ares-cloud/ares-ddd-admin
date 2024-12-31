package dto

// UserDto 用户数据传输对象
type UserDto struct {
	ID             string  `json:"id"`             // 用户ID
	TenantID       string  `json:"tenantId"`       // 租户ID
	Username       string  `json:"username"`       // 用户名
	Avatar         string  `json:"avatar"`         // 头像
	Name           string  `json:"name"`           // 姓名
	Nickname       string  `json:"nickname"`       // 昵称
	Phone          string  `json:"phone"`          // 手机号
	Email          string  `json:"email"`          // 邮箱
	Remark         string  `json:"remark"`         // 备注
	InvitationCode string  `json:"invitationCode"` // 邀请码
	Status         int8    `json:"status"`         // 状态,1启用,2禁用
	RoleIds        []int64 `json:"roleIds"`        // 角色ID列表
	CreatedAt      int64   `json:"createdAt"`      // 创建时间
	UpdatedAt      int64   `json:"updatedAt"`      // 更新时间
}
