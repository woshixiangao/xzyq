package main

import (
	"log"
	"xzyq/database"
	"xzyq/handlers"
	"xzyq/middleware"
	"xzyq/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	database.InitDB()

	// 自动迁移数据库表
	db := database.GetDB()

	// 删除现有的外键约束
	if err := db.Migrator().DropConstraint(&models.User{}, "fk_users_org"); err != nil {
		log.Printf("删除外键约束失败: %v", err)
	}

	// 自动迁移数据库表
	db.AutoMigrate(&models.User{}, &models.Log{}, &models.Organization{}, &models.ObjectClass{})

	// 手动添加外键约束
	if err := db.Exec(`ALTER TABLE users 
		ADD CONSTRAINT fk_users_org 
		FOREIGN KEY (org_id) 
		REFERENCES organization(id) 
		ON DELETE SET NULL`).Error; err != nil {
		log.Printf("添加外键约束失败: %v", err)
	}

	// 创建Gin路由
	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/register", handlers.RegisterUser)
		public.POST("/login", handlers.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		protected.POST("/logout", handlers.Logout)
		protected.GET("/users", handlers.GetUsers)
		protected.GET("/users/:id", handlers.GetUser)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)

		// 个人资料相关路由
		protected.GET("/user/profile", handlers.GetProfile)
		protected.PUT("/user/profile", handlers.UpdateProfile)
		protected.PUT("/user/change-password", handlers.ChangePassword)

		// 组织管理路由
		protected.GET("/organizations", handlers.GetOrganizations)
		protected.GET("/organizations/all", handlers.GetAllOrganizations)
		protected.GET("/organizations/:id", handlers.GetOrganization)
		protected.GET("/organizations/:id/users", handlers.GetOrganizationUsers)
		protected.POST("/organizations", handlers.CreateOrganization)
		protected.PUT("/organizations/:id", handlers.UpdateOrganization)
		protected.DELETE("/organizations/:id", handlers.DeleteOrganization)

		// 对象类管理路由
		protected.GET("/object-classes", handlers.GetObjectClasses)
		protected.GET("/object-classes/:id", handlers.GetObjectClass)
		protected.POST("/object-classes", handlers.CreateObjectClass)
		protected.PUT("/object-classes/:id", handlers.UpdateObjectClass)
		protected.DELETE("/object-classes/:id", handlers.DeleteObjectClass)
		protected.GET("/object-classes/:id/children", handlers.GetObjectClassChildren)
		protected.POST("/object-classes/:id/children", handlers.CreateChildObjectClass)
	}

	// 管理员路由
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// 这里可以添加管理员特有的路由
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
