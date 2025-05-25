package handlers

import (
	"net/http"
	"strconv"

	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// ListRoles 获取角色列表
func ListRoles(c *gin.Context) {
	var roles []models.Role

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 获取总数
	var total int64
	db.Model(&models.Role{}).Count(&total)

	// 获取角色列表
	result := db.Offset(offset).Limit(pageSize).Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      roles,
	})
}

// CreateRole 创建新角色
func CreateRole(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查角色名是否已存在
	var existingRole models.Role
	if result := db.Where("name = ?", req.Name).First(&existingRole); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists"})
		return
	}

	// 创建新角色
	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	if result := db.Create(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetRole 获取单个角色信息
func GetRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	result := db.First(&role, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole 更新角色信息
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if result := db.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查角色名是否与其他角色冲突
	if role.Name != req.Name {
		var existingRole models.Role
		if result := db.Where("name = ? AND id != ?", req.Name, id).First(&existingRole); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists"})
			return
		}
	}

	// 更新角色信息
	role.Name = req.Name
	role.Description = req.Description
	role.Permissions = req.Permissions

	if result := db.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有用户关联此角色
	var userCount int64
	db.Table("user_roles").Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete role with associated users"})
		return
	}

	// 执行删除操作
	result := db.Delete(&models.Role{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
