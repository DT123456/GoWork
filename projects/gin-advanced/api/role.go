package api

import (
	"gin-advanced/models"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.CreateRole(&role); err != nil {
		utils.Error(c, "创建角色失败")
		return
	}
	utils.Success(c, role)
}

// GetRoles 获取所有角色
func GetRoles(c *gin.Context) {
	roles, err := service.GetRoles()
	if err != nil {
		utils.Error(c, "获取角色失败")
		return
	}
	utils.Success(c, roles)
}

// AssignRole 给用户分配角色
func AssignRole(c *gin.Context) {
	var req struct {
		UserID uint `json:"userID" binding:"required"`
		RoleID uint `json:"roleID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.AssignRoleToUser(req.UserID, req.RoleID); err != nil {
		utils.Error(c, "分配角色失败")
		return
	}
	utils.Success(c, "分配成功")
}

// AddPermission 给角色添加权限
func AddPermission(c *gin.Context) {
	var req struct {
		RoleID       uint `json:"roleID" binding:"required"`
		PermissionID uint `json:"permissionID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.AddPermissionToRole(req.RoleID, req.PermissionID); err != nil {
		utils.Error(c, "添加权限失败")
		return
	}
	utils.Success(c, "添加成功")
}

// SetRolePermissions 批量设置角色权限
func SetRolePermissions(c *gin.Context) {
	var req struct {
		RoleID        uint   `json:"role_id" binding:"required"`
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.SetRolePermissions(req.RoleID, req.PermissionIDs); err != nil {
		utils.Error(c, "设置权限失败")
		return
	}
	utils.Success(c, "设置成功")
}

// SeedRoles 初始化默认角色和权限
func SeedRoles(c *gin.Context) {
	if err := service.SeedRoles(); err != nil {
		utils.Error(c, "初始化失败")
		return
	}
	utils.Success(c, "初始化成功")
}

// GetPermissions 获取所有权限
func GetPermissions(c *gin.Context) {
	permissions, err := service.GetPermissions()
	if err != nil {
		utils.Error(c, "获取权限失败")
		return
	}
	utils.Success(c, permissions)
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var perm models.Permission
	if err := c.ShouldBindJSON(&perm); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.CreatePermission(&perm); err != nil {
		utils.Error(c, "创建权限失败")
		return
	}
	utils.Success(c, perm)
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	var perm models.Permission
	if err := c.ShouldBindJSON(&perm); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.UpdatePermission(&perm); err != nil {
		utils.Error(c, "更新权限失败")
		return
	}
	utils.Success(c, perm)
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, "参数错误")
		return
	}
	if err := service.DeletePermission(uint(id)); err != nil {
		utils.Error(c, "删除权限失败")
		return
	}
	utils.Success(c, "删除成功")
}
