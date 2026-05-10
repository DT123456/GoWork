package utils

import (
	"math/rand"
	"time"

	"github.com/mojocn/base64Captcha"
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Width  int    // 宽度
	Height int    // 高度
	Length int    // 验证码长度
	Chars  string // 可用字符集
	Expire int64  // 过期时间(秒)
}

// 默认配置
var DefaultCaptchaConfig = CaptchaConfig{
	Width:  120,
	Height: 40,
	Length: 4,
	Chars:  "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxy2345678",
	Expire: 300, // 5分钟
}

// CaptchaResult 验证码结果
type CaptchaResult struct {
	CaptchaID string `json:"captcha_id"` // 验证码ID
	Image     string `json:"image"`      // Base64编码的图片
}

// 初始化验证码驱动
var store = base64Captcha.DefaultMemStore

// GetCaptcha 生成验证码，返回ID和Base64图片
func GetCaptcha() (*CaptchaResult, error) {
	cfg := DefaultCaptchaConfig
	return GenerateCaptcha(cfg)
}

// GenerateCaptcha 生成自定义配置的验证码
func GenerateCaptcha(cfg CaptchaConfig) (*CaptchaResult, error) {
	// 使用 base64Captcha 生成图片
	driver := base64Captcha.NewDriverDigit(int(cfg.Height), int(cfg.Width), int(cfg.Length), 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaResult{
		CaptchaID: id,
		Image:     b64s,
	}, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaID, code string) bool {
	if captchaID == "" || code == "" {
		return false
	}

	// 使用 store 验证
	result := store.Verify(captchaID, code, true)
	return result
}

// generateCode 生成随机验证码（备用）
func generateCode(length int, chars string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return string(code)
}
