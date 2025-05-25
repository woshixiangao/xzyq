package handlers

import (
	"net/http"
	"strconv"
	"time"

	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// ProjectRequest 项目请求结构
type ProjectRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description"`
	StartDate      time.Time  `json:"start_date" binding:"required"`
	EndDate        *time.Time `json:"end_date"`
	Status         string     `json:"status" binding:"required"`
	OrganizationID uint       `json:"organization_id" binding:"required"`
}

// ListProjects 获取项目列表
func ListProjects(c *gin.Context) {
	var projects []models.Project

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&models.Project{})

	// 支持按组织ID筛选
	if orgID := c.Query("organization_id"); orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}

	// 支持按状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取项目列表
	result := query.Preload("Organization").Offset(offset).Limit(pageSize).Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      projects,
	})
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证组织是否存在
	var org models.Organization
	if result := db.First(&org, req.OrganizationID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization not found"})
		return
	}

	// 创建新项目
	project := models.Project{
		Name:           req.Name,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         req.Status,
		OrganizationID: req.OrganizationID,
	}

	if result := db.Create(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// 返回创建的项目（包含组织信息）
	db.Preload("Organization").First(&project, project.ID)
	c.JSON(http.StatusCreated, project)
}

// GetProject 获取单个项目信息
func GetProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	result := db.Preload("Organization").Preload("Products").First(&project, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if result := db.First(&project, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证新组织是否存在
	if project.OrganizationID != req.OrganizationID {
		var org models.Organization
		if result := db.First(&org, req.OrganizationID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization not found"})
			return
		}
	}

	// 更新项目信息
	project.Name = req.Name
	project.Description = req.Description
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate
	project.Status = req.Status
	project.OrganizationID = req.OrganizationID

	if result := db.Save(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// 返回更新后的项目（包含组织信息）
	db.Preload("Organization").First(&project, project.ID)
	c.JSON(http.StatusOK, project)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	// 检查是否存在关联的产品
	var productCount int64
	db.Model(&models.Product{}).Where("project_id = ?", id).Count(&productCount)
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete project with associated products"})
		return
	}

	// 执行删除操作
	result := db.Delete(&models.Project{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
