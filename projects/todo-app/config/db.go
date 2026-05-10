package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Conf = new(Config)

type Config struct {
	MySQLConfig struct {
		User     string `ini:"user"`
		Password string `ini:"password"`
		DBName   string `ini:"dbname"`
		Host     string `ini:"host"`
		Port     string `ini:"port"`
	} `ini:"mysql"`
}

func InitConfig() {
	Conf.MySQLConfig.User = "root"
	Conf.MySQLConfig.Password = "root123"
	Conf.MySQLConfig.Host = "127.0.0.1"
	Conf.MySQLConfig.Port = "3306"
	Conf.MySQLConfig.DBName = "bubble"
}

func InitDB(cfg *Config) error {
	// 初始化日志文件
	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer f.Close() // 确保文件在函数结束时关闭
		log.SetOutput(f) // 将日志输出设置为文件
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLConfig.User,
		cfg.MySQLConfig.Password,
		cfg.MySQLConfig.Host,
		cfg.MySQLConfig.Port,
		cfg.MySQLConfig.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 连接数据库失败，记录错误日志并退出程序
		log.Fatal("failed to connect database: ", err)
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	// 防止 DB == nil 时 panic
	if DB != nil {
		// 获取底层的 *sql.DB 对象以关闭连接
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}