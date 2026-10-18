// Package errorx 定义 admin-zero 统一错误码，与 simple-admin 前端约定一致。
package errorx

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"google.golang.org/grpc/status"
)

const (
	OK              = 0
	NotExist        = 4
	Unauthorized    = 7 // 无权限访问（casbin 拒绝）
	AuthFailed      = 8 // 登录失败（验证码/密码错误）
	ServerError     = 9
	InvalidArgument = 10
)

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

func New(code int, msg string) error {
	return &Error{Code: code, Msg: msg}
}

func Newf(code int, format string, args ...any) error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, args...)}
}

// codeMsgRe 匹配 Error.Error() 的输出格式。业务错误经 RPC 传递时会被 grpc 包成
// status.Error(codes.Unknown, "code=N msg=...")，从 message 里还原业务码。
var codeMsgRe = regexp.MustCompile(`^code=(\d+) msg=(.*)$`)

// FromErr 把任意 error 归一为 (code, msg)，供 REST 层写响应。
func FromErr(err error) (int, string) {
	if err == nil {
		return OK, "success"
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code, e.Msg
	}
	// RPC 返回的业务错误被 grpc 包成 status（code=Unknown, message="code=N msg=..."），
	// 类型信息已丢失，只能从 desc 解析还原业务码
	if st, ok := status.FromError(err); ok {
		if m := codeMsgRe.FindStringSubmatch(st.Message()); m != nil {
			if code, convErr := strconv.Atoi(m[1]); convErr == nil {
				return code, m[2]
			}
		}
		// 其它 grpc 状态错误（如 DeadlineExceeded）无业务码可映射，按服务错误返回
		return ServerError, st.Message()
	}
	return ServerError, err.Error()
}
