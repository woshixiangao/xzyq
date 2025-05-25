package models

import (
	"time"

	"gorm.io/gorm"
)

// Log 日志模型
type Log struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UserID    uint      `json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	Action    string    `gorm:"size:50" json:"action"` // login或logout
	IP        string    `gorm:"size:50" json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

// TableName 指定表名
func (Log) TableName() string {
	return "logs"
}
