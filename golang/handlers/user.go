package handlers

import (
	"net/http"
	"strconv"

	"xzyq/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`
	RoleIDs  []uint `json:"role_ids"`
	IsActive bool   `json:"is_active"`
}

// RoleRequest 角色请求结构
type RoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Permissions string `json:"permissions" binding:"required"`
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	var users []models.User

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&models.User{})

	// 支持按角色ID筛选
	if roleID := c.Query("role_id"); roleID != "" {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").Where("user_roles.role_id = ?", roleID)
	}

	// 支持按活跃状态筛选
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取用户列表
	result := query.Preload("Roles").Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      users,
	})
}

// CreateUser 创建新用户
func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if result := db.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 检查邮箱是否已存在
	if result := db.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		IsActive: true,
	}

	// 开始事务
	tx := db.Begin()

	// 创建用户
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 添加角色
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := tx.Find(&roles, req.RoleIDs).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role IDs"})
			return
		}
		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles"})
			return
		}
	}

	// 提交事务
	tx.Commit()

	// 返回创建的用户（包含角色信息）
	db.Preload("Roles").First(&user, user.ID)
	c.JSON(http.StatusCreated, user)
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := db.Preload("Roles").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否与其他用户冲突
	if user.Username != req.Username {
		var existingUser models.User
		if result := db.Where("username = ? AND id != ?", req.Username, id).First(&existingUser); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	// 检查邮箱是否与其他用户冲突
	if user.Email != req.Email {
		var existingUser models.User
		if result := db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	// 开始事务
	tx := db.Begin()

	// 更新用户基本信息
	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.IsActive = req.IsActive

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// 更新用户信息
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 更新角色
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := tx.Find(&roles, req.RoleIDs).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role IDs"})
			return
		}
		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update roles"})
			return
		}
	}

	// 提交事务
	tx.Commit()

	// 返回更新后的用户（包含角色信息）
	db.Preload("Roles").First(&user, user.ID)
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 执行删除操作
	result := db.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
