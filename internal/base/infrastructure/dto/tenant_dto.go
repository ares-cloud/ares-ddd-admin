package dto

// TenantDto 租户数据传输对象
type TenantDto struct {
	ID          string   `json:"id"`          // ID
	Code        string   `json:"code"`        // 租户编码
	Name        string   `json:"name"`        // 租户名称
	Domain      string   `json:"domain"`      // 域名
	Description string   `json:"description"` // 描述
	IsDefault   int8     `json:"isDefault"`   // 是否默认租户
	Status      int8     `json:"status"`      // 状态
	AdminUser   *UserDto `json:"adminUser"`   // 管理员用户
	ExpireTime  int64    `json:"expireTime"`  // 过期时间
	CreatedAt   int64    `json:"createdAt"`   // 创建时间
	UpdatedAt   int64    `json:"updatedAt"`   // 更新时间
}
