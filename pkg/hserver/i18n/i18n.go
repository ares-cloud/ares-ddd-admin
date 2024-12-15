package i18n

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	hertzI18n "github.com/hertz-contrib/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func Handler(rootPath string, def language.Tag, tags ...language.Tag) app.HandlerFunc {
	return hertzI18n.Localize(
		hertzI18n.WithBundle(&hertzI18n.BundleCfg{
			RootPath:         rootPath,
			AcceptLanguage:   tags,
			DefaultLanguage:  def,
			FormatBundleFile: "yaml",
			UnmarshalFunc:    yaml.Unmarshal,
		}),
		hertzI18n.WithGetLangHandle(func(c context.Context, ctx *app.RequestContext, defaultLang string) string {
			lang := ctx.GetHeader("lang")
			if len(lang) == 0 {
				return defaultLang
			}
			return string(lang)
		}),
	)
}
