package hserver

import (
	"context"
	"net/http"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
)

/*
ResponseSuccess 返回成功响应数据
调用示例:
1.成功时，不需要返回数据: server.ResponseSuccess(c, nil)
2.成功时，需要返回数据: server.ResponseSuccess(c, gin.H{"name": "xim","age": 18})
*/
func ResponseSuccess(ctx context.Context, c *app.RequestContext, data interface{}) {
	if data == nil {
		data = utils.H{}
	}
	message := hertzI18n.MustGetMessage(ctx, constant.ReasonSuccess)
	if message == "" {
		message = "ok"
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: http.StatusOK, constant.RespMsg: message, constant.RespData: data, constant.RespReason: constant.ReasonSuccess})
}

/*
ResponseFailure 返回失败响应数据
调用示例:
*/
func ResponseFailure(ctx context.Context, c *app.RequestContext, code int, reason, msg string, data interface{}) {
	if data == nil {
		data = utils.H{}
	}
	if reason != "" {
		i18Mag := hertzI18n.MustGetMessage(ctx, reason)
		if i18Mag != "" {
			msg = i18Mag
		}
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: code, constant.RespMsg: msg, constant.RespData: data, constant.RespReason: reason, constant.RespTimestamp: time.Now().Format("2006-01-02 15:04:05")})
}

/*
ResponseFailureErr 返回失败响应数据
调用示例:
*/
func ResponseFailureErr(ctx context.Context, c *app.RequestContext, err *herrors.HError) {
	code := err.Code
	if code == 0 {
		code = http.StatusInternalServerError
	}
	msg := err.DefMessage
	if err.Reason != "" {
		i18Mag := hertzI18n.MustGetMessage(ctx, err.Reason)
		if i18Mag != "" {
			msg = i18Mag
		}
	}
	errMsg := err.DefMessage
	if err.BusinessError != nil {
		errMsg = err.BusinessError.Error()
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: code, constant.RespMsg: msg, constant.ErrMsg: errMsg, constant.RespReason: err.Reason, constant.RespTimestamp: time.Now().Format("2006-01-02 15:04:05")})
}
