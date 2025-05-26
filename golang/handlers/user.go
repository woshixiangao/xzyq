package handlers

import (
	"net/http"
	"time"
	"xzyq/database"
	"xzyq/models"
	"xzyq/utils"

	"github.com/gin-gonic/gin"
)

// RegisterUser 注册新用户
func RegisterUser(c *gin.Context) {
	var user models.User

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 创建用户
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// 查找用户
	var user models.User
	result := database.DB.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 验证密码 - 直接比较原文
	if user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	database.DB.Save(&user)

	// 记录登录日志
	log := models.Log{
		UserID:    user.ID,
		Username:  user.Username,
		Action:    "login",
		IP:        c.ClientIP(),
		Timestamp: time.Now(),
	}
	database.DB.Create(&log)

	// 返回token和用户信息
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户退出
func Logout(c *gin.Context) {
	// 从上下文中获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	username, _ := c.Get("username")

	// 记录退出日志
	log := models.Log{
		UserID:    userID.(uint),
		Username:  username.(string),
		Action:    "logout",
		IP:        c.ClientIP(),
		Timestamp: time.Now(),
	}
	database.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 绑定更新数据
	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	// 更新用户信息
	if err := database.DB.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
