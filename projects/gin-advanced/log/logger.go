package log

import (
	"fmt"
	"gin-advanced/core"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ============ 字段构建器 ============

// Field 构建 zap.Field
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Fields 批量构建 zap.Fields
func Fields(fields []interface{}) []zap.Field {
	var zapFields []zap.Field
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			zapFields = append(zapFields, zap.Any(key, fields[i+1]))
		}
	}
	return zapFields
}

// String 构建字符串字段
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Uint 构建无符号整数字段
func Uint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

// Int 构建整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 构建64位整数字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Bool 构建布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Duration 构建时间间隔字段
func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// Sprintf 格式化字符串
func Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// ============ 日志方法 ============

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	core.ZapLog.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	core.ZapLog.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	core.ZapLog.Warn(msg, fields...)
}

// Fatal 致命日志
func Fatal(msg string, fields ...zap.Field) {
	core.ZapLog.Fatal(msg, fields...)
}

// ============ 快捷方法 ============

// D 调试（带默认字段）
func D(msg string) {
	core.ZapLog.Debug(msg)
}

// I 信息
func I(msg string) {
	core.ZapLog.Info(msg)
}

// W 警告
func W(msg string) {
	core.ZapLog.Warn(msg)
}

// E 错误
func E(msg string) {
	core.ZapLog.Error(msg)
}

// ============ 带结构化字段的快捷方法 ============

// WithField 带单个字段
func WithField(key string, value interface{}) *zap.Logger {
	return core.ZapLog.With(zap.Any(key, value))
}

// WithFields 带多个字段
func WithFields(fields map[string]interface{}) *zap.Logger {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return core.ZapLog.With(zapFields...)
}

// ============ 模块化日志 ============

// MySQL MySQL模块日志
func MySQL(msg string, fields ...zap.Field) {
	core.ZapMySQL.Info(msg, fields...)
}

// Redis Redis模块日志
func Redis(msg string, fields ...zap.Field) {
	core.ZapRedis.Info(msg, fields...)
}

// Gin Gin框架日志
func Gin(msg string, fields ...zap.Field) {
	core.ZapGin.Info(msg, fields...)
}

// System 系统日志
func System(msg string, fields ...zap.Field) {
	core.ZapSystem.Info(msg, fields...)
}

// ============ 业务日志快捷方法 ============

// Operation 操作日志（记录用户操作）
func Operation(userID uint, action string, details map[string]interface{}) {
	fields := []zap.Field{
		zap.Uint("user_id", userID),
		zap.String("action", action),
		zap.Time("timestamp", time.Now()),
	}
	for k, v := range details {
		fields = append(fields, zap.Any(k, v))
	}
	core.ZapLog.Info("用户操作", fields...)
}

// Request 请求日志
func Request(method, path string, statusCode int, duration time.Duration, userID uint) {
	core.ZapGin.Info("HTTP请求",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
		zap.Uint("user_id", userID),
	)
}

// ErrorWithStack 带栈信息的错误日志
func ErrorWithStack(err error, msg string) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	core.ZapLog.Error(msg,
		zap.Error(err),
		zap.String("caller", fmt.Sprintf("%s:%d", file, line)),
		zap.String("function", fn.Name()),
	)
}

// ============ 分级日志文件输出 ============

// LogByLevel 根据日志级别记录到不同文件
func LogByLevel(level string, msg string, fields ...zap.Field) {
	switch strings.ToLower(level) {
	case "debug":
		core.ZapLog.Debug(msg, fields...)
	case "info":
		core.ZapLog.Info(msg, fields...)
	case "warn":
		core.ZapLog.Warn(msg, fields...)
	case "error":
		core.ZapLog.Error(msg, fields...)
	default:
		core.ZapLog.Info(msg, fields...)
	}
}
