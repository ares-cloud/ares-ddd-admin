package commands

// LoginType 登录类型
type LoginType int8

const (
	LoginTypeAdmin  LoginType = 1 // 管理端登录
	LoginTypeMember LoginType = 2 // 前台用户登录
)

type LoginCommand struct {
	Username    string    `json:"username" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	CaptchaKey  string    `json:"captchaKey" binding:"required"`
	CaptchaCode string    `json:"captchaCode" binding:"required"`
	Platform    string    `json:"platform" binding:"required"`
	LoginType   LoginType `json:"login_type" binding:"required"` // 登录类型
	IP          string    `json:"ip"`                            // 登录IP
	Location    string    `json:"location"`                      // 登录地点
	UserAgent   string    `json:"user_agent"`                    // User-Agent
}

type RefreshTokenCommand struct {
	Token string `json:"token" binding:"required"`
}
