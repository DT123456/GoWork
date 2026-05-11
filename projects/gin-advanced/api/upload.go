package api

import (
	"gin-advanced/log"
	"gin-advanced/utils"

	"github.com/gin-gonic/gin"
)

const (
	// MaxFileSize 文件最大大小 (10MB)
	MaxFileSize = 10 * 1024 * 1024
	// MaxImageSize 图片最大大小 (5MB)
	MaxImageSize = 5 * 1024 * 1024
)

// UploadFile
// @Summary 上传文件
// @Description 上传通用文件到服务器或OSS
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /upload [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 检查文件大小
	if file.Size > MaxFileSize {
		utils.BadRequest(c, "文件大小不能超过10MB")
		return
	}

	var result interface{}

	// 判断是否启用OSS
	if utils.IsOSSEnabled() {
		ossResult, err := utils.UploadFileToOSS(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = ossResult
	} else {
		localResult, err := utils.UploadFile(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = localResult
	}

	utils.Success(c, result)
}

// UploadImage
// @Summary 上传图片
// @Description 上传图片到服务器或OSS
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "图片文件"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /upload/image [post]
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的图片")
		return
	}

	// 检查文件大小
	if file.Size > MaxImageSize {
		utils.BadRequest(c, "图片大小不能超过5MB")
		return
	}

	var result interface{}

	// 判断是否启用OSS
	if utils.IsOSSEnabled() {
		ossResult, err := utils.UploadImageToOSS(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = ossResult
	} else {
		localResult, err := utils.UploadImage(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = localResult
	}

	utils.Success(c, result)
}

// UploadAvatar
// @Summary 上传头像
// @Description 上传用户头像图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "头像文件"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /upload/avatar [post]
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的头像")
		return
	}

	// 限制图片大小
	if file.Size > 2*1024*1024 {
		utils.BadRequest(c, "头像大小不能超过2MB")
		return
	}

	var result interface{}

	// 判断是否启用OSS
	if utils.IsOSSEnabled() {
		ossResult, err := utils.UploadAvatarToOSS(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = ossResult
	} else {
		localResult, err := utils.UploadAvatar(file)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		result = localResult
	}

	utils.Success(c, result)
}

// UploadFiles
// @Summary 批量上传文件
// @Description 一次上传多个文件
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /upload/multiple [post]
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "获取文件失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	var results []interface{}
	var failedCount int

	for _, f := range files {
		// 检查每个文件的大小
		if f.Size > MaxFileSize {
			failedCount++
			log.W(log.Sprintf("批量上传跳过: 文件 %s 大小超过限制 (%d bytes)", f.Filename, f.Size))
			continue
		}

		if utils.IsOSSEnabled() {
			ossResult, err := utils.UploadFileToOSS(f)
			if err != nil {
				failedCount++
				log.E(log.Sprintf("OSS上传失败 [%s]: %v", f.Filename, err))
				continue
			}
			results = append(results, ossResult)
		} else {
			localResult, err := utils.UploadFile(f)
			if err != nil {
				failedCount++
				log.E(log.Sprintf("本地上传失败 [%s]: %v", f.Filename, err))
				continue
			}
			results = append(results, localResult)
		}
	}

	// 如果全部失败，返回错误
	if len(results) == 0 && failedCount > 0 {
		utils.Error(c, "所有文件上传失败")
		return
	}

	// 如果有部分失败，返回成功但带警告信息
	if failedCount > 0 {
		log.W(log.Sprintf("批量上传完成: 成功 %d 个, 失败 %d 个", len(results), failedCount))
	}

	utils.Success(c, gin.H{
		"success_count": len(results),
		"failed_count":  failedCount,
		"files":         results,
	})
}

// UploadConfigResponse 上传配置响应
type UploadConfigResponse struct {
	OSSEnabled    bool     `json:"oss_enabled"`
	MaxSize       int      `json:"max_size"`
	MaxImageSize  int      `json:"max_image_size"`
	MaxAvatarSize int      `json:"max_avatar_size"`
	AllowedExt    []string `json:"allowed_ext"`
	AllowedImage  []string `json:"allowed_image"`
	AllowedAvatar []string `json:"allowed_avatar"`
}

// GetUploadConfig
// @Summary 获取上传配置
// @Description 获取文件上传的配置信息，包括大小限制和允许的文件类型
// @Tags 文件上传
// @Produce json
// @Success 200 {object} utils.Response{data=UploadConfigResponse}
// @Router /upload/config [get]
func GetUploadConfig(c *gin.Context) {
	config := UploadConfigResponse{
		OSSEnabled:    utils.IsOSSEnabled(),
		MaxSize:       10,
		MaxImageSize:  5,
		MaxAvatarSize: 2,
		AllowedExt:    []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx"},
		AllowedImage:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		AllowedAvatar: []string{".jpg", ".jpeg", ".png", ".webp"},
	}
	utils.Success(c, config)
}
