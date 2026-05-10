package models

import "gorm.io/gorm"

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
