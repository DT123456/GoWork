package main

import (
	"fmt"
	"log"

	"todo-app/config"
	"todo-app/models"
	"todo-app/router"
)

func main() {

	config.InitConfig()

	// 初始化数据库连接
	if err := config.InitDB(config.Conf); err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}

	defer config.Close() // 程序退出时关闭数据库连接

	// 自动迁移表结构
	if err := config.DB.AutoMigrate(&models.Todo{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server startup failed: %v", err)
	}
}