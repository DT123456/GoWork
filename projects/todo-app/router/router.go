package router

import (
	"todo-app/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 定义路由
	r.POST("/todos", api.CreateTodoHandler)
	r.GET("/todos", api.GetTodosHandler)
	r.GET("/todos/:id", api.GetTodoHandler)
	r.PUT("/todos/:id", api.UpdateTodoHandler)
	r.DELETE("/todos/:id", api.DeleteTodoHandler)

	return r
}