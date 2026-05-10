package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============ 统一响应结构 ============

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Page    int         `json:"page"`
	PageSize int       `json:"page_size"`
	Total   int64       `json:"total"`
	List    interface{} `json:"list"`
}

// ============ 成功响应 ============

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Created 创建成功
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code: 201,
		Msg:  "创建成功",
		Data: data,
	})
}

// ============ 分页响应 ============

// Page 分页成功响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:     200,
		Msg:      "success",
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	})
}

// ============ 错误响应 ============

// Error 错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 500,
		Msg:  msg,
	})
}

// ErrorWithCode 自定义错误码
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: 400,
		Msg:  msg,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Msg:  msg,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{
		Code: 403,
		Msg:  msg,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{
		Code: 404,
		Msg:  msg,
	})
}

// ============ 错误码定义 ============

const (
	CodeSuccess       = 200
	CodeCreated       = 201
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeServerError   = 500
	CodeParamError    = 4001
	CodeTokenExpired  = 4002
	CodeTokenInvalid  = 4003
	CodeNoPermission  = 4004
	CodeResourceExist = 4005
	CodeResourceNotFound = 4006
)

// ErrorCode 错误码响应
func ErrorCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context) {
	ErrorCode(c, CodeParamError, "参数错误")
}

// TokenExpired Token过期
func TokenExpired(c *gin.Context) {
	ErrorCode(c, CodeTokenExpired, "登录已过期")
}

// TokenInvalid Token无效
func TokenInvalid(c *gin.Context) {
	ErrorCode(c, CodeTokenInvalid, "Token无效")
}

// NoPermission 无权限
func NoPermission(c *gin.Context) {
	ErrorCode(c, CodeNoPermission, "没有权限")
}

// ResourceExist 资源已存在
func ResourceExist(c *gin.Context) {
	ErrorCode(c, CodeResourceExist, "资源已存在")
}

// ResourceNotFound 资源不存在
func ResourceNotFound(c *gin.Context) {
	ErrorCode(c, CodeResourceNotFound, "资源不存在")
}
