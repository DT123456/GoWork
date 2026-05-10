package core

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	user := viper.GetString("mysql.username")
	pwd := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")

	dsn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		ZapMySQL.Fatal("数据库连接失败",
			zap.String("host", host),
			zap.String("port", port),
			zap.String("database", dbname),
			zap.Error(err),
		)
	}
	DB = db

	ZapMySQL.Info("数据库连接成功",
		zap.String("host", host),
		zap.String("port", port),
		zap.String("database", dbname),
	)
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
