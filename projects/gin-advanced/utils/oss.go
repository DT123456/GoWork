package utils

import (
	"bytes"
	"fmt"
	"gin-advanced/core"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSConfig OSS配置
type OSSConfig struct {
	Enabled   bool   // 是否启用OSS
	Endpoint  string // OSS endpoint，如：oss-cn-beijing.aliyuncs.com
	AccessKey string // AccessKey ID
	SecretKey string // AccessKey Secret
	Bucket    string // Bucket名称
	Domain    string // 自定义域名，如：https://cdn.example.com（可选）
	SavePath  string // 上传路径前缀，如：images/
}

// 获取OSS配置（从配置文件读取）
func GetOSSConfig() OSSConfig {
	cfg := core.GetConfig()
	return OSSConfig{
		Enabled:   cfg.OSS.Enabled,
		Endpoint:  cfg.OSS.Endpoint,
		AccessKey: cfg.OSS.AccessKey,
		SecretKey: cfg.OSS.SecretKey,
		Bucket:    cfg.OSS.Bucket,
		Domain:    cfg.OSS.Domain,
		SavePath:  cfg.OSS.SavePath,
	}
}

var ossClient *oss.Client
var ossBucket *oss.Bucket

// InitOSS 初始化OSS客户端
func InitOSS() error {
	cfg := GetOSSConfig()
	if !cfg.Enabled {
		return nil
	}

	var err error
	ossClient, err = oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("初始化OSS客户端失败: %w", err)
	}

	ossBucket, err = ossClient.Bucket(cfg.Bucket)
	if err != nil {
		return fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	return nil
}

// OSSUploadResult OSS上传结果
type OSSUploadResult struct {
	URL       string `json:"url"`       // 完整URL
	FileName  string `json:"filename"` // 文件名
	FileSize  int64  `json:"filesize"` // 文件大小
	OSSKey    string `json:"oss_key"`  // OSS对象Key
	IsOSS     bool   `json:"is_oss"`   // 是否上传到OSS
}

// UploadToOSS 上传文件到OSS
func UploadToOSS(file *multipart.FileHeader, objectKey string) (*OSSUploadResult, error) {
	cfg := GetOSSConfig()
	if !cfg.Enabled {
		return nil, fmt.Errorf("OSS未启用")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 上传到OSS
	err = ossBucket.PutObject(objectKey, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("上传到OSS失败: %w", err)
	}

	// 生成URL
	url := buildOSSURL(objectKey)

	return &OSSUploadResult{
		URL:      url,
		FileName: file.Filename,
		FileSize: file.Size,
		OSSKey:   objectKey,
		IsOSS:    true,
	}, nil
}

// UploadImageToOSS 上传图片到OSS
func UploadImageToOSS(file *multipart.FileHeader) (*OSSUploadResult, error) {
	ext := getFileExt(file.Filename)
	objectKey := fmt.Sprintf("%simages/%d/%s%s",
		GetOSSConfig().SavePath,
		time.Now().Unix(),
		generateFileName(),
		ext,
	)
	return UploadToOSS(file, objectKey)
}

// UploadAvatarToOSS 上传头像到OSS
func UploadAvatarToOSS(file *multipart.FileHeader) (*OSSUploadResult, error) {
	ext := getFileExt(file.Filename)
	objectKey := fmt.Sprintf("%savatars/%d/%s%s",
		GetOSSConfig().SavePath,
		time.Now().Unix(),
		generateFileName(),
		ext,
	)
	return UploadToOSS(file, objectKey)
}

// UploadFileToOSS 上传任意文件到OSS
func UploadFileToOSS(file *multipart.FileHeader) (*OSSUploadResult, error) {
	ext := getFileExt(file.Filename)
	objectKey := fmt.Sprintf("%sfiles/%d/%s%s",
		GetOSSConfig().SavePath,
		time.Now().Unix(),
		generateFileName(),
		ext,
	)
	return UploadToOSS(file, objectKey)
}

// DeleteFromOSS 从OSS删除文件
func DeleteFromOSS(objectKey string) error {
	cfg := GetOSSConfig()
	if !cfg.Enabled {
		return nil
	}
	return ossBucket.DeleteObject(objectKey)
}

// buildOSSURL 构建OSS文件访问URL
func buildOSSURL(objectKey string) string {
	cfg := GetOSSConfig()
	if cfg.Domain != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.Domain, "/"), objectKey)
	}
	return fmt.Sprintf("https://%s.%s/%s", cfg.Bucket, cfg.Endpoint, objectKey)
}

// getFileExt 获取文件扩展名
func getFileExt(filename string) string {
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		ext = filename[idx:]
	}
	return strings.ToLower(ext)
}

// generateFileName 生成唯一文件名
func generateFileName() string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), randString(8))
}

// randString 生成随机字符串
func randString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}

// IsOSSEnabled 检查OSS是否启用
func IsOSSEnabled() bool {
	return GetOSSConfig().Enabled
}
