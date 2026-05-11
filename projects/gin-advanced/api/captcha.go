package api

import (
	"gin-advanced/utils"

	"github.com/gin-gonic/gin"
)

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	Image     string `json:"image"`
}

// VerifyCaptchaRequest 验证验证码请求
type VerifyCaptchaRequest struct {
	CaptchaID string `json:"captcha_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// GetCaptcha
// @Summary 获取图形验证码
// @Description 获取验证码图片，返回验证码ID和Base64编码的图片
// @Tags 验证码
// @Produce json
// @Success 200 {object} utils.Response{data=CaptchaResponse}
// @Failure 500 {object} utils.Response
// @Router /captcha [get]
func GetCaptcha(c *gin.Context) {
	result, err := utils.GetCaptcha()
	if err != nil {
		utils.Error(c, "生成验证码失败")
		return
	}
	utils.Success(c, result)
}

// VerifyCaptcha
// @Summary 验证验证码
// @Description 验证用户输入的验证码是否正确
// @Tags 验证码
// @Accept json
// @Produce json
// @Param body body VerifyCaptchaRequest true "验证码信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /captcha/verify [post]
func VerifyCaptcha(c *gin.Context) {
	var req VerifyCaptchaRequest

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
