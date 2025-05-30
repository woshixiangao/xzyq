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
		Preload("ParentClass").
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
		Preload("ParentClass").
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

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	class.CreatedBy = userID.(uint)
	class.UpdatedAt = time.Now()

	// 如果指定了父类，检查父类是否存在
	if class.ParentClassID != nil {
		var parentClass models.ObjectClass
		if err := database.DB.First(&parentClass, *class.ParentClassID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父类不存在"})
			return
		}
	}

	// 检查组织是否存在
	var org models.Organization
	if err := database.DB.First(&org, class.OrgID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "组织不存在"})
		return
	}

	if err := database.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建对象类失败"})
		return
	}

	// 重新获取完整的对象类信息
	database.DB.Preload("Organization").
		Preload("ParentClass").
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
	var updateData models.ObjectClass
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的更新数据"})
		return
	}

	// 更新时间
	updateData.UpdatedAt = time.Now()

	// 如果更新了父类ID，检查父类是否存在
	if updateData.ParentClassID != nil {
		var parentClass models.ObjectClass
		if err := database.DB.First(&parentClass, *updateData.ParentClassID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父类不存在"})
			return
		}
	}

	// 如果更新了组织ID，检查组织是否存在
	if updateData.OrgID != 0 {
		var org models.Organization
		if err := database.DB.First(&org, updateData.OrgID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "组织不存在"})
			return
		}
	}

	// 更新对象类
	if err := database.DB.Model(&class).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新对象类失败"})
		return
	}

	// 重新获取完整的对象类信息
	database.DB.Preload("Organization").
		Preload("ParentClass").
		Preload("CreatedByUser").
		First(&class, id)

	c.JSON(http.StatusOK, class)
}

// DeleteObjectClass 删除对象类
func DeleteObjectClass(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有子类引用此对象类
	var childCount int64
	if err := database.DB.Model(&models.ObjectClass{}).Where("parent_class_id = ?", id).Count(&childCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查子类失败"})
		return
	}

	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法删除：该对象类存在子类"})
		return
	}

	// 执行删除
	if err := database.DB.Delete(&models.ObjectClass{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除对象类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "对象类删除成功"})
}
