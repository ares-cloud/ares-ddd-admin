package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate   = validator.New()
	translator *ut.UniversalTranslator
	trans      ut.Translator
	customTags = make(map[string]string)
)

func init() {
	// 1. 注册标签处理器
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// 2. 初始化翻译器
	zhT := zh.New()
	enT := en.New()
	translator = ut.New(enT, zhT, enT)

	// 默认使用中文
	trans, _ = translator.GetTranslator("zh")

	// 3. 注册翻译器
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)

	// 4. 注册自定义验证规则
	registerCustomValidations()

	// 5. 注册自定义翻译
	registerCustomTranslations()
}

// 注册自定义验证规则
func registerCustomValidations() {
	// 手机号验证
	_ = validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) == 11
	})

	// 密码强度验证
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		pwd := fl.Field().String()
		// 至少8位,包含大小写字母和数字
		if len(pwd) < 8 {
			return false
		}
		var hasUpper, hasLower, hasNumber bool
		for _, c := range pwd {
			switch {
			case c >= 'A' && c <= 'Z':
				hasUpper = true
			case c >= 'a' && c <= 'z':
				hasLower = true
			case c >= '0' && c <= '9':
				hasNumber = true
			}
		}
		return hasUpper && hasLower && hasNumber
	})

}

// 注册自定义翻译
func registerCustomTranslations() {
	// 注册自定义tag翻译
	customTags = map[string]string{
		"mobile":   "{0}必须是有效的手机号",
		"password": "{0}必须至少8位,包含大小写字母和数字",
		"chinese":  "{0}只能包含中文和中文标点",
	}

	for tag, msg := range customTags {
		registerTranslation(tag, msg)
	}
}

// 注册翻译
func registerTranslation(tag string, msg string) {
	_ = validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

// Validate 结构体验证
func Validate(v interface{}) herrors.Herr {
	err := validate.Struct(v)
	if err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			return herrors.NewBadReqHError(err)
		}

		var errMsgs []string
		for _, e := range errs {
			// 优先使用翻译后的错误信息
			if msg, ok := customTags[e.Tag()]; ok {
				errMsgs = append(errMsgs, fmt.Sprintf(msg, e.Field()))
			} else {
				errMsgs = append(errMsgs, e.Translate(trans))
			}
		}
		return herrors.NewBadReqHError(fmt.Errorf(strings.Join(errMsgs, "; ")))
	}
	return nil
}

// SetLanguage 设置语言
func SetLanguage(lang string) error {
	if t, ok := translator.GetTranslator(lang); ok {
		trans = t
		switch lang {
		case "en":
			return enTranslations.RegisterDefaultTranslations(validate, trans)
		case "zh":
			return zhTranslations.RegisterDefaultTranslations(validate, trans)
		}
	}
	return fmt.Errorf("unsupported language: %s", lang)
}
