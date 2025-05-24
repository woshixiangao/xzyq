package utils

import (
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// GetUserList 获取用户列表
func GetUserList(c *gin.Context, db *sql.DB) {
    rows, err := db.Query("SELECT username, created_at FROM users ORDER BY created_at DESC")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
        return
    }
    defer rows.Close()

    var users []gin.H
    for rows.Next() {
        var username string
        var createdAt time.Time
        if err := rows.Scan(&username, &createdAt); err != nil {
            continue
        }
        users = append(users, gin.H{
            "username": username,
            "created_at": createdAt,
        })
    }
    c.JSON(http.StatusOK, users)
}

// AddUser 添加新用户
func AddUser(c *gin.Context, db *sql.DB) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 检查用户名是否已存在
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    if exists {
        c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
        return
    }

    // 对密码进行加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    // 将用户信息插入数据库
    _, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "添加用户失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "添加用户成功"})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context, db *sql.DB) {
    username := c.Param("username")
    
    // 不允许删除当前登录用户
    session := sessions.Default(c)
    currentUser := session.Get("username").(string)
    if username == currentUser {
        c.JSON(http.StatusForbidden, gin.H{"error": "不能删除当前登录用户"})
        return
    }

    result, err := db.Exec("DELETE FROM users WHERE username = $1", username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "删除用户成功"})
}

// UpdateUserPassword 修改用户密码
func UpdateUserPassword(c *gin.Context, db *sql.DB) {
    username := c.Param("username")
    var data struct {
        NewPassword string `json:"new_password"`
    }

    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 对新密码进行加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
        return
    }

    result, err := db.Exec("UPDATE users SET password = $1 WHERE username = $2", string(hashedPassword), username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "修改密码失败"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "修改密码成功"})
}