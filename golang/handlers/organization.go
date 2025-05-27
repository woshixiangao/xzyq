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
	orgID, exists := c.Get("OrgId")

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
	// 将 orgID 转换为指针类型
	parentID := orgID.(uint)
	organization.ParentID = &parentID

	result := database.DB.Create(&organization)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建组织失败"})
		return
	}

	c.JSON(http.StatusCreated, organization)
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
