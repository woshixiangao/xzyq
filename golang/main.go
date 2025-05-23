package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// 设置静态文件目录
	r.Static("/static", "./static")

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	// 启动服务器
	r.Run(":8080")
}