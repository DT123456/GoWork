package utils

import (
	"fmt"
	"gin-advanced/core"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize    int64         // 最大文件大小（MB）
	AllowedExt []string      // 允许的扩展名
	SavePath   string        // 保存路径
	BaseURL    string        // 访问基础URL
}

// 默认上传配置
var DefaultUploadConfig = UploadConfig{
	MaxSize:    10, // 10MB
	AllowedExt: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx"},
	SavePath:   "./uploads",
	BaseURL:    "/uploads",
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`       // 访问URL
	FileName string `json:"filename"` // 文件名
	FileSize int64  `json:"filesize"` // 文件大小
}

// UploadFile 上传单个文件
func UploadFile(file *multipart.FileHeader, config ...UploadConfig) (*UploadResult, error) {
	cfg := DefaultUploadConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 检查文件大小
	if file.Size > cfg.MaxSize*1024*1024 {
		return nil, fmt.Errorf("文件大小超过限制(最大%dMB)", cfg.MaxSize)
	}

	// 检查扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, e := range cfg.AllowedExt {
		if e == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 创建保存目录
	dateDir := time.Now().Format("2006/01/02")
	saveDir := filepath.Join(cfg.SavePath, dateDir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 生成唯一文件名
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(saveDir, newFileName)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 返回结果
	url := fmt.Sprintf("%s/%s/%s", cfg.BaseURL, dateDir, newFileName)

	return &UploadResult{
		URL:      url,
		FileName: newFileName,
		FileSize: file.Size,
	}, nil
}

// UploadImage 上传图片（限制图片格式）
func UploadImage(file *multipart.FileHeader) (*UploadResult, error) {
	config := UploadConfig{
		MaxSize:    5, // 图片默认5MB
		AllowedExt: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		SavePath:   "./uploads/images",
		BaseURL:    "/uploads/images",
	}
	return UploadFile(file, config)
}

// UploadAvatar 上传头像（限制图片格式和大小）
func UploadAvatar(file *multipart.FileHeader) (*UploadResult, error) {
	config := UploadConfig{
		MaxSize:    2, // 头像2MB
		AllowedExt: []string{".jpg", ".jpeg", ".png", ".webp"},
		SavePath:   "./uploads/avatars",
		BaseURL:    "/uploads/avatars",
	}
	return UploadFile(file, config)
}

// UploadFiles 上传多个文件
func UploadFiles(files []*multipart.FileHeader, config ...UploadConfig) ([]*UploadResult, error) {
	results := make([]*UploadResult, 0, len(files))
	for _, f := range files {
		result, err := UploadFile(f, config...)
		if err != nil {
			continue // 跳过失败的文件
		}
		results = append(results, result)
	}
	return results, nil
}

// DeleteFile 删除文件
func DeleteFile(url string) error {
	// 从URL提取文件路径
	path := strings.TrimPrefix(url, "/uploads/")
	fullPath := filepath.Join("./uploads", path)

	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

// GenerateFileKey 生成文件访问Key（用于验证）
func GenerateFileKey(fileName, userID string) string {
	return fmt.Sprintf("file:%s:%s:%d", fileName, uuid.New().String(), time.Now().Unix())
}

// CheckFileAccess 检查文件访问权限
func CheckFileAccess(fileKey string) bool {
	key := "file_access:" + fileKey
	result, err := core.RDB.Exists(core.Ctx, key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

// SetFileAccess 设置文件访问权限
func SetFileAccess(fileKey string, userID uint, expire time.Duration) error {
	key := "file_access:" + fileKey
	return core.RDB.Set(core.Ctx, key, userID, expire).Err()
}
