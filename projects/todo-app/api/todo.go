package api

import (
	"net/http"
	"todo-app/models"

	"github.com/gin-gonic/gin"
)

func CreateTodoHandler(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := todo.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func GetTodosHandler(c *gin.Context) {
	var todos []models.Todo
	if err := models.FindAll(&todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := models.GetTodoByID(id, &todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func UpdateTodoHandler(c *gin.Context) {
	id := c.Param("id")

	// 先查询现有记录，确保我们有正确的 ID 和数据库中的实际数据
	var todo models.Todo
	if err := models.GetTodoByID(id, &todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// 绑定新值
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新并返回数据库中的实际数据
	if err := todo.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新查询确保返回数据库实际值
	models.GetTodoByID(id, &todo)
	c.JSON(http.StatusOK, todo)
}

func DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteTodoByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}