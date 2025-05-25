package database

import (
	"fmt"
	"log"
	"os"

	"xzyq/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	// 数据库连接配置
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("Connecting to database at %s:%s...", host, port)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")

	// 自动迁移数据库表结构
	migrateDB(db)

	return db
}

// migrateDB 执行数据库迁移
func migrateDB(db *gorm.DB) {
	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Organization{},
		&models.Project{},
		&models.Product{},
		&models.Log{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// 创建默认管理员角色（如果不存在）
	createDefaultAdminRole(db)

	// 创建默认管理员用户（如果不存在）
	createDefaultAdminUser(db)
}

// createDefaultAdminRole 创建默认管理员角色
func createDefaultAdminRole(db *gorm.DB) {
	var adminRole models.Role
	result := db.Where("name = ?", "admin").First(&adminRole)

	if result.Error == gorm.ErrRecordNotFound {
		adminRole = models.Role{
			Name:        "admin",
			Description: "System Administrator",
			Permissions: `{
				"organizations": ["create", "read", "update", "delete"],
				"projects": ["create", "read", "update", "delete"],
				"products": ["create", "read", "update", "delete"],
				"users": ["create", "read", "update", "delete"],
				"roles": ["create", "read", "update", "delete"],
				"logs": ["read"]
			}`,
		}

		if err := db.Create(&adminRole).Error; err != nil {
			log.Printf("Failed to create admin role: %v", err)
		}
	}
}

// createDefaultAdminUser 创建默认管理员用户
func createDefaultAdminUser(db *gorm.DB) {
	var adminUser models.User
	result := db.Where("username = ?", "admin").First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// 获取管理员角色
		var adminRole models.Role
		if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
			log.Printf("Failed to find admin role: %v", err)
			return
		}

		// 创建管理员用户
		adminUser = models.User{
			Username: "admin",
			Password: "$2a$10$XgXB8ZUhPmHXkXG0R5yn6.9VRCHqw0Zy8OgIm1vGAn.YQNi8rKKyq", // 默认密码：admin123
			Email:    "admin@xzyq.com",
			FullName: "System Administrator",
			IsActive: true,
			Roles:    []models.Role{adminRole},
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		}
	}
}
