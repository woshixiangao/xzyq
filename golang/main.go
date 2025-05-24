package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/xzyq/golang/utils"
	"net/http"
)

func main() {
    // 初始化数据库连接
    db := utils.InitDB()
    defer db.Close()

    r := gin.Default()

    // 设置session中间件
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))

    // 设置静态文件目录
    r.Static("/static", "./static")

    // 加载HTML模板
    r.LoadHTMLGlob("templates/*")

    // 公开路由
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", gin.H{})
    })

    r.GET("/register", func(c *gin.Context) {
        c.HTML(http.StatusOK, "register.html", gin.H{})
    })

    // API路由
    r.POST("/api/register", func(c *gin.Context) {
        utils.HandleRegister(c, db)
    })
    r.POST("/api/login", func(c *gin.Context) {
        utils.HandleLogin(c, db)
    })
    r.POST("/api/logout", utils.HandleLogout)

    // 需要认证的路由
    authorized := r.Group("/")
    authorized.Use(utils.AuthRequired)
    {
        authorized.GET("/home", func(c *gin.Context) {
            session := sessions.Default(c)
            username := session.Get("username")
            c.HTML(http.StatusOK, "home.html", gin.H{
                "username": username,
            })
        })

        // 添加日志页面路由
        authorized.GET("/logs", func(c *gin.Context) {
            c.HTML(http.StatusOK, "logs.html", gin.H{})
        })

        // API路由
        authorized.GET("/api/logs", func(c *gin.Context) {
            utils.GetLogs(c, db)
        })

        // 用户管理API
        authorized.GET("/api/users", func(c *gin.Context) {
            utils.GetUserList(c, db)
        })
        authorized.POST("/api/users", func(c *gin.Context) {
            utils.AddUser(c, db)
        })
        authorized.DELETE("/api/users/:username", func(c *gin.Context) {
            utils.DeleteUser(c, db)
        })
        authorized.PUT("/api/users/:username/password", func(c *gin.Context) {
            utils.UpdateUserPassword(c, db)
        })
    }

    // 启动服务器
    r.Run(":8080")
}