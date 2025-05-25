package handlers

import (
	"net/http"
	"strconv"

	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// OrganizationRequest 组织请求结构
type OrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

// ListOrganizations 获取组织列表
func ListOrganizations(c *gin.Context) {
	var organizations []models.Organization

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	db.Model(&models.Organization{}).Count(&total)

	// 获取组织列表
	result := db.Offset(offset).Limit(pageSize).Find(&organizations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      organizations,
	})
}

// CreateOrganization 创建新组织
func CreateOrganization(c *gin.Context) {
	var req OrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查组织名是否已存在
	var existingOrg models.Organization
	if result := db.Where("name = ?", req.Name).First(&existingOrg); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Organization name already exists"})
		return
	}

	// 如果指定了父组织，检查其是否存在
	if req.ParentID != nil {
		var parentOrg models.Organization
		if result := db.First(&parentOrg, *req.ParentID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent organization not found"})
			return
		}
	}

	// 创建新组织
	org := models.Organization{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if result := db.Create(&org); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// GetOrganization 获取单个组织信息
func GetOrganization(c *gin.Context) {
	id := c.Param("id")

	var org models.Organization
	result := db.Preload("Parent").Preload("Projects").First(&org, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrganization 更新组织信息
func UpdateOrganization(c *gin.Context) {
	id := c.Param("id")

	var org models.Organization
	if result := db.First(&org, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var req OrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查新名称是否与其他组织冲突
	if org.Name != req.Name {
		var existingOrg models.Organization
		if result := db.Where("name = ? AND id != ?", req.Name, id).First(&existingOrg); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Organization name already exists"})
			return
		}
	}

	// 更新组织信息
	org.Name = req.Name
	org.Description = req.Description
	org.ParentID = req.ParentID

	if result := db.Save(&org); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganization 删除组织
func DeleteOrganization(c *gin.Context) {
	id := c.Param("id")

	// 检查是否存在子组织
	var childCount int64
	db.Model(&models.Organization{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete organization with child organizations"})
		return
	}

	// 检查是否存在关联的项目
	var projectCount int64
	db.Model(&models.Project{}).Where("organization_id = ?", id).Count(&projectCount)
	if projectCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete organization with associated projects"})
		return
	}

	// 执行删除操作
	result := db.Delete(&models.Organization{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
