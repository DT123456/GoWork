package core

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 配置结构
type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	MySQL struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"mysql"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret string `mapstructure:"secret"`
		Expire int    `mapstructure:"expire"`
	} `mapstructure:"jwt"`
	Log struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	} `mapstructure:"log"`
	OSS struct {
		Enabled   bool   `mapstructure:"enabled"`
		Endpoint  string `mapstructure:"endpoint"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Bucket    string `mapstructure:"bucket"`
		Domain    string `mapstructure:"domain"`
		SavePath  string `mapstructure:"save_path"`
	} `mapstructure:"oss"`
}

var (
	AppConfig *Config
	configMux sync.RWMutex
)

func InitViper() {
	viper.SetConfigFile("config/config.yml")
	viper.SetConfigType("yml")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("配置文件读取失败:", err)
	}

	// 映射到结构体
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatal("配置映射失败:", err)
	}

	// 监听配置文件变化（热更新）
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件已更改，重新加载...")
		if err := viper.Unmarshal(AppConfig); err != nil {
			log.Println("配置重载失败:", err)
		}
	})
}

// GetConfig 获取当前配置（线程安全）
func GetConfig() *Config {
	configMux.RLock()
	defer configMux.RUnlock()
	return AppConfig
}

// GetString 获取字符串配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取Int配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool 获取Bool配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}
