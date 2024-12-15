package sql_injection

import (
	"context"

	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"

	"net/http"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
)

// PreventSQLInjection 中间件函数
func PreventSQLInjection() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 过滤查询参数
		c.QueryArgs().VisitAll(func(key, value []byte) {
			if isSQLInjection(string(value)) {
				hlog.CtxErrorf(ctx, "Potential SQL Injection detected in query parameter %s: %s", key, value)
				i18Mag := hertzI18n.MustGetMessage(ctx, constant.SQLInjectionDetected)
				c.AbortWithStatusJSON(http.StatusOK, utils.H{constant.RespCode: http.StatusForbidden, constant.RespMsg: i18Mag, constant.RespReason: constant.SQLInjectionDetected, constant.RespData: utils.H{}})
				return
			}
			cleanedValue := filterSQLInjection(string(value))
			c.QueryArgs().Set(string(key), cleanedValue)
		})
		// 过滤表单数据
		if c.IsPost() {
			c.PostArgs().VisitAll(func(key, value []byte) {
				if isSQLInjection(string(value)) {
					hlog.CtxErrorf(ctx, "Potential SQL Injection detected in query parameter %s: %s", key, value)
					i18Mag := hertzI18n.MustGetMessage(ctx, constant.SQLInjectionDetected)
					c.AbortWithStatusJSON(http.StatusOK, utils.H{constant.RespCode: http.StatusForbidden, constant.RespMsg: i18Mag, constant.RespReason: constant.SQLInjectionDetected, constant.RespData: utils.H{}})
					return
				}
				cleanedValue := filterSQLInjection(string(value))
				c.PostArgs().Set(string(key), cleanedValue)
			})
		}
		// 如果没有检测到SQL注入，继续处理请求
		c.Next(ctx)
	}
}
func filterSQLInjection(value string) string {
	// 更加全面的过滤列表
	dangerousPatterns := []string{"--", ";", "'", "\"", "/*", "*/", "@@", "@", "char", "nchar", "varchar", "nvarchar",
		"alter", "begin", "cast", "create", "cursor", "declare", "delete", "drop", "end",
		"exec", "execute", "fetch", "insert", "kill", "select", "sys", "sysobjects", "syscolumns",
		"table", "update"}
	for _, pattern := range dangerousPatterns {
		value = strings.ReplaceAll(value, pattern, "")
	}
	return value
}
func isSQLInjection(value string) bool {
	// 简单的正则表达式检测
	injectionPatterns := []string{
		`(?i)union(\s)+select`,
		`(?i)select(\s)+\*`,
		`(?i)insert(\s)+into`,
		`(?i)drop(\s)+table`,
		`(?i)alter(\s)+table`,
		`(?i)delete(\s)+from`,
		`(?i)update(\s)+\w+(\s)+set`,
		`(?i)exec(\s)+\w+`,
	}

	for _, pattern := range injectionPatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return true
		}
	}
	return false
}
