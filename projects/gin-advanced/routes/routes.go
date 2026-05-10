package routes

import (
	"gin-advanced/api"
	_ "gin-advanced/docs"
	"gin-advanced/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.RequestLog()) // 统一请求日志

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 公共路由
	r.POST("/register", api.Register)
	r.POST("/login", api.Login)

	// 验证码
	r.GET("/captcha", api.GetCaptcha)
	r.POST("/captcha/verify", api.VerifyCaptcha)

	// 文件上传
	r.POST("/upload", api.UploadFile)
	r.POST("/upload/image", api.UploadImage)
	r.POST("/upload/avatar", api.UploadAvatar)
	r.POST("/upload/multiple", api.UploadFiles)
	r.GET("/upload/config", api.GetUploadConfig)

	// 需登录路由
	auth := r.Group("/user")
	auth.Use(middleware.JWTAuth())
	{
		auth.PUT("/update/:id", api.UpdateUser)
		auth.GET("/info/:id", api.UserInfo)
		auth.DELETE("/delete/:id", api.DeleteUser)
	}

	// 角色管理 (需要 admin 权限)
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RequireRoles("admin"))
	{
		admin.POST("/roles", api.CreateRole)
		admin.GET("/roles", api.GetRoles)
		admin.POST("/roles/assign", api.AssignRole)
		admin.GET("/roles/permission", api.GetPermissions)
		admin.POST("/roles/permission", api.AddPermission)
		admin.PUT("/roles/permission", api.SetRolePermissions)
		admin.POST("/permissions", api.CreatePermission)
		admin.PUT("/permissions/:id", api.UpdatePermission)
		admin.DELETE("/permissions/:id", api.DeletePermission)
		admin.GET("/seed", api.SeedRoles)

		// 系统日志
		admin.GET("/logs/system", api.GetSystemLogs)
		admin.GET("/logs/system/:id", api.GetSystemLogDetail)
		admin.GET("/logs/system/stats", api.GetSystemLogStatistics)
		admin.DELETE("/logs/system", api.DeleteSystemLogs)

		// 访问日志
		admin.GET("/logs/access", api.GetAccessLogs)
		admin.GET("/logs/access/stats", api.GetAccessLogStatistics)
		admin.DELETE("/logs/access", api.CleanAccessLogs)
	}

	// 权限测试
	r.GET("/test/permission", middleware.JWTAuth(), middleware.PermissionAuth("user:read"))

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
