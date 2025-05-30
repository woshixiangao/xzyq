package handlers

import (
	"fmt"
	"net/http"
	"xzyq/database"
	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// GetOrganizations 获取所有组织
func GetOrganizations(c *gin.Context) {
	type OrgWithDetails struct {
		models.Organization
		UserCount int64                `json:"user_count"`
		ParentOrg *models.Organization `json:"parent_org,omitempty"`
	}

	// 从上下文中获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的操作"})
		return
	}

	var orgsWithDetails []OrgWithDetails
	db := database.GetDB()

	// 查询当前用户创建的组织
	var organizations []models.Organization
	if err := db.Where("created_by = ?", userID).Find(&organizations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组织列表失败"})
		return
	}

	// 为每个组织获取详细信息
	for _, org := range organizations {
		var count int64
		if err := db.Unscoped().Model(&models.User{}).Where("org_id = ?", org.ID).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "统计用户数量失败"})
			return
		}

		details := OrgWithDetails{
			Organization: org,
			UserCount:    count,
		}

		// 如果有父组织ID，查询父组织信息
		if org.ParentID != nil {
			var parentOrg models.Organization
			if err := db.First(&parentOrg, *org.ParentID).Error; err == nil {
				details.ParentOrg = &parentOrg
			}
		}

		orgsWithDetails = append(orgsWithDetails, details)
	}

	c.JSON(http.StatusOK, orgsWithDetails)
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

	// 开启事务
	tx := database.DB.Begin()

	// 查找组织
	var organization models.Organization
	if err := tx.First(&organization, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("组织(ID:%s)不存在", id)})
		return
	}

	// 查找该组织下的用户数量
	var userCount int64
	if err := tx.Model(&models.User{}).Where("org_id = ?", organization.ID).Count(&userCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("查询组织(ID:%d)的用户数量失败: %v", organization.ID, err)})
		return
	}

	// 如果组织下还有用户，则不允许删除
	if userCount > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("无法删除组织[%s]，该组织下还有 %d 个用户", organization.Name, userCount),
		})
		return
	}

	// 删除组织相关的所有数据
	// 1. 删除组织下的用户（虽然已经确认数量为0，但为了保险起见）
	if err := tx.Where("org_id = ?", organization.ID).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除组织[%s]的用户数据失败: %v", organization.Name, err)})
		return
	}

	// 2. 删除组织本身
	if err := tx.Delete(&organization).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除组织[%s]失败: %v", organization.Name, err)})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除组织[%s]时提交事务失败: %v", organization.Name, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("组织[%s]及其相关数据已成功删除", organization.Name),
	})
}

// GetAllOrganizations 获取所有组织（用于父级租户选择）
func GetAllOrganizations(c *gin.Context) {
	var organizations []models.Organization
	if err := database.DB.Find(&organizations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组织列表失败"})
		return
	}
	c.JSON(http.StatusOK, organizations)
}

// GetOrganizationUsers 获取组织下的用户列表
func GetOrganizationUsers(c *gin.Context) {
	id := c.Param("id")

	// 验证组织是否存在
	var organization models.Organization
	if err := database.DB.First(&organization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "组织不存在"})
		return
	}

	// 获取组织下的用户列表，包括软删除的用户
	var users []models.User
	if err := database.DB.Unscoped().Where("org_id = ?", id).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	// 如果没有找到用户，返回空数组而不是 null
	if users == nil {
		users = make([]models.User, 0)
	}

	c.JSON(http.StatusOK, users)
}
