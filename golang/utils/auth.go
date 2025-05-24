package utils

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// User 结构体用于表示用户数据
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleRegister 处理用户注册
func HandleRegister(c *gin.Context, db *sql.DB) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        metadata := &LogMetadata{
            Action: "register_attempt",
            Status: http.StatusBadRequest,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", fmt.Sprintf("注册失败：无效的请求数据 - %v", err), metadata)
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 检查用户名是否已存在
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
    if err != nil {
        metadata := &LogMetadata{
            Action: "register_check",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", fmt.Sprintf("注册失败：数据库查询错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    if exists {
        metadata := &LogMetadata{
            Action: "register_duplicate",
            Status: http.StatusConflict,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
            },
        }
        InfoLogger("auth", fmt.Sprintf("注册失败：用户名 %s 已存在", user.Username), metadata)
        c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
        return
    }

    // 对密码进行加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        metadata := &LogMetadata{
            Action: "register_hash",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", fmt.Sprintf("注册失败：密码加密错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    // 将用户信息插入数据库
    _, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
    if err != nil {
        metadata := &LogMetadata{
            Action: "register_insert",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", fmt.Sprintf("注册失败：数据库插入错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
        return
    }

    metadata := &LogMetadata{
        Action: "register_success",
        Status: http.StatusOK,
        ExtraInfo: map[string]interface{}{
            "username": user.Username,
        },
    }
    InfoLogger("auth", fmt.Sprintf("用户 %s 注册成功", user.Username), metadata)
    DbLogger("auth", fmt.Sprintf("新用户注册：%s", user.Username), metadata)
    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// HandleLogin 处理用户登录
func HandleLogin(c *gin.Context, db *sql.DB) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        metadata := &LogMetadata{
            IP:     c.ClientIP(),
            Action: "login_attempt",
            Status: http.StatusBadRequest,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", "登录失败：无效的请求数据", metadata)
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 从数据库获取用户信息
    var storedPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedPassword)
    if err == sql.ErrNoRows {
        metadata := &LogMetadata{
            Action: "login_attempt",
            Status: http.StatusUnauthorized,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
                "reason": "user_not_found",
            },
        }
        InfoLogger("auth", fmt.Sprintf("登录失败：用户 %s 不存在", user.Username), metadata)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    } else if err != nil {
        metadata := &LogMetadata{
            Action: "login_attempt",
            Status: http.StatusInternalServerError,
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        }
        ErrorLogger("auth", fmt.Sprintf("登录失败：数据库查询错误 - %v", err), metadata)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    // 验证密码
    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
    if err != nil {
        metadata := &LogMetadata{
            Action: "login_attempt",
            Status: http.StatusUnauthorized,
            ExtraInfo: map[string]interface{}{
                "username": user.Username,
                "reason": "invalid_password",
            },
        }
        InfoLogger("auth", fmt.Sprintf("登录失败：用户 %s 密码错误", user.Username), metadata)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    // 设置session
    session := sessions.Default(c)
    session.Set("username", user.Username)
    session.Save()

    metadata := &LogMetadata{
        Action: "login_success",
        Status: http.StatusOK,
        ExtraInfo: map[string]interface{}{
            "username": user.Username,
            "ip": c.ClientIP(),
        },
    }
    InfoLogger("auth", fmt.Sprintf("用户 %s 登录成功", user.Username), metadata)
    DbLogger("auth", fmt.Sprintf("用户登录：%s", user.Username), metadata)
    c.JSON(http.StatusOK, gin.H{"message": "登录成功", "redirect": "/home"})
}

// HandleLogout 处理用户登出
func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	
	// 记录退出日志
	metadata := &LogMetadata{
		Action: "logout",
		Status: http.StatusOK,
		ExtraInfo: map[string]interface{}{
			"username": username,
			"ip":      c.ClientIP(),
		},
	}
	InfoLogger("auth", fmt.Sprintf("用户 %v 退出登录", username), metadata)
	
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// AuthRequired 认证中间件
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
	c.Next()
}