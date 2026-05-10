package middleware

import (
	"bytes"
	"gin-advanced/log"
	"gin-advanced/models"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUsernameByID 根据用户ID获取用户名（支持缓存和数据库fallback）
func GetUsernameByID(userID uint) string {
	if userID == 0 {
		return "未登录"
	}
	// 从Redis缓存获取用户名
	username, err := utils.GetUsernameByCache(userID)
	if err == nil && username != "" {
		return username
	}
	// 缓存未命中，从数据库查询
	user, err := service.GetUserByID(int(userID))
	if err != nil {
		return ""
	}
	return user.Username
}

// ============ 跳过日志记录的路径 ============

var skipLogPaths = map[string]bool{
	"/swagger/":          true, // Swagger 文档
	"/uploads/":          true, // 静态资源
	"/favicon.ico":        true, // 图标
	"/health":             true, // 健康检查
	"/metrics":            true, // 监控指标
}

// shouldSkipLog 判断是否跳过日志记录
func shouldSkipLog(path string) bool {
	// 完全匹配
	if skipLogPaths[path] {
		return true
	}
	// 前缀匹配
	for prefix := range skipLogPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// ============ 请求体捕获 ============

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ============ 请求日志中间件 ============

// RequestLog 统一请求日志中间件（同时记录到 zap 和数据库）
// 替代原来的 AccessLog + SystemLog，减少重复记录
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的请求
		if shouldSkipLog(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()

		// 获取用户信息 - 直接从请求头解析token，避免依赖中间件执行顺序
		uid := getUserIDFromRequest(c)

		// 记录请求基本信息到 zap
		go log.Gin("请求开始",
			log.String("method", c.Request.Method),
			log.String("path", c.Request.URL.Path),
			log.String("query", c.Request.URL.RawQuery),
			log.String("ip", c.ClientIP()),
			log.Uint("user_id", uid),
		)

		// 处理请求
		c.Next()

		// 记录请求完成信息到 zap
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		logLevel := "info"
		if statusCode >= 400 {
			logLevel = "warn"
		}
		if statusCode >= 500 {
			logLevel = "error"
		}

		logFields := []interface{}{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
			"ip", c.ClientIP(),
			"user_id", uid,
		}

		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		switch logLevel {
		case "warn":
			log.W(log.Sprintf("[WARN] %s %s - status=%d, duration=%dms", c.Request.Method, c.Request.URL.Path, statusCode, duration.Milliseconds()))
		case "error":
			errMsg := ""
			if len(c.Errors) > 0 {
				errMsg = c.Errors.String()
			}
			log.E(log.Sprintf("[ERROR] %s %s - status=%d, duration=%dms, errors=%s", c.Request.Method, c.Request.URL.Path, statusCode, duration.Milliseconds(), errMsg))
		default:
			log.I(log.Sprintf("[%d] %s %s - %dms", statusCode, c.Request.Method, c.Request.URL.Path, duration.Milliseconds()))
		}

		// 异步写入数据库（只记录关键操作，不记录请求/响应体）
		go func() {
			accessLog := &models.AccessLog{
				UserID:     uid,
				Username:   getUsername(uid),
				Method:     c.Request.Method,
				Path:       c.Request.URL.Path,
				Query:      c.Request.URL.RawQuery,
				IP:         c.ClientIP(),
				UserAgent:  c.Request.UserAgent(),
				StatusCode: statusCode,
				Duration:   duration.Milliseconds(),
			}
			if len(c.Errors) > 0 {
				accessLog.ErrorMsg = c.Errors.String()
			}
			service.CreateAccessLog(accessLog)
		}()
	}
}

// ============ 辅助函数 ============

// getUserIDFromRequest 从请求头直接解析token获取用户ID
func getUserIDFromRequest(c *gin.Context) uint {
	token := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "bearer ") {
		token = authHeader[7:]
	}
	if token == "" {
		token = c.GetHeader("token")
	}
	if token == "" {
		return 0
	}
	claims, err := utils.ParseToken(token)
	if err != nil {
		return 0
	}
	return claims.UserID
}

func getUsername(userID uint) string {
	if userID == 0 {
		return "未登录"
	}
	// 从Redis缓存获取用户名
	username, err := utils.GetUsernameByCache(userID)
	if err != nil {
		// 缓存未命中，从数据库查询
		user, err := service.GetUserByID(int(userID))
		if err != nil {
			return ""
		}
		return user.Username
	}
	return username
}
