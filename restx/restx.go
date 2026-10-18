// Package restx REST 层响应工具：把 errorx 错误写成前端约定的 {code, msg} 包裹。
// 前端 axios 拦截器契约：HTTP 200 且 res.data.code === 0 视为成功，否则弹 res.data.msg。
package restx

import (
	"net/http"

	"github.com/wangminghui0622-bit/admin-common/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 统一响应体。
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ErrorHandler 业务错误统一写 HTTP 200 + {code, msg}（前端按 code !== 0 弹错）。
// 在 main 里用 httpx.SetErrorHandler(restx.ErrorHandler) 挂载。
func ErrorHandler(err error) (int, any) {
	code, msg := errorx.FromErr(err)
	return http.StatusOK, &Body{Code: code, Msg: msg}
}

// WriteErr 供中间件直接写错误响应（可指定 HTTP 状态，如 401 触发前端跳登录）。
func WriteErr(w http.ResponseWriter, status, code int, msg string) {
	httpx.WriteJson(w, status, &Body{Code: code, Msg: msg})
}
