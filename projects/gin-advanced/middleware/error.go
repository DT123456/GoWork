package middleware

import (
	"fmt"
	"gin-advanced/log"
	"gin-advanced/models"
	"gin-advanced/service"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorCode 错误码定义
type ErrorCode struct {
	Code    int
	Message string
}

var (
	// 通用错误码
	ErrBadRequest      = ErrorCode{400, "请求参数错误"}
	ErrUnauthorized    = ErrorCode{401, "未授权，请登录"}
	ErrForbidden       = ErrorCode{403, "禁止访问"}
	ErrNotFound        = ErrorCode{404, "资源不存在"}
	ErrServerError     = ErrorCode{500, "服务器内部错误"}
	ErrTimeout         = ErrorCode{504, "请求超时"}
	ErrTooManyRequests = ErrorCode{429, "请求过于频繁"}

	// 业务错误码
	ErrUserNotFound  = ErrorCode{1001, "用户不存在"}
	ErrUserExist     = ErrorCode{1002, "用户已存在"}
	ErrPasswordWrong = ErrorCode{1003, "密码错误"}
	ErrTokenExpired  = ErrorCode{1004, "Token已过期"}
	ErrTokenInvalid  = ErrorCode{1005, "Token无效"}
	ErrNoPermission  = ErrorCode{1006, "没有权限"}
	ErrRoleExist     = ErrorCode{2001, "角色已存在"}
	ErrRoleNotFound  = ErrorCode{2002, "角色不存在"}
	ErrResourceExist = ErrorCode{3001, "资源已存在"}
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录错误日志
				stack := string(debug.Stack())
				log.E(fmt.Sprintf("[PANIC] %v", r))
				log.E(fmt.Sprintf("[PANIC STACK] %s", stack))

				// 异步记录系统日志
				go recordErrorLog(c, r)

				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "服务器内部错误，请稍后重试",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// guessModule 猜测模块名称
func guessModule(path string) string {
	if strings.HasPrefix(path, "/admin") {
		return "管理后台"
	}
	if strings.HasPrefix(path, "/user") {
		return "用户模块"
	}
	if strings.HasPrefix(path, "/api") {
		return "API接口"
	}
	return "其他"
}

// handleError 处理错误
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 记录错误日志
	log.E(fmt.Sprintf("[ERROR] %s %s - %v", c.Request.Method, c.Request.URL.Path, err))

	switch {
	case c.Request.URL.Path == "/login" && c.Request.Method == "POST":
		c.JSON(http.StatusOK, gin.H{
			"code": ErrUnauthorized.Code,
			"msg":  "账号或密码错误",
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": ErrServerError.Code,
			"msg":  ErrServerError.Message,
		})
	}
}

// recordErrorLog 记录错误日志
func recordErrorLog(c *gin.Context, err interface{}) {
	userID, _ := c.Get("userID")
	var uid uint
	if u, ok := userID.(uint); ok {
		uid = u
	}
	module := guessModule(c.Request.URL.Path)

	systemLog := &models.SystemLog{
		UserID:   uid,
		Action:   fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		Module:   module,
		Status:   0,
		Code:     "ERROR",
		ErrorMsg: fmt.Sprintf("%v", err),
		Message:  "系统异常",
		IP:       c.ClientIP(),
		Duration: 0,
	}

	service.CreateSystemLog(systemLog)
}

// CustomError 自定义错误响应
func CustomError(c *gin.Context, errCode ErrorCode, msg ...string) {
	message := errCode.Message
	if len(msg) > 0 {
		message = msg[0]
	}
	c.JSON(http.StatusOK, gin.H{
		"code": errCode.Code,
		"msg":  message,
	})
}

// BadRequestError 请求参数错误
func BadRequestError(c *gin.Context, msg ...string) {
	CustomError(c, ErrBadRequest, msg...)
}

// UnauthorizedError 未授权错误
func UnauthorizedError(c *gin.Context, msg ...string) {
	CustomError(c, ErrUnauthorized, msg...)
}

// ForbiddenError 禁止访问错误
func ForbiddenError(c *gin.Context, msg ...string) {
	CustomError(c, ErrForbidden, msg...)
}

// NotFoundError 资源不存在错误
func NotFoundError(c *gin.Context, msg ...string) {
	CustomError(c, ErrNotFound, msg...)
}

// ServerError 服务器内部错误
func ServerError(c *gin.Context, msg ...string) {
	CustomError(c, ErrServerError, msg...)
}

// ValidationError 参数校验错误
func ValidationError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": ErrBadRequest.Code,
		"msg":  msg,
	})
}

// BusinessError 业务错误
func BusinessError(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, errCode ErrorCode, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": errCode.Code,
		"msg":  errCode.Message,
		"data": data,
	})
}

// RequestTimeout 请求超时处理
func RequestTimeout(c *gin.Context) {
	log.W(fmt.Sprintf("[TIMEOUT] %s %s", c.Request.Method, c.Request.URL.Path))

	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code": ErrTimeout.Code,
		"msg":  ErrTimeout.Message,
	})
}

// LogErrors 记录请求错误日志
func LogErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 记录错误
		if len(c.Errors) > 0 {
			duration := time.Since(start).Milliseconds()
			log.E(fmt.Sprintf("[FAILED] %s %s - status=%d, duration=%dms, errors=%s",
				c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, c.Errors.String()))
		}
	}
}
