// 对象类管理路由
objectClassRoutes := router.Group("/api/object-classes")
objectClassRoutes.Use(middleware.AuthMiddleware())
{
	objectClassRoutes.GET("", handlers.GetObjectClasses)
	objectClassRoutes.GET("/:id", handlers.GetObjectClass)
	objectClassRoutes.POST("", handlers.CreateObjectClass)
	objectClassRoutes.PUT("/:id", handlers.UpdateObjectClass)
	objectClassRoutes.DELETE("/:id", handlers.DeleteObjectClass)
} 