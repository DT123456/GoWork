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

// UploadFile 上传文件
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

// UploadImage 上传图片
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

// UploadAvatar 上传头像
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

// UploadFiles 批量上传文件
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

// GetUploadConfig 获取上传配置
func GetUploadConfig(c *gin.Context) {
	config := map[string]interface{}{
		"oss_enabled":    utils.IsOSSEnabled(),
		"max_size":       10,
		"max_image_size": 5,
		"max_avatar_size": 2,
		"allowed_ext":    []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx"},
		"allowed_image":  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		"allowed_avatar": []string{".jpg", ".jpeg", ".png", ".webp"},
	}
	utils.Success(c, config)
}
