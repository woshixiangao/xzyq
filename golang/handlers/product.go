package handlers

import (
	"net/http"
	"strconv"

	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// ProductRequest 产品请求结构
type ProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required"`
	Unit        string `json:"unit" binding:"required"`
	ProjectID   uint   `json:"project_id" binding:"required"`
}

// ListProducts 获取产品列表
func ListProducts(c *gin.Context) {
	var products []models.Product

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&models.Product{})

	// 支持按项目ID筛选
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// 支持按类别筛选
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// 支持按产品代码或名称搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取产品列表
	result := query.Preload("Project").Offset(offset).Limit(pageSize).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      products,
	})
}

// CreateProduct 创建新产品
func CreateProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证项目是否存在
	var project models.Project
	if result := db.First(&project, req.ProjectID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
		return
	}

	// 检查产品代码是否已存在
	var existingProduct models.Product
	if result := db.Where("code = ?", req.Code).First(&existingProduct); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product code already exists"})
		return
	}

	// 创建新产品
	product := models.Product{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		ProjectID:   req.ProjectID,
	}

	if result := db.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// 返回创建的产品（包含项目信息）
	db.Preload("Project").First(&product, product.ID)
	c.JSON(http.StatusCreated, product)
}

// GetProduct 获取单个产品信息
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	result := db.Preload("Project").First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct 更新产品信息
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if result := db.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证项目是否存在
	if product.ProjectID != req.ProjectID {
		var project models.Project
		if result := db.First(&project, req.ProjectID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
			return
		}
	}

	// 检查产品代码是否与其他产品冲突
	if product.Code != req.Code {
		var existingProduct models.Product
		if result := db.Where("code = ? AND id != ?", req.Code, id).First(&existingProduct); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Product code already exists"})
			return
		}
	}

	// 更新产品信息
	product.Name = req.Name
	product.Code = req.Code
	product.Description = req.Description
	product.Category = req.Category
	product.Unit = req.Unit
	product.ProjectID = req.ProjectID

	if result := db.Save(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 返回更新后的产品（包含项目信息）
	db.Preload("Project").First(&product, product.ID)
	c.JSON(http.StatusOK, product)
}

// DeleteProduct 删除产品
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// 执行删除操作
	result := db.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
