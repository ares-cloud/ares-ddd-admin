package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, string, error) {
	// Configure captcha parameters
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	// Generate the captcha
	id, b64s, a, err := c.Generate()
	if err != nil {
		return "", "", "", err
	}
	return id, b64s, a, nil
}

// GetMathCaptcha create return id, b64s, err
func GetMathCaptcha(width, height int64) (string, string, string, error) {
	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 60
	}
	// 配置算术验证码
	driver := base64Captcha.NewDriverMath(
		60,                                 // 高度
		200,                                // 宽度
		2,                                  // 噪声数量
		base64Captcha.OptionShowHollowLine, // 干扰线选项
		&color.RGBA{R: 99, G: 253, B: 124, A: 100}, // 背景颜色
		nil, // 使用默认字体存储
		nil, // 使用默认字体
	)
	// 生成验证码实例
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验码
	return captcha.Generate()
}

// GetDigitCaptcha create return id, b64s, err
func GetDigitCaptcha() (string, string, string, error) {
	// 配置算术验证码
	driver := base64Captcha.NewDriverDigit(80, 240, 3, 0.5, 80)
	// 生成验证码实例
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验码
	return captcha.Generate()
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaId, value string) bool {
	// Verify the captcha
	if base64Captcha.DefaultMemStore.Verify(captchaId, value, true) {
		return true
	}
	return false
}