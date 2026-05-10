package middleware

import (
	"fmt"
	"gin-advanced/core"
	"gin-advanced/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 限流阈值（次数）
	MaxCount int
	// 限流窗口（秒）
	Window time.Duration
	// 限流类型: "ip" / "user" / "path"
	Type string
}

// DefaultRateLimitConfig 默认限流配置（每IP每分钟100次）
var DefaultRateLimitConfig = RateLimitConfig{
	MaxCount: 100,
	Window:   time.Minute,
	Type:     "ip",
}

// RateLimit 限流中间件
func RateLimit(config ...RateLimitConfig) gin.HandlerFunc {
	cfg := DefaultRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		key := getRateLimitKey(c, cfg.Type)
		count, err := core.RDB.Incr(core.Ctx, key).Uint64()
		if err != nil {
			c.Next()
			return
		}

		// 设置过期时间
		if count == 1 {
			core.RDB.Expire(core.Ctx, key, cfg.Window)
		}

		// 获取剩余时间
		ttl, _ := core.RDB.TTL(core.Ctx, key).Result()
		remain := int64(cfg.MaxCount) - int64(count)
		if remain < 0 {
			remain = 0
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.MaxCount))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remain))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

		if count > uint64(cfg.MaxCount) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
				"data": gin.H{
					"retry_after": int(ttl.Seconds()),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getRateLimitKey 生成限流key
func getRateLimitKey(c *gin.Context, limitType string) string {
	switch limitType {
	case "ip":
		return fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
	case "user":
		userID, _ := c.Get("userID")
		if id, ok := userID.(uint); ok {
			return fmt.Sprintf("ratelimit:user:%d", id)
		}
		return fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
	case "path":
		return fmt.Sprintf("ratelimit:path:%s:%s", c.Request.URL.Path, c.ClientIP())
	default:
		return fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
	}
}

// RateLimitByIP 基于IP的限流
func RateLimitByIP() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: 60,
		Window:   time.Minute,
		Type:     "ip",
	})
}

// RateLimitByUser 基于用户的限流
func RateLimitByUser() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: 1000,
		Window:   time.Hour,
		Type:     "user",
	})
}

// RateLimitByPath 基于路径的限流
func RateLimitByPath() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: 100,
		Window:   time.Minute,
		Type:     "path",
	})
}

// RateLimitLogin 登录限流（防爆破）
func RateLimitLogin() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: 5,
		Window:   time.Minute,
		Type:     "ip",
	})
}

// RateLimitAPI API通用限流
func RateLimitAPI() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: 100,
		Window:   time.Minute,
		Type:     "ip",
	})
}

// BurstRateLimit 突发限流（允许短暂的高并发）
func BurstRateLimit(maxCount int, window time.Duration) gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		MaxCount: maxCount,
		Window:   window,
		Type:     "ip",
	})
}

// WhiteListRateLimit 白名单限流（跳过检查）
func WhiteListRateLimit(whiteList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		for _, w := range whiteList {
			if ip == w {
				c.Next()
				return
			}
		}
		RateLimitByIP()(c)
	}
}

// MultiRateLimit 多维度限流（同时限制IP和用户）
func MultiRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipKey := fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
		userID, _ := c.Get("userID")

		// IP限流
		ipCount, _ := core.RDB.Incr(core.Ctx, ipKey).Uint64()
		if ipCount == 1 {
			core.RDB.Expire(core.Ctx, ipKey, time.Minute)
		}
		if ipCount > 60 {
			utils.ErrorWithCode(c, 429, "请求过于频繁")
			c.Abort()
			return
		}

		// 用户限流
		if id, ok := userID.(uint); ok {
			userKey := fmt.Sprintf("ratelimit:user:%d", id)
			userCount, _ := core.RDB.Incr(core.Ctx, userKey).Uint64()
			if userCount == 1 {
				core.RDB.Expire(core.Ctx, userKey, time.Minute)
			}
			if userCount > 100 {
				utils.ErrorWithCode(c, 429, "请求过于频繁")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GetRateLimitStatus 获取限流状态
func GetRateLimitStatus(c *gin.Context) {
	key := getRateLimitKey(c, "ip")
	count, _ := core.RDB.Get(core.Ctx, key).Uint64()
	ttl, _ := core.RDB.TTL(core.Ctx, key).Result()

	utils.Success(c, gin.H{
		"current_count": count,
		"remain":         60 - int(count),
		"reset_in":       int(ttl.Seconds()),
	})
}
