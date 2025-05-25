package handlers

import "gorm.io/gorm"

// DB 全局数据库实例
var db *gorm.DB

// SetDB 设置全局数据库实例
func SetDB(database *gorm.DB) {
	db = database
}
