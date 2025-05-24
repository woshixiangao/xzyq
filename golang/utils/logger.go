package utils

import (
    "database/sql"
    "encoding/json"
    "fmt"
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