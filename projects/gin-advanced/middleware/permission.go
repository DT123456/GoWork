package middleware

import (
	"gin-advanced/core"
	"gin-advanced/models"
	"gin-advanced/utils"

	"github.com/gin-gonic/gin"
)

// PermissionAuth 权限认证中间件
func PermissionAuth(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 context 获取用户 ID
		userID, exists := c.Get("userID")
		if !exists {
			utils.Error(c, "用户未登录")
			c.Abort()
			return
		}

		// 查询用户的角色和权限
		var permissions []string
		var user models.User
		if err := core.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			utils.Error(c, "用户不存在")
			c.Abort()
			return
		}

		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				permissions = append(permissions, perm.Name)
			}
		}

		// 检查是否有所需权限
		hasPermission := false
		for _, p := range permissions {
			if p == requiredPerm || p == "*" { // * 表示超级管理员
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.Error(c, "没有权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles 要求用户拥有指定角色
func RequireRoles(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Error(c, "用户未登录")
			c.Abort()
			return
		}

		var user models.User
		if err := core.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			utils.Error(c, "用户不存在")
			c.Abort()
			return
		}

		userRoleMap := make(map[string]bool)
		for _, role := range user.Roles {
			userRoleMap[role.Name] = true
		}

		for _, required := range roleNames {
			if userRoleMap[required] {
				c.Next()
				return
			}
		}

		utils.Error(c, "没有该角色权限")
		c.Abort()
	}
}
