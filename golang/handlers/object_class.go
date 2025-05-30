package handlers

import (
	"net/http"
	"time"
	"xzyq/database"
	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// GetObjectClasses 获取对象类列表
func GetObjectClasses(c *gin.Context) {
	var classes []models.ObjectClass
	if err := database.DB.Preload("Organization").
		Preload("CreatedByUser").
		Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取对象类列表失败"})
		return
	}
	c.JSON(http.StatusOK, classes)
}

// GetObjectClass 获取单个对象类
func GetObjectClass(c *gin.Context) {
	id := c.Param("id")
	var class models.ObjectClass

	if err := database.DB.Preload("Organization").
		Preload("CreatedByUser").
		First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "对象类不存在"})
		return
	}

	c.JSON(http.StatusOK, class)
}

// CreateObjectClass 创建对象类
func CreateObjectClass(c *gin.Context) {
	var class models.ObjectClass
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取当前用户ID和组织ID
	userID, _ := c.Get("userID")
	class.CreatedBy = userID.(uint)

	// 获取当前用户信息以获取其组织ID
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 使用当前用户的组织ID
	class.OrgID = *currentUser.OrgID
	class.UpdatedAt = time.Now()

	if err := database.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建对象类失败"})
		return
	}

	// 重新获取完整的对象类信息
	database.DB.Preload("Organization").
		Preload("CreatedByUser").
		First(&class, class.ID)

	c.JSON(http.StatusCreated, class)
}

// UpdateObjectClass 更新对象类
func UpdateObjectClass(c *gin.Context) {
	id := c.Param("id")
	var class models.ObjectClass

	// 检查对象类是否存在
	if err := database.DB.First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "对象类不存在"})
		return
	}

	// 绑定更新数据
	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的更新数据"})
		return
	}

	// 更新对象类（只更新允许的字段）
	updates := map[string]interface{}{
		"name":        updateData.Name,
		"description": updateData.Description,
		"updated_at":  time.Now(),
	}

	if err := database.DB.Model(&class).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新对象类失败"})
		return
	}

	// 重新获取完整的对象类信息
	database.DB.Preload("Organization").
		Preload("CreatedByUser").
		First(&class, id)

	c.JSON(http.StatusOK, class)
}

// DeleteObjectClass 删除对象类
func DeleteObjectClass(c *gin.Context) {
	id := c.Param("id")

	// 执行删除
	if err := database.DB.Delete(&models.ObjectClass{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除对象类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "对象类删除成功"})
}
