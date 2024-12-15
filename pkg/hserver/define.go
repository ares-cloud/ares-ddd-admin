package hserver

import (
	"context"
	"fmt"

	"github.com/ares-cloud/ares-ddd-admin/pkg/constant"
	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
	"github.com/ares-cloud/ares-ddd-admin/pkg/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"

	"net/http"
	"reflect"
	"strconv"
	"time"
)

type Router interface {
	ConfigRoutes(r *server.Hertz, t token.IToken) // 配置路由
}

// ResponseResult 响应结果
type ResponseResult struct {
	Data  interface{}
	Error *herrors.HError
}

func DefaultResponseResult() *ResponseResult {
	return &ResponseResult{}
}
func (r *ResponseResult) WithData(data interface{}) *ResponseResult {
	r.Data = data
	return r
}

func (r *ResponseResult) WithError(err *herrors.HError) *ResponseResult {
	r.Error = err
	return r
}

// ServiceFunc 实际提供服务的函数
type ServiceFunc[T any] func(ctx context.Context, par *T) *ResponseResult

// ServiceNotParFunc 实际提供服务的函数(无参数)
type ServiceNotParFunc func(ctx context.Context) *ResponseResult

// Handler 接口处理器
type Handler[T any] struct {
	Context        context.Context
	RequestContext *app.RequestContext
	Param          *T
	Error          error
}

// NewHandler [T any] ， handler 工厂函数
// 参数：
//
//	ctx ： desc
//	c ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func NewHandler[T any](ctx context.Context, c *app.RequestContext) *Handler[T] {
	return &Handler[T]{Context: ctx, RequestContext: c}
}

// NewHandlerFu [T any] ， handlerFun 工厂函数
// 参数：
//
//	fun ： desc
//
// 返回值：
//
//	app.HandlerFunc ：desc
func NewHandlerFu[T any](fun ServiceFunc[T]) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		NewHandler[T](c, ctx).WithBinder().WithBinderPage().Do(fun)
	}
}

// NewNotParHandlerFu [T any] ， 无参数的处理器
// 参数：
//
//	fun ： desc
//
// 返回值：
//
//	app.HandlerFunc ：desc
func NewNotParHandlerFu(fun ServiceNotParFunc) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		NewHandler[int](c, ctx).DoNotPar(fun)
	}
}

// WithBinder ， 绑定并验证参数
// 参数：
//
//	param ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func (h *Handler[T]) WithBinder() *Handler[T] {
	h.Param = new(T)

	if err := h.RequestContext.BindAndValidate(h.Param); err != nil {
		h.ParamErr(err)
	}
	return h
}

// WithBinderPage ， 绑定并验证分页
// 参数：
//
//	param ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func (h *Handler[T]) WithBinderPage() *Handler[T] {
	// 使用反射检查是否有 Page 字段
	paramValue := reflect.ValueOf(h.Param).Elem()
	paramType := paramValue.Type()

	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)
		if field.Name == "Page" && field.Type.Kind() == reflect.Ptr {
			pageValue := paramValue.Field(i)
			if pageValue.IsNil() {
				// 如果 Page 字段是 nil，则初始化它
				pageValue.Set(reflect.New(field.Type.Elem()))
			}
			// 获取并转换 current
			currentStr := h.RequestContext.Query("current")
			if currentStr != "" {
				current, err := strconv.ParseInt(currentStr, 10, 64)
				if err != nil {
					h.ParamErr(err)
					return h
				}
				pageValue.Elem().FieldByName("Current").SetInt(current)
			}
			// 获取并转换 size
			sizeStr := h.RequestContext.Query("size")
			if sizeStr != "" {
				size, err := strconv.ParseInt(sizeStr, 10, 64)
				if err != nil {
					h.ParamErr(err)
					return h
				}
				pageValue.Elem().FieldByName("Size").SetInt(size)
			}
			break
		}
	}
	return h
}

// Do ， 执行 server 函数
// 参数：
//
//	serviceFunc ： desc
//
// 返回值：
func (h *Handler[T]) Do(serviceFunc ServiceFunc[T]) {
	// 错误处理
	if h.Error != nil {
		return
	}
	if h.RequestContext.IsAborted() {
		return
	}
	// 调用服务函数
	res := serviceFunc(h.Context, h.Param)
	if res.Error != nil {
		ResponseFailureErr(h.Context, h.RequestContext, res.Error)
	} else {
		ResponseSuccess(h.Context, h.RequestContext, res.Data)
	}
}

// DoNotPar ， 无参数服务函数
// 参数：
//
//	serviceFunc ： desc
//
// 返回值：
func (h *Handler[T]) DoNotPar(serviceFunc ServiceNotParFunc) {
	// 错误处理
	if h.Error != nil {
		return
	}
	if h.RequestContext.IsAborted() {
		return
	}
	// 调用服务函数
	res := serviceFunc(h.Context)
	if res.Error != nil {
		ResponseFailureErr(h.Context, h.RequestContext, res.Error)
	} else {
		ResponseSuccess(h.Context, h.RequestContext, res.Data)
	}
}

// ParamErr ， 参数错误
// 参数：
//
//	param ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func (h *Handler[T]) ParamErr(err error) {
	h.Error = err
	msg := hertzI18n.MustGetMessage(h.Context, herrors.ReasonParameterError)
	if msg == "" {
		msg = "Param err"
	}
	h.RequestContext.JSON(http.StatusOK, utils.H{
		constant.RespCode:      constant.StatusInvalidParam,
		constant.RespMsg:       msg,
		constant.ErrMsg:        fmt.Sprintf("%s，fail err ：%+v", msg, h.Error),
		constant.RespReason:    herrors.ReasonParameterError,
		constant.RespTimestamp: time.Now().Unix(),
	})
	h.RequestContext.Abort()
}
