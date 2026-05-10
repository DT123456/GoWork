package service

import (
	"gin-advanced/core"
	"gin-advanced/log"
	"gin-advanced/models"
	"gin-advanced/utils"
)

// CreateRole 创建角色
func CreateRole(role *models.Role) error {
	return core.DB.Create(role).Error
}

// GetRoles 获取所有角色
func GetRoles() ([]models.Role, error) {
	var roles []models.Role
	err := core.DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

// GetRoleByID 根据ID获取角色
func GetRoleByID(id uint) (models.Role, error) {
	var role models.Role
	err := core.DB.Preload("Permissions").First(&role, id).Error
	return role, err
}

// UpdateRole 更新角色
func UpdateRole(role *models.Role) error {
	return core.DB.Save(role).Error
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	return core.DB.Delete(&models.Role{}, id).Error
}

// AssignRoleToUser 给用户分配角色
func AssignRoleToUser(userID, roleID uint) error {
	// 先检查是否已存在
	var count int64
	core.DB.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	if count > 0 {
		return nil // 已存在，直接返回成功
	}
	return core.DB.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error
}

// RemoveRoleFromUser 移除用户角色
func RemoveRoleFromUser(userID, roleID uint) error {
	return core.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := core.DB.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error
	return roles, err
}

// AddPermissionToRole 给角色添加权限
func AddPermissionToRole(roleID, permissionID uint) error {
	var role models.Role
	if err := core.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	var perm models.Permission
	if err := core.DB.First(&perm, permissionID).Error; err != nil {
		return err
	}

	return core.DB.Model(&role).Association("Permissions").Append(&perm)
}

// SetRolePermissions 批量设置角色权限（替换模式）
func SetRolePermissions(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := core.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	// 先清空现有权限
	core.DB.Model(&role).Association("Permissions").Clear()

	// 再添加新权限
	if len(permissionIDs) == 0 {
		return nil
	}

	var permissions []models.Permission
	if err := core.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	return core.DB.Model(&role).Association("Permissions").Append(&permissions)
}

// RemovePermissionFromRole 移除角色权限
func RemovePermissionFromRole(roleID, permissionID uint) error {
	var role models.Role
	if err := core.DB.First(&role, roleID).Error; err != nil {
		return err
	}
	var perm models.Permission
	if err := core.DB.First(&perm, permissionID).Error; err != nil {
		return err
	}
	return core.DB.Model(&role).Association("Permissions").Delete(&perm)
}

// GetPermissions 获取所有权限（包括已删除的）
func GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := core.DB.Unscoped().Find(&permissions).Error
	return permissions, err
}

// CreatePermission 创建权限
func CreatePermission(perm *models.Permission) error {
	return core.DB.Create(perm).Error
}

// UpdatePermission 更新权限
func UpdatePermission(perm *models.Permission) error {
	return core.DB.Save(perm).Error
}

// DeletePermission 删除权限
func DeletePermission(id uint) error {
	return core.DB.Delete(&models.Permission{}, id).Error
}

// SeedRoles 初始化默认角色和权限
func SeedRoles() error {
	permissions := []models.Permission{
		{Name: "user:create", Path: "/register", Method: "POST", Description: "创建用户"},
		{Name: "user:read", Path: "/user/*", Method: "GET", Description: "读取用户"},
		{Name: "user:update", Path: "/user/*", Method: "PUT", Description: "更新用户"},
		{Name: "user:delete", Path: "/user/*", Method: "DELETE", Description: "删除用户"},
		{Name: "user:manage", Path: "/admin/*", Method: "*", Description: "管理所有用户"},
		{Name: "role:manage", Path: "/admin/roles/*", Method: "*", Description: "管理角色"},
		{Name: "log:view", Path: "/admin/logs/*", Method: "*", Description: "查看日志"},
	}

	for i := range permissions {
		core.DB.FirstOrCreate(&permissions[i], models.Permission{Name: permissions[i].Name})
	}

	// 创建 admin 角色
	admin := models.Role{Name: "admin", Description: "管理员", Permissions: permissions}
	core.DB.FirstOrCreate(&admin, models.Role{Name: "admin"})

	// 创建普通用户角色
	user := models.Role{Name: "user", Description: "普通用户"}
	core.DB.FirstOrCreate(&user, models.Role{Name: "user"})

	// 创建默认管理员用户
	var adminUser models.User
	core.DB.Where("username = ?", "admin").First(&adminUser)
	if adminUser.ID == 0 {
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			log.E(log.Sprintf("创建默认管理员密码加密失败: %v", err))
			return err
		}
		adminUser = models.User{
			Username: "admin",
			Password: hashedPassword,
			Nickname: "管理员",
		}
		if err := core.DB.Create(&adminUser).Error; err != nil {
			log.E(log.Sprintf("创建默认管理员用户失败: %v", err))
			return err
		}
		log.I(log.Sprintf("创建默认管理员用户成功: username=admin, user_id=%d", adminUser.ID))
	}

	// 给管理员用户分配admin角色
	var adminRole models.Role
	core.DB.Where("name = ?", "admin").First(&adminRole)
	if adminRole.ID > 0 {
		var count int64
		core.DB.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).Count(&count)
		if count == 0 {
			if err := core.DB.Create(&models.UserRole{UserID: adminUser.ID, RoleID: adminRole.ID}).Error; err != nil {
				log.E(log.Sprintf("分配admin角色失败: %v", err))
			}
		}

		// 设置用户缓存（无论用户是新创建还是已存在）
		core.DB.Preload("Permissions").First(&adminRole, adminRole.ID)
		var roles []string
		var permNames []string
		roles = append(roles, adminRole.Name)
		for _, p := range adminRole.Permissions {
			permNames = append(permNames, p.Name)
		}
		utils.SetUserCache(adminUser.ID, adminUser.Username, adminUser.Nickname, roles, permNames)
		log.I(log.Sprintf("设置admin用户缓存: user_id=%d", adminUser.ID))
	}

	return nil
}
