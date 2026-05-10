package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidationError 校验错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RegisterCustomValidators 注册自定义校验器
func RegisterCustomValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册手机号校验
		v.RegisterValidation("phone", validatePhone)

		// 注册密码强度校验（至少8位，包含字母和数字）
		v.RegisterValidation("password_strong", validatePasswordStrong)

		// 注册IP地址校验
		v.RegisterValidation("ip", validateIP)

		// 注册URL校验
		v.RegisterValidation("url", validateURL)

		// 注册身份证校验
		v.RegisterValidation("idcard", validateIDCard)

		// 注册自定义错误消息
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	return nil
}

// validatePhone 手机号校验（简单版）
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) == 11 && phone[0] == '1' {
		return true
	}
	return false
}

// validatePasswordStrong 强密码校验
func validatePasswordStrong(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 8 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, c := range pwd {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}

// validateIP IP地址校验
func validateIP(fl validator.FieldLevel) bool {
	ip := fl.Field().String()
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		var num int
		fmt.Sscanf(part, "%d", &num)
		if num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// validateURL URL校验
func validateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// validateIDCard 身份证校验（简单校验长度和格式）
func validateIDCard(fl validator.FieldLevel) bool {
	idcard := fl.Field().String()
	if len(idcard) != 18 {
		return false
	}
	// 检查前17位是否为数字
	for i := 0; i < 17; i++ {
		if idcard[i] < '0' || idcard[i] > '9' {
			return false
		}
	}
	// 检查最后一位是否为数字或X
	lastChar := idcard[17]
	if !((lastChar >= '0' && lastChar <= '9') || lastChar == 'X') {
		return false
	}
	return true
}

// ValidateRequest 校验请求参数，返回格式化错误
func ValidateRequest(obj interface{}) []ValidationError {
	var errors []ValidationError

	if err := binding.Validator.ValidateStruct(obj); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s不能为空", field)
			case "email":
				message = fmt.Sprintf("%s必须是有效的邮箱", field)
			case "min":
				message = fmt.Sprintf("%s长度不能少于%s", field, err.Param())
			case "max":
				message = fmt.Sprintf("%s长度不能超过%s", field, err.Param())
			case "phone":
				message = fmt.Sprintf("%s必须是有效的手机号", field)
			case "password_strong":
				message = fmt.Sprintf("%s至少8位，包含字母和数字", field)
			case "len":
				message = fmt.Sprintf("%s长度必须为%s", field, err.Param())
			case "gte":
				message = fmt.Sprintf("%s必须大于等于%s", field, err.Param())
			case "lte":
				message = fmt.Sprintf("%s必须小于等于%s", field, err.Param())
			default:
				message = fmt.Sprintf("%s校验失败", field)
			}

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}
	return errors
}

// BindAndValidate 绑定并校验参数
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		Error(c, "参数错误: "+err.Error())
		return false
	}
	errors := ValidateRequest(obj)
	if len(errors) > 0 {
		Error(c, errors[0].Message)
		return false
	}
	return true
}
