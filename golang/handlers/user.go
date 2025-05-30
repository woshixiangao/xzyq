package handlers

import (
	"fmt"
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

	// 如果指定了组织ID，检查组织是否存在
	if user.OrgID != nil {
		var org models.Organization
		if err := database.DB.First(&org, *user.OrgID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization not found"})
			return
		}
	}

	// 对密码进行加密
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入用户名和密码"})
		return
	}

	// 查找用户（包括软删除的用户）
	var user models.User
	result := database.DB.Unscoped().Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		fmt.Printf("用户登录失败: %s, 错误: %v\n", loginData.Username, result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户是否被删除
	if user.DeletedAt.Valid {
		fmt.Printf("已删除用户尝试登录: %s\n", loginData.Username)
		c.JSON(http.StatusForbidden, gin.H{"error": "该账号已被禁用"})
		return
	}

	// 首先尝试验证密码是否已经是加密的
	passwordValid := utils.CheckPassword(loginData.Password, user.Password)

	// 如果密码验证失败，检查是否是明文密码
	if !passwordValid && loginData.Password == user.Password {
		// 密码是明文且匹配，更新为加密密码
		hashedPassword, err := utils.HashPassword(loginData.Password)
		if err != nil {
			fmt.Printf("加密密码失败: %s, 错误: %v\n", loginData.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			return
		}

		// 更新用户密码
		if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
			fmt.Printf("更新加密密码失败: %s, 错误: %v\n", loginData.Username, err)
			// 不返回错误，继续登录流程
		}

		passwordValid = true
	}

	if !passwordValid {
		fmt.Printf("密码验证失败: %s\n", loginData.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		fmt.Printf("生成token失败: %s, 错误: %v\n", loginData.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成登录凭证失败"})
		return
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	if err := database.DB.Save(&user).Error; err != nil {
		fmt.Printf("更新最后登录时间失败: %s, 错误: %v\n", loginData.Username, err)
	}

	// 记录登录日志
	log := models.Log{
		UserID:    user.ID,
		Username:  user.Username,
		Action:    "login",
		IP:        c.ClientIP(),
		Timestamp: time.Now(),
	}
	if err := database.DB.Create(&log).Error; err != nil {
		fmt.Printf("创建登录日志失败: %s, 错误: %v\n", loginData.Username, err)
	}

	fmt.Printf("用户登录成功: %s\n", loginData.Username)

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
	if err := database.DB.Preload("Org").Find(&users).Error; err != nil {
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

	// 如果指定了组织ID，检查组织是否存在
	if updateData.OrgID != nil {
		var org models.Organization
		if err := database.DB.First(&org, *updateData.OrgID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization not found"})
			return
		}
	}

	// 更新用户信息
	if err := database.DB.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 重新查询用户信息以获取关联的组织数据
	database.DB.Preload("Org").First(&user, id)

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 开启事务
	tx := database.DB.Begin()

	// 查找用户 - 使用 Unscoped() 忽略软删除
	var user models.User
	if err := tx.Unscoped().First(&user, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 执行删除操作（硬删除）
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除用户失败: %v", err)})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("提交事务失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetProfile 获取当前用户的个人资料
func GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile 更新当前用户的个人资料
func UpdateProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 绑定更新数据
	var updateData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	// 如果要更新用户名，检查是否已存在
	if updateData.Username != "" && updateData.Username != user.Username {
		var existingUser models.User
		result := database.DB.Where("username = ?", updateData.Username).First(&existingUser)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		user.Username = updateData.Username
	}

	// 更新其他信息
	if updateData.Password != "" {
		user.Password = updateData.Password
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Phone != "" {
		user.Phone = updateData.Phone
	}

	// 保存更新
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 获取请求数据
	var passwordData struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password data"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证原密码
	if !utils.CheckPassword(passwordData.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// 对新密码进行加密
	hashedPassword, err := utils.HashPassword(passwordData.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 更新密码
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
