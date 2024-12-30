package errors

import "errors"

// 创建领域错误
func New(text string) error {
	return errors.New(text)
}

// 包装领域错误
func Wrap(err error, message string) error {
	return errors.New(message + ": " + err.Error())
}
