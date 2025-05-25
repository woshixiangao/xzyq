package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 包含共同的字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string    `gorm:"uniqueIndex;not null" json:"username"`
	Password    string    `gorm:"not null" json:"-"`
	Email       string    `gorm:"uniqueIndex" json:"email"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	Roles       []Role    `gorm:"many2many:user_roles;" json:"roles"`
	LastLoginAt time.Time `json:"last_login_at"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	Permissions string `gorm:"type:json" json:"permissions"`
}

// Organization 组织模型
type Organization struct {
	BaseModel
	Name        string        `gorm:"uniqueIndex;not null" json:"name"`
	Description string        `json:"description"`
	ParentID    *uint         `json:"parent_id"`
	Parent      *Organization `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Projects    []Project     `json:"projects,omitempty"`
}

// Project 项目模型
type Project struct {
	BaseModel
	Name           string       `gorm:"not null" json:"name"`
	Description    string       `json:"description"`
	StartDate      time.Time    `json:"start_date"`
	EndDate        *time.Time   `json:"end_date"`
	Status         string       `json:"status"`
	OrganizationID uint         `json:"organization_id"`
	Organization   Organization `json:"organization,omitempty"`
	Products       []Product    `json:"products,omitempty"`
}

// Product 产品模型
type Product struct {
	BaseModel
	Name        string  `gorm:"not null" json:"name"`
	Code        string  `gorm:"uniqueIndex;not null" json:"code"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"`
	ProjectID   uint    `json:"project_id"`
	Project     Project `json:"project,omitempty"`
}

// Log 系统日志模型
type Log struct {
	BaseModel
	UserID     uint   `json:"user_id"`
	User       User   `json:"user,omitempty"`
	Action     string `gorm:"not null" json:"action"`
	Resource   string `json:"resource"`
	ResourceID uint   `json:"resource_id"`
	Details    string `gorm:"type:json" json:"details"`
	IP         string `json:"ip"`
}
