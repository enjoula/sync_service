// response 包提供统一的HTTP响应格式封装
// 所有API响应都使用统一的JSON格式，包含code、message、data和trace_id字段
package response

import (
	"net/http"
	"video-service/internal/errors"

	"github.com/gin-gonic/gin"
)

// TraceIDKey 是context中存储追踪ID的键名
// 与middleware包中的TraceIDKey保持一致
const TraceIDKey = "trace_id"

// Response 定义统一的API响应结构
type Response struct {
	Code    int         `json:"code"`               // 业务状态码：0表示成功，其他值表示错误
	Message string      `json:"message"`            // 响应消息
	Data    interface{} `json:"data"`               // 响应数据（成功时包含业务数据，失败时为nil）
	TraceID string      `json:"trace_id,omitempty"` // 请求追踪ID（可选，用于分布式追踪）
}

// 错误码常量从 errors 包导入
const (
	CodeSuccess      = errors.CodeSuccess      // 成功
	CodeBadRequest   = errors.CodeBadRequest   // 请求参数错误
	CodeUnauthorized = errors.CodeUnauthorized // 未授权（需要登录或token无效）
	CodeConflict     = errors.CodeConflict     // 资源冲突（如用户已存在）
	CodeInternalErr  = errors.CodeInternalErr  // 服务器内部错误
)

// Success 返回成功响应
// 参数：
//
//	c: Gin上下文
//	data: 要返回的业务数据
func Success(c *gin.Context, data interface{}) {
	// 从context中获取追踪ID（由Trace中间件设置）
	traceID, _ := c.Get(TraceIDKey)

	// 返回JSON响应，HTTP状态码为200
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: errors.MsgSuccess,
		Data:    data,
		TraceID: toStr(traceID),
	})
}

// SuccessMsg 返回成功响应（带自定义消息）
// 参数：
//
//	c: Gin上下文
//	msg: 自定义成功消息
//	data: 要返回的业务数据
func SuccessMsg(c *gin.Context, msg string, data interface{}) {
	traceID, _ := c.Get(TraceIDKey)
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: msg,
		Data:    data,
		TraceID: toStr(traceID),
	})
}

// Error 返回错误响应
// 参数：
//
//	c: Gin上下文
//	code: 业务错误码
//	msg: 错误消息
func Error(c *gin.Context, code int, msg string) {
	traceID, _ := c.Get(TraceIDKey)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
		TraceID: toStr(traceID),
	})
}

// InternalError 返回服务器内部错误响应
// 参数：
//
//	c: Gin上下文
//	err: 错误对象
func InternalError(c *gin.Context, err error) {
	traceID, _ := c.Get(TraceIDKey)
	c.JSON(http.StatusOK, Response{
		Code:    CodeInternalErr,
		Message: err.Error(),
		Data:    nil,
		TraceID: toStr(traceID),
	})
}

// toStr 将interface{}类型转换为string类型
// 用于安全地获取context中的追踪ID
func toStr(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
