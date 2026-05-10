package models

import (
	"strconv"

	"gorm.io/gorm"

	"todo-app/config"
)

type Todo struct {
	gorm.Model // 自带 ID、CreatedAt、UpdatedAt、DeletedAt
	//Id          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Completed   bool   `gorm:"default:false" json:"completed"`
}

func (Todo) TableName() string {
	return "todos1"
}

// Create 创建新记录
func (t *Todo) Create() error {
	return config.DB.Create(t).Error
}

// FindAll 查询所有记录
func FindAll(todos *[]Todo) error {
	//GORM在查询时添加 WHERE deleted_at IS NULL 条件
	return config.DB.Find(todos).Error
}

// GetTodoByID 根据ID查询记录
func GetTodoByID(id string, todo *Todo) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	todo.ID = uint(intID)
	return config.DB.First(todo).Error
}

// Update 更新记录
func (t *Todo) Update() error {
	return config.DB.Save(t).Error
}

// DeleteTodoByID 根据ID删除记录
func DeleteTodoByID(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return config.DB.Delete(&Todo{}, intID).Error
}