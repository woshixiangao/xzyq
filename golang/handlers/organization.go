package handlers

import (
	"net/http"
	"xzyq/database"
	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// GetOrganizations 获取所有组织
func GetOrganizations(c *gin.Context) {
	var organizations []models.Organization
	db := database.GetDB()
	result := db.Find(&organizations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组织列表失败"})
		return
	}
	c.JSON(http.StatusOK, organizations)
}

// GetOrganization 获取单个组织
func GetOrganization(c *gin.Context) {
	id := c.Param("id")
	var organization models.Organization
	result := database.DB.First(&organization, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "组织不存在"})
		return
	}
	c.JSON(http.StatusOK, organization)
}

// CreateOrganization 创建组织
func CreateOrganization(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的操作"})
		return
	}

	var organization models.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 设置创建者ID
	organization.CreatedBy = userID.(uint)
	// 设置父组织ID为null，因为这是一个新的顶级组织
	organization.ParentID = nil

	// 开启数据库事务
	tx := database.DB.Begin()

	// 创建组织
	result := tx.Create(&organization)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建组织失败"})
		return
	}

	// 创建管理员用户
	adminUser := models.User{
		Username:  "admin_" + organization.Name, // 使用组织名称创建唯一的管理员用户名
		Password:  "123456",                     // 默认密码，建议后续要求用户修改
		Email:     "",
		Phone:     "",
		IsActive:  true,
		Role:      "admin",
		OrgID:     &organization.ID, // 设置组织ID
		CreatedBy: userID.(uint),    // 设置创建者ID为当前用户
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := tx.Where("username = ?", adminUser.Username).First(&existingUser).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "管理员用户名已存在"})
		return
	}

	// 创建管理员用户
	if err := tx.Create(&adminUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建管理员用户失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 返回组织信息和管理员账号信息
	c.JSON(http.StatusCreated, gin.H{
		"organization": organization,
		"admin_user": gin.H{
			"username": adminUser.Username,
			"password": "123456", // 返回默认密码给前端显示
		},
	})
}

// UpdateOrganization 更新组织
func UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	var organization models.Organization
	if err := database.DB.First(&organization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "组织不存在"})
		return
	}

	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	result := database.DB.Save(&organization)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新组织失败"})
		return
	}

	c.JSON(http.StatusOK, organization)
}

// DeleteOrganization 删除组织
func DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	var organization models.Organization
	if err := database.DB.First(&organization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "组织不存在"})
		return
	}

	result := database.DB.Delete(&organization)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除组织失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组织已删除"})
}
