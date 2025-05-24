package utils

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() *sql.DB {
	// 使用 PostgreSQL URL 格式
	connStr := "postgresql://myuser:mysecretpassword@47.121.141.235:5432/postgres?sslmode=disable"
	
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 测试连接
	err = DB.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	fmt.Println("数据库连接成功！")
	
	// 初始化日志器
	InitLogger(DB)

    // 创建system_logs表
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS system_logs (
            id SERIAL PRIMARY KEY,
            level VARCHAR(10) NOT NULL,
            component VARCHAR(50) NOT NULL,
            message TEXT NOT NULL,
            metadata JSONB,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal("创建system_logs表失败:", err)
    }

    // 创建索引以提高查询性能
    _, err = DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs(level);
        CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at);
        CREATE INDEX IF NOT EXISTS idx_system_logs_component ON system_logs(component)
    `)
    if err != nil {
        log.Fatal("创建system_logs索引失败:", err)
    }

    // 创建租户表
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS tenants (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            code VARCHAR(50) NOT NULL UNIQUE,
            address TEXT,
            contact_person VARCHAR(50),
            contact_phone VARCHAR(20),
            email VARCHAR(100),
            status BOOLEAN DEFAULT true,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal("创建租户表失败:", err)
    }

    return DB
}