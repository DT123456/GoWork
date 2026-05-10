package service

import (
	"errors"

	"gin-advanced/core"
	"gin-advanced/models"
	"gin-advanced/utils"
)

func CreateUser(user *models.User) error {
	// 检查用户名是否已存在
	var count int64
	core.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 密码加密存储
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return core.DB.Create(user).Error
}

func GetUserByUsername(username string) (user models.User, err error) {
	err = core.DB.Where("username = ?", username).First(&user).Error
	return
}

func GetUserByID(id int) (user models.User, err error) {
	err = core.DB.Where("id = ?", id).First(&user).Error
	return
}

func UpdateUser(id int, user *models.User) error {
	return core.DB.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func DeleteUser(id int) error {
	return core.DB.Delete(&models.User{}, id).Error
}

func GetUserPermissions(userID uint) ([]string, error) {
	var permissions []string
	err := core.DB.Raw(`
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?`, userID).Scan(&permissions).Error
	return permissions, err
}