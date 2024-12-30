package validator

import "regexp"

// 常用正则表达式
var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

// ValidateUsername 验证用户名
// 规则：只允许字母、数字、下划线和连字符
func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePhone 验证手机号
// 规则：1开头的11位数字，第二位是3-9
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// ValidatePassword 验证密码
// 规则：长度至少6位
func ValidatePassword(password string) bool {
	return len(password) >= 6
}

// ValidateLength 验证字符串长度是否在指定范围内
func ValidateLength(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

// ValidateRequired 验证必填字段
func ValidateRequired(str string) bool {
	return str != ""
}
