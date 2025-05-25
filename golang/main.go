package main

import (
	"log"
	"os"

	"xzyq/database"
	"xzyq/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化数据库连接
	db = database.InitDB()

	// 设置handlers包的数据库实例
	handlers.SetDB(db)
}

func main() {
	// 设置gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由引擎
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

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.HandleRegister)
			auth.POST("/login", handlers.HandleLogin)
		}

		// 需要认证的路由组
		protected := api.Group("/")
		protected.Use(handlers.AuthMiddleware())
		{
			// 组织管理
			org := protected.Group("/organizations")
			{
				org.GET("", handlers.ListOrganizations)
				org.POST("", handlers.CreateOrganization)
				org.GET("/:id", handlers.GetOrganization)
				org.PUT("/:id", handlers.UpdateOrganization)
				org.DELETE("/:id", handlers.DeleteOrganization)
			}

			// 项目管理
			proj := protected.Group("/projects")
			{
				proj.GET("", handlers.ListProjects)
				proj.POST("", handlers.CreateProject)
				proj.GET("/:id", handlers.GetProject)
				proj.PUT("/:id", handlers.UpdateProject)
				proj.DELETE("/:id", handlers.DeleteProject)
			}

			// 产品管理
			prod := protected.Group("/products")
			{
				prod.GET("", handlers.ListProducts)
				prod.POST("", handlers.CreateProduct)
				prod.GET("/:id", handlers.GetProduct)
				prod.PUT("/:id", handlers.UpdateProduct)
				prod.DELETE("/:id", handlers.DeleteProduct)
			}

			// 用户管理
			users := protected.Group("/users")
			{
				users.GET("", handlers.ListUsers)
				users.POST("", handlers.CreateUser)
				users.GET("/:id", handlers.GetUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// 角色管理
			roles := protected.Group("/roles")
			{
				roles.GET("", handlers.ListRoles)
				roles.POST("", handlers.CreateRole)
				roles.GET("/:id", handlers.GetRole)
				roles.PUT("/:id", handlers.UpdateRole)
				roles.DELETE("/:id", handlers.DeleteRole)
			}

			// 日志管理
			logs := protected.Group("/logs")
			{
				logs.GET("", handlers.ListLogs)
				logs.GET("/:id", handlers.GetLog)
			}
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
