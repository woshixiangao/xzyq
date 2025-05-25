package handlers

import (
	"net/http"
	"strconv"
	"time"

	"xzyq/models"

	"github.com/gin-gonic/gin"
)

// LogRequest 日志请求结构
type LogRequest struct {
	Action     string `json:"action" binding:"required"`
	Resource   string `json:"resource" binding:"required"`
	ResourceID uint   `json:"resource_id" binding:"required"`
	Details    string `json:"details"`
}

// CreateLog 创建系统日志
func CreateLog(userID uint, action, resource string, resourceID uint, details string, c *gin.Context) error {
	log := models.Log{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IP:         c.ClientIP(),
	}

	return db.Create(&log).Error
}

// ListLogs 获取日志列表
func ListLogs(c *gin.Context) {
	var logs []models.Log

	// 支持分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&models.Log{})

	// 支持按用户ID筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 支持按操作类型筛选
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// 支持按资源类型筛选
	if resource := c.Query("resource"); resource != "" {
		query = query.Where("resource = ?", resource)
	}

	// 支持按时间范围筛选
	if startDate := c.Query("start_date"); startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// 将结束日期设置为当天的最后一刻
			end = end.Add(24*time.Hour - time.Second)
			query = query.Where("created_at <= ?", end)
		}
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取日志列表
	result := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      logs,
	})
}

// GetLog 获取单个日志详情
func GetLog(c *gin.Context) {
	id := c.Param("id")

	var log models.Log
	result := db.Preload("User").First(&log, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// LogMiddleware 日志记录中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行请求
		c.Next()

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			return
		}

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path

		// 确定操作类型
		action := method
		resource := "unknown"
		var resourceID uint

		// 解析资源类型和ID
		switch {
		case path == "/api/auth/login":
			resource = "auth"
			action = "login"
		case path == "/api/auth/register":
			resource = "auth"
			action = "register"
		default:
			// 解析其他API路径
			// TODO: 根据实际API路径结构完善此部分
		}

		// 记录日志
		CreateLog(
			userID.(uint),
			action,
			resource,
			resourceID,
			"", // 可以添加更多详细信息
			c,
		)
	}
}
