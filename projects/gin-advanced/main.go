package main

import (
	"gin-advanced/core"
	"gin-advanced/crontab"
	"gin-advanced/models"
	"gin-advanced/routes"
	"gin-advanced/utils"
)

func main() {
	// 初始化
	core.InitViper()
	core.InitZap()
	core.InitMySQL()
	core.InitRedis()

	// 注册自定义校验器
	utils.RegisterCustomValidators()

	// 自动建表
	core.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.SystemLog{},
		&models.AccessLog{},
	)

	// 启动定时任务
	crontab.Start()

	// 启动
	r := routes.InitRouter()
	r.Run(":8080")
}
