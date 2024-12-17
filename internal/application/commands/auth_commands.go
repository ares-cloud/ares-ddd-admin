package commands

type LoginCommand struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaKey  string `json:"captchaKey" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
}

type RefreshTokenCommand struct {
	Token string `json:"token" binding:"required"`
}
