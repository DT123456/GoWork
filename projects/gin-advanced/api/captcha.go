package api

import (
	"gin-advanced/utils"

	"github.com/gin-gonic/gin"
)

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	result, err := utils.GetCaptcha()
	if err != nil {
		utils.Error(c, "生成验证码失败")
		return
	}
	utils.Success(c, result)
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID string `json:"captcha_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "参数错误")
		return
	}

	if utils.VerifyCaptcha(req.CaptchaID, req.Code) {
		utils.SuccessWithMessage(c, "验证成功", true)
	} else {
		utils.Error(c, "验证码错误")
	}
}
