package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/xzyq/golang/utils"
	"log"
	"net/http"
)

var db *sql.DB

func initDB() {
	// 使用 PostgreSQL URL 格式
	connStr := "postgresql://myuser:mysecretpassword@47.121.141.235:5432/postgres?sslmode=disable"
	
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	fmt.Println("数据库连接成功！")
	
	// 初始化日志器
	utils.InitLogger(db)

    // 创建system_logs表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS system_logs (
            id SERIAL PRIMARY KEY,
            level VARCHAR(10) NOT NULL,
            component VARCHAR(50) NOT NULL,
            message TEXT NOT NULL,
            metadata JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal("创建system_logs表失败:", err)
    }

    // 创建索引以提高查询性能
    _, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs(level);
        CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at);
        CREATE INDEX IF NOT EXISTS idx_system_logs_component ON system_logs(component)
    `)
    if err != nil {
        log.Fatal("创建system_logs索引失败:", err)
    }
}

func main() {
    // 初始化数据库连接
    initDB()
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