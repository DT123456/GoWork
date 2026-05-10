package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRateLimitConfig(t *testing.T) {
	config := RateLimitConfig{
		MaxCount: 10,
		Window:   60,
		Type:     "ip",
	}

	if config.MaxCount != 10 {
		t.Errorf("Expected MaxCount 10, got %d", config.MaxCount)
	}
	if config.Type != "ip" {
		t.Errorf("Expected Type 'ip', got '%s'", config.Type)
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	if DefaultRateLimitConfig.MaxCount != 100 {
		t.Errorf("Expected default MaxCount 100, got %d", DefaultRateLimitConfig.MaxCount)
	}
}

func TestErrorCode(t *testing.T) {
	// 测试错误码定义
	tests := []struct {
		err     ErrorCode
		code    int
		message string
	}{
		{ErrBadRequest, 400, "请求参数错误"},
		{ErrUnauthorized, 401, "未授权，请登录"},
		{ErrForbidden, 403, "禁止访问"},
		{ErrNotFound, 404, "资源不存在"},
		{ErrServerError, 500, "服务器内部错误"},
		{ErrTooManyRequests, 429, "请求过于频繁"},
	}

	for _, tt := range tests {
		if tt.err.Code != tt.code {
			t.Errorf("Expected code %d for %s, got %d", tt.code, tt.message, tt.err.Code)
		}
		if tt.err.Message != tt.message {
			t.Errorf("Expected message '%s', got '%s'", tt.message, tt.err.Message)
		}
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRecoveryNoPanic(t *testing.T) {
	r := gin.New()
	r.Use(Recovery())
	r.GET("/normal", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/normal", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
