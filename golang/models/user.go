package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Username    string        `gorm:"size:50;not null;unique" json:"username"`
	Password    string        `gorm:"size:255;not null" json:"-"` // 密码不返回给前端
	Email       string        `gorm:"size:100" json:"email"`
	Phone       string        `gorm:"size:20" json:"phone"`
	LastLoginAt time.Time     `json:"last_login_at"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
	Role        string        `gorm:"size:20;default:'user'" json:"role"` // admin或user
	OrgID       *uint         `json:"org_id"`                             // 组织ID
	Org         *Organization `gorm:"foreignKey:OrgID" json:"org"`        // 组织关联
	CreatedBy   uint          `json:"created_by"`                         // 创建者ID
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
