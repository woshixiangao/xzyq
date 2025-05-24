package utils

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

var db *sql.DB

func InitLogger(database *sql.DB) {
    db = database
}

type LogMetadata struct {
    UserID    string `json:"user_id,omitempty"`
    IP        string `json:"ip,omitempty"`
    Action    string `json:"action,omitempty"`
    Status    int    `json:"status,omitempty"`
    ExtraInfo map[string]interface{} `json:"extra_info,omitempty"`
}

func logToDB(level, component string, message string, metadata *LogMetadata) error {
    metadataJSON, err := json.Marshal(metadata)
    if err != nil {
        return fmt.Errorf("marshal metadata failed: %v", err)
    }

    _, err = db.Exec(
        "INSERT INTO system_logs (level, component, message, metadata) VALUES ($1, $2, $3, $4)",
        level,
        component,
        message,
        metadataJSON,
    )
    return err
}

func InfoLogger(component string, message string, metadata *LogMetadata) {
    if err := logToDB("INFO", component, message, metadata); err != nil {
        fmt.Printf("Failed to write info log: %v\n", err)
    }
}

func ErrorLogger(component string, message string, metadata *LogMetadata) {
    if err := logToDB("ERROR", component, message, metadata); err != nil {
        fmt.Printf("Failed to write error log: %v\n", err)
    }
}

func DbLogger(component string, message string, metadata *LogMetadata) {
    if err := logToDB("DB", component, message, metadata); err != nil {
        fmt.Printf("Failed to write db log: %v\n", err)
    }
}

// GetLogs 获取系统日志列表
func GetLogs(c *gin.Context, db *sql.DB) {
    // 获取查询参数
    level := c.Query("level")
    component := c.Query("component")
    startDate := c.Query("startDate")
    endDate := c.Query("endDate")
    keyword := c.Query("keyword")

    // 构建基础查询
    query := `SELECT id, level, component, message, metadata, created_at 
              FROM system_logs 
              WHERE 1=1`
    var params []interface{}
    paramCount := 1

    // 添加过滤条件
    if level != "" {
        query += fmt.Sprintf(" AND level = $%d", paramCount)
        params = append(params, level)
        paramCount++
    }

    if component != "" {
        query += fmt.Sprintf(" AND component = $%d", paramCount)
        params = append(params, component)
        paramCount++
    }

    if startDate != "" {
        query += fmt.Sprintf(" AND created_at >= $%d", paramCount)
        params = append(params, startDate)
        paramCount++
    }

    if endDate != "" {
        query += fmt.Sprintf(" AND created_at <= $%d", paramCount)
        params = append(params, endDate)
        paramCount++
    }

    if keyword != "" {
        query += fmt.Sprintf(" AND (message ILIKE $%d OR metadata::text ILIKE $%d)", paramCount, paramCount)
        searchPattern := "%" + keyword + "%"
        params = append(params, searchPattern)
        paramCount++
    }

    // 添加排序和分页
    query += " ORDER BY created_at DESC LIMIT 100"

    // 执行查询
    rows, err := db.Query(query, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "查询日志失败"})
        return
    }
    defer rows.Close()

    // 构建结果
    var logs []gin.H
    for rows.Next() {
        var (
            id int
            level, component, message string
            metadataBytes []byte
            createdAt time.Time
        )

        if err := rows.Scan(&id, &level, &component, &message, &metadataBytes, &createdAt); err != nil {
            continue
        }

        var metadata map[string]interface{}
        if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
            metadata = make(map[string]interface{})
        }

        logs = append(logs, gin.H{
            "id": id,
            "level": level,
            "component": component,
            "message": message,
            "metadata": metadata,
            "created_at": createdAt.Format("2006-01-02 15:04:05"),
        })
    }

    c.JSON(http.StatusOK, logs)
}