// errors 包提供统一的错误码和错误信息定义
// 所有业务错误码和错误信息都在此集中管理
package errors

import "fmt"

// 业务错误码常量
const (
	CodeSuccess      = 0   // 成功
	CodeBadRequest   = 400 // 请求参数错误
	CodeUnauthorized = 401 // 未授权（需要登录或token无效）
	CodeForbidden    = 403 // 禁止访问
	CodeConflict     = 409 // 资源冲突（如用户已存在）
	CodeInternalErr  = 500 // 服务器内部错误
)

// 错误信息常量
const (
	// 通用错误信息
	MsgSuccess       = "Success"
	MsgBadRequest    = "请求参数无效"
	MsgInternalError = "服务器内部错误"
	MsgUnauthorized  = "未授权"

	// 用户相关错误信息
	MsgUsernamePasswordEmpty = "用户名或密码不能为空"
	MsgUsernameLengthInvalid = "用户名长度必须在4-15个字符之间"
	MsgUsernameFormatInvalid = "用户名格式验证失败"
	MsgUsernameInvalidChars  = "用户名只能包含字母和数字"
	MsgUsernameDuplicate     = "用户名重复"
	MsgUsernamePasswordError = "用户名或密码错误"
	MsgUserNotFound          = "用户不存在"
	MsgUserQueryFailed       = "查询用户失败"
	MsgUserCreateFailed      = "创建用户失败"
	MsgUserInfoFormatError   = "用户信息格式错误"
	MsgNotLoggedIn           = "未登录"

	// 密码相关错误信息
	MsgPasswordEncryptFailed = "密码加密失败"

	// Token相关错误信息
	MsgTokenGenerateFailed   = "生成token失败"
	MsgTokenQueryFailed      = "查询token数量失败"
	MsgTokenDeactivateFailed = "停用旧token失败"
	MsgTokenSaveFailed       = "保存token失败"
	MsgTokenDuplicate        = "token已存在，请稍后重试"
	MsgTokenMissing          = "missing authorization header"
	MsgTokenInvalidFormat    = "invalid authorization header format"
	MsgTokenInvalid          = "invalid token"

	// 服务器错误信息
	MsgServerPanic = "server panic"
)

// BusinessError 业务错误类型
// 包含错误码和错误信息
type BusinessError struct {
	Code    int    // 业务错误码
	Message string // 错误信息
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	return e.Message
}

// New 创建业务错误
func New(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// Newf 创建带格式化的业务错误
func Newf(code int, format string, args ...interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// GetCode 获取错误码
func (e *BusinessError) GetCode() int {
	return e.Code
}

// GetMessage 获取错误信息
func (e *BusinessError) GetMessage() string {
	return e.Message
}

// 预定义的业务错误
var (
	// 通用错误
	ErrBadRequest    = New(CodeBadRequest, MsgBadRequest)
	ErrInternalError = New(CodeInternalErr, MsgInternalError)
	ErrUnauthorized  = New(CodeUnauthorized, MsgUnauthorized)

	// 用户相关错误
	ErrUsernamePasswordEmpty = New(CodeBadRequest, MsgUsernamePasswordEmpty)
	ErrUsernameLengthInvalid = New(CodeBadRequest, MsgUsernameLengthInvalid)
	ErrUsernameFormatInvalid = New(CodeBadRequest, MsgUsernameFormatInvalid)
	ErrUsernameInvalidChars  = New(CodeBadRequest, MsgUsernameInvalidChars)
	ErrUsernameDuplicate     = New(CodeConflict, MsgUsernameDuplicate)
	ErrUsernamePasswordError = New(CodeUnauthorized, MsgUsernamePasswordError)
	ErrUserQueryFailed       = New(CodeInternalErr, MsgUserQueryFailed)
	ErrUserCreateFailed      = New(CodeInternalErr, MsgUserCreateFailed)
	ErrUserInfoFormatError   = New(CodeInternalErr, MsgUserInfoFormatError)
	ErrNotLoggedIn           = New(CodeUnauthorized, MsgNotLoggedIn)

	// 密码相关错误
	ErrPasswordEncryptFailed = New(CodeInternalErr, MsgPasswordEncryptFailed)

	// Token相关错误
	ErrTokenGenerateFailed   = New(CodeInternalErr, MsgTokenGenerateFailed)
	ErrTokenQueryFailed      = New(CodeInternalErr, MsgTokenQueryFailed)
	ErrTokenDeactivateFailed = New(CodeInternalErr, MsgTokenDeactivateFailed)
	ErrTokenSaveFailed       = New(CodeInternalErr, MsgTokenSaveFailed)
	ErrTokenDuplicate        = New(CodeConflict, MsgTokenDuplicate)
	ErrTokenMissing          = New(CodeUnauthorized, MsgTokenMissing)
	ErrTokenInvalidFormat    = New(CodeUnauthorized, MsgTokenInvalidFormat)
)

// NewTokenInvalid 创建token无效错误（需要传入具体错误信息）
func NewTokenInvalid(err error) *BusinessError {
	return Newf(CodeUnauthorized, "%s: %s", MsgTokenInvalid, err.Error())
}

// NewServerPanic 创建服务器panic错误
func NewServerPanic(panicValue interface{}) *BusinessError {
	return Newf(CodeInternalErr, "%s: %v", MsgServerPanic, panicValue)
}
