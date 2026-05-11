package api

import (
	"gin-advanced/log"
	"gin-advanced/models"
	"gin-advanced/service"
	"gin-advanced/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Token       string   `json:"token"`
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// Register
// @Summary 用户注册
// @Description 用户注册接口，需要提供用户名、密码和验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		utils.BadRequest(c, "验证码错误")
		return
	}

	// 验证用户名
	if req.Username == "" {
		utils.BadRequest(c, "用户名不能为空")
		return
	}

	// 验证密码长度
	if len(req.Password) < 6 {
		utils.BadRequest(c, "密码长度不能少于6位")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	err := service.CreateUser(&user)
	if err != nil {
		log.E(log.Sprintf("用户注册失败: %v", err))
		utils.Error(c, err.Error())
		return
	}
	log.I(log.Sprintf("用户 %s 注册成功", user.Username))
	utils.Success(c, "注册成功")
}

// Login
// @Summary 用户登录
// @Description 用户登录接口，需要用户名、密码和验证码，返回JWT token
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.W(log.Sprintf("登录参数错误: %v", err))
		utils.Error(c, "参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		log.W(log.Sprintf("验证码错误: captchaID=%s", req.CaptchaID))
		utils.BadRequest(c, "验证码错误")
		return
	}

	user, err := service.GetUserByUsername(req.Username)
	if err != nil || !utils.CheckPassword(req.Password, user.Password) {
		log.W(log.Sprintf("登录失败: 用户名或密码错误 [%s]", req.Username))
		utils.Error(c, "账号或密码错误")
		return
	}

	// 查询用户角色和权限
	role, err := service.GetUserRoles(user.ID)
	if err != nil {
		log.E(log.Sprintf("获取用户角色失败: user_id=%d, err=%v", user.ID, err))
		utils.Error(c, "获取用户信息失败")
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		log.E(log.Sprintf("生成 Token 失败: user_id=%d, err=%v", user.ID, err))
		utils.Error(c, "生成 token 失败")
		return
	}

	// 构建返回数据
	var roles []string
	var permissions []string
	for _, v := range role {
		roles = append(roles, v.Name)
		for _, perm := range v.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	// 写入Redis缓存
	utils.SetUserCache(user.ID, user.Username, user.Nickname, roles, permissions)
	utils.SetTokenCache(token, user.ID)

	log.I(log.Sprintf("用户登录成功: username=%s, user_id=%d", user.Username, user.ID))

	utils.Success(c, LoginResponse{
		Token:       token,
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Roles:       roles,
		Permissions: permissions,
	})
}

// UserInfo
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=models.User}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /user/info/{id} [get]
func UserInfo(c *gin.Context) {
	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	user, err := service.GetUserByID(intId)
	if err != nil {
		utils.Error(c, "用户不存在")
		return
	}
	utils.Success(c, user)
}

// UpdateUser
// @Summary 更新用户信息
// @Description 更新指定用户的详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param body body models.User true "用户信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /user/update/{id} [put]
func UpdateUser(c *gin.Context) {
	// 获取当前用户ID
	currentUserID := utils.GetUserIDFromContext(c)
	if currentUserID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	// 获取当前用户的权限
	permissions, _ := service.GetUserPermissions(currentUserID)
	hasPermission := false
	for _, p := range permissions {
		if p == "user:update" || p == "user:manage" {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		utils.NoPermission(c)
		return
	}

	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	err := service.UpdateUser(intId, &user)
	if err != nil {
		utils.Error(c, "更新失败")
		return
	}
	utils.Success(c, "更新成功")
}

// DeleteUser
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /user/delete/{id} [delete]
func DeleteUser(c *gin.Context) {
	// 获取当前用户ID
	currentUserID := utils.GetUserIDFromContext(c)
	if currentUserID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	// 获取当前用户的权限
	permissions, _ := service.GetUserPermissions(currentUserID)
	hasPermission := false
	for _, p := range permissions {
		if p == "user:delete" || p == "user:manage" {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		utils.NoPermission(c)
		return
	}

	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	err := service.DeleteUser(intId)
	if err != nil {
		utils.Error(c, "删除失败")
		return
	}
	utils.Success(c, "删除成功")
}
