package api

import (
	"gin-advanced/models"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRole
// @Summary 创建角色
// @Description 创建一个新的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.Role true "角色信息"
// @Success 200 {object} utils.Response{data=models.Role}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles [post]
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

// GetRoles
// @Summary 获取角色列表
// @Description 获取所有角色的列表
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.Role}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles [get]
func GetRoles(c *gin.Context) {
	roles, err := service.GetRoles()
	if err != nil {
		utils.Error(c, "获取角色失败")
		return
	}
	utils.Success(c, roles)
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID uint `json:"userID" binding:"required"`
	RoleID uint `json:"roleID" binding:"required"`
}

// AssignRole
// @Summary 给用户分配角色
// @Description 为指定用户分配一个角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body AssignRoleRequest true "分配信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles/assign [post]
func AssignRole(c *gin.Context) {
	var req AssignRoleRequest
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

// AddPermissionRequest 添加权限请求
type AddPermissionRequest struct {
	RoleID       uint `json:"roleID" binding:"required"`
	PermissionID uint `json:"permissionID" binding:"required"`
}

// AddPermission
// @Summary 给角色添加权限
// @Description 为指定角色添加一个权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body AddPermissionRequest true "权限信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles/permission [post]
func AddPermission(c *gin.Context) {
	var req AddPermissionRequest
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

// SetRolePermissionsRequest 批量设置角色权限请求
type SetRolePermissionsRequest struct {
	RoleID        uint   `json:"role_id" binding:"required"`
	PermissionIDs []uint `json:"permission_ids"`
}

// SetRolePermissions
// @Summary 批量设置角色权限
// @Description 替换指定角色的所有权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SetRolePermissionsRequest true "权限列表"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles/permission [put]
func SetRolePermissions(c *gin.Context) {
	var req SetRolePermissionsRequest
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

// SeedRoles
// @Summary 初始化默认角色和权限
// @Description 创建默认的admin角色和权限，并分配给admin用户
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/seed [get]
func SeedRoles(c *gin.Context) {
	if err := service.SeedRoles(); err != nil {
		utils.Error(c, "初始化失败")
		return
	}
	utils.Success(c, "初始化成功")
}

// GetPermissions
// @Summary 获取所有权限
// @Description 获取系统中所有可用的权限列表
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.Permission}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/roles/permission [get]
func GetPermissions(c *gin.Context) {
	permissions, err := service.GetPermissions()
	if err != nil {
		utils.Error(c, "获取权限失败")
		return
	}
	utils.Success(c, permissions)
}

// CreatePermission
// @Summary 创建权限
// @Description 创建一个新的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.Permission true "权限信息"
// @Success 200 {object} utils.Response{data=models.Permission}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/permissions [post]
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

// UpdatePermission
// @Summary 更新权限
// @Description 更新指定权限的信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.Permission true "权限信息"
// @Success 200 {object} utils.Response{data=models.Permission}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/permissions/{id} [put]
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

// DeletePermission
// @Summary 删除权限
// @Description 删除指定权限
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "权限ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /admin/permissions/{id} [delete]
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
