package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/xzyq/golang/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

// User 结构体用于表示用户数据
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

// 注册处理函数
func handleRegister(c *gin.Context) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        metadata := &utils.LogMetadata{
            Action: "register_attempt",
            Status: http.StatusBadRequest,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", fmt.Sprintf("注册失败：无效的请求数据 - %v", err), metadata)
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 检查用户名是否已存在
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
    if err != nil {
        metadata := &utils.LogMetadata{
            Action: "register_check",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", fmt.Sprintf("注册失败：数据库查询错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    if exists {
        metadata := &utils.LogMetadata{
            Action: "register_duplicate",
            Status: http.StatusConflict,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
            },
        }
        utils.InfoLogger("auth", fmt.Sprintf("注册失败：用户名 %s 已存在", user.Username), metadata)
        c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
        return
    }

    // 对密码进行加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        metadata := &utils.LogMetadata{
            Action: "register_hash",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", fmt.Sprintf("注册失败：密码加密错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    // 将用户信息插入数据库
    _, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
    if err != nil {
        metadata := &utils.LogMetadata{
            Action: "register_insert",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", fmt.Sprintf("注册失败：数据库插入错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
        return
    }

    metadata := &utils.LogMetadata{
        Action: "register_success",
        Status: http.StatusOK,
        ExtraInfo: map[string]interface{}{
            "username": user.Username,
        },
    }
    utils.InfoLogger("auth", fmt.Sprintf("用户 %s 注册成功", user.Username), metadata)
    utils.DbLogger("auth", fmt.Sprintf("新用户注册：%s", user.Username), metadata)
    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func handleLogin(c *gin.Context) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        metadata := &utils.LogMetadata{
            IP:     c.ClientIP(),
            Action: "login_attempt",
            Status: http.StatusBadRequest,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", "登录失败：无效的请求数据", metadata)
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 从数据库获取用户信息
    var storedPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedPassword)
    if err == sql.ErrNoRows {
        metadata := &utils.LogMetadata{
            Action: "login_attempt",
            Status: http.StatusUnauthorized,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
                "reason": "user_not_found",
            },
        }
        utils.InfoLogger("auth", fmt.Sprintf("登录失败：用户 %s 不存在", user.Username), metadata)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    } else if err != nil {
        metadata := &utils.LogMetadata{
            Action: "login_attempt",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        utils.ErrorLogger("auth", fmt.Sprintf("登录失败：数据库查询错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    // 验证密码
    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
    if err != nil {
        metadata := &utils.LogMetadata{
            Action: "login_attempt",
            Status: http.StatusUnauthorized,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
                "reason": "invalid_password",
            },
        }
        utils.InfoLogger("auth", fmt.Sprintf("登录失败：用户 %s 密码错误", user.Username), metadata)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 设置session
    session := sessions.Default(c)
    session.Set("username", user.Username)
    session.Save()

    metadata := &utils.LogMetadata{
        Action: "login_success",
        Status: http.StatusOK,
        ExtraInfo: map[string]interface{}{
            "username": user.Username,
            "ip": c.ClientIP(),
        },
    }
    utils.InfoLogger("auth", fmt.Sprintf("用户 %s 登录成功", user.Username), metadata)
    utils.DbLogger("auth", fmt.Sprintf("用户登录：%s", user.Username), metadata)
    c.JSON(http.StatusOK, gin.H{"message": "登录成功", "redirect": "/home"})
}

// 登出处理函数
func handleLogout(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	
	// 记录退出日志
	metadata := &utils.LogMetadata{
		Action: "logout",
		Status: http.StatusOK,
		ExtraInfo: map[string]interface{}{
			"username": username,
			"ip":      c.ClientIP(),
		},
	}
	utils.InfoLogger("auth", fmt.Sprintf("用户 %v 退出登录", username), metadata)
	
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// 认证中间件
func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
	c.Next()
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

		// 获取日志列表
		authorized.GET("/api/logs", func(c *gin.Context) {
		    // 获取查询参数
		    level := c.Query("level")
		    component := c.Query("component")
		    startDate := c.Query("startDate")
		    endDate := c.Query("endDate")
		    keyword := c.Query("keyword")
		
		    // 构建基础查询
		    query := `SELECT id, level, component, message, metadata, created_at 
		              FROM system_logs 
		              WHERE 1=1`
		    var params []interface{}
		    paramCount := 1
		
		    // 添加过滤条件
		    if level != "" {
		        query += fmt.Sprintf(" AND level = $%d", paramCount)
		        params = append(params, level)
		        paramCount++
		    }
		
		    if component != "" {
		        query += fmt.Sprintf(" AND component = $%d", paramCount)
		        params = append(params, component)
		        paramCount++
		    }
		
		    if startDate != "" {
		        query += fmt.Sprintf(" AND created_at >= $%d", paramCount)
		        params = append(params, startDate)
		        paramCount++
		    }
		
		    if endDate != "" {
		        query += fmt.Sprintf(" AND created_at <= $%d", paramCount)
		        params = append(params, endDate)
		        paramCount++
		    }
		
		    if keyword != "" {
		        query += fmt.Sprintf(" AND (message ILIKE $%d OR metadata::text ILIKE $%d)", paramCount, paramCount)
		        searchPattern := "%" + keyword + "%"
		        params = append(params, searchPattern)
		        paramCount++
		    }
		
		    // 添加排序和分页
		    query += " ORDER BY created_at DESC LIMIT 100"
		
		    // 执行查询
		    rows, err := db.Query(query, params...)
		    if err != nil {
		        c.JSON(http.StatusInternalServerError, gin.H{"error": "查询日志失败"})
		        return
		    }
		    defer rows.Close()
		
		    // 构建结果
		    var logs []gin.H
		    for rows.Next() {
		        var (
		            id int
		            level, component, message string
		            metadataBytes []byte
		            createdAt time.Time
		        )
		
		        if err := rows.Scan(&id, &level, &component, &message, &metadataBytes, &createdAt); err != nil {
		            continue
		        }

		        var metadata map[string]interface{}
		        if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		            metadata = make(map[string]interface{})
		        }
		
		        logs = append(logs, gin.H{
		            "id": id,
		            "level": level,
		            "component": component,
		            "message": message,
		            "metadata": metadata,
		            "created_at": createdAt.Format("2006-01-02 15:04:05"),
		        })
		    }

		    c.JSON(http.StatusOK, logs)
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