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
	db.AutoMigrate(&models.User{}, &models.Log{})

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
	}

	// 管理员路由
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// 添加管理员特有的路由
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
