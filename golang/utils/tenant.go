package utils

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

type Tenant struct {
    ID            int       `json:"id"`
    Name          string    `json:"name" binding:"required"`
    Code          string    `json:"code" binding:"required"`
    Address       string    `json:"address"`
    ContactPerson string    `json:"contact_person"`
    ContactPhone  string    `json:"contact_phone"`
    Email         string    `json:"email"`
    Status        bool      `json:"status"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// GetTenantList 获取租户列表
func GetTenantList(c *gin.Context) {
    rows, err := DB.Query(`
        SELECT id, name, code, address, contact_person, contact_phone, email, status, created_at, updated_at 
        FROM tenants
        ORDER BY id DESC
    `)
    if err != nil {
        ErrorLogger("tenant", "获取租户列表失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取租户列表失败"})
        return
    }
    defer rows.Close()

    var tenants []Tenant
    for rows.Next() {
        var t Tenant
        err := rows.Scan(&t.ID, &t.Name, &t.Code, &t.Address, &t.ContactPerson, &t.ContactPhone, &t.Email, &t.Status, &t.CreatedAt, &t.UpdatedAt)
        if err != nil {
            ErrorLogger("tenant", "扫描租户数据失败", &LogMetadata{
                ExtraInfo: map[string]interface{}{
                    "error": err.Error(),
                },
            })
            continue
        }
        tenants = append(tenants, t)
    }

    c.JSON(http.StatusOK, tenants)
}

// AddTenant 添加新租户
func AddTenant(c *gin.Context) {
    var tenant Tenant
    if err := c.ShouldBindJSON(&tenant); err != nil {
        ErrorLogger("tenant", "租户数据绑定失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户数据"})
        return
    }

    // 检查租户代码是否已存在
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tenants WHERE code = $1)", tenant.Code).Scan(&exists)
    if err != nil {
        ErrorLogger("tenant", "检查租户代码失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
        return
    }
    if exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "租户代码已存在"})
        return
    }

    // 插入新租户
    result := DB.QueryRow(`
        INSERT INTO tenants (name, code, address, contact_person, contact_phone, email, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `, tenant.Name, tenant.Code, tenant.Address, tenant.ContactPerson, tenant.ContactPhone, tenant.Email, tenant.Status)

    err = result.Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
    if err != nil {
        ErrorLogger("tenant", "添加租户失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "添加租户失败"})
        return
    }

    InfoLogger("tenant", "添加租户成功", &LogMetadata{
        Action: "create",
        UserID: c.GetString("username"),
        ExtraInfo: map[string]interface{}{
            "tenant_id":   tenant.ID,
            "tenant_code": tenant.Code,
        },
    })

    c.JSON(http.StatusOK, tenant)
}

// UpdateTenant 更新租户信息
func UpdateTenant(c *gin.Context) {
    id := c.Param("id")
    var tenant Tenant
    if err := c.ShouldBindJSON(&tenant); err != nil {
        ErrorLogger("tenant", "租户数据绑定失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户数据"})
        return
    }

    // 检查租户代码是否与其他租户重复
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tenants WHERE code = $1 AND id != $2)", tenant.Code, id).Scan(&exists)
    if err != nil {
        ErrorLogger("tenant", "检查租户代码失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
        return
    }
    if exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "租户代码已存在"})
        return
    }

    result, err := DB.Exec(`
        UPDATE tenants 
        SET name = $1, code = $2, address = $3, contact_person = $4, 
            contact_phone = $5, email = $6, status = $7, updated_at = CURRENT_TIMESTAMP
        WHERE id = $8
    `, tenant.Name, tenant.Code, tenant.Address, tenant.ContactPerson, 
       tenant.ContactPhone, tenant.Email, tenant.Status, id)

    if err != nil {
        ErrorLogger("tenant", "更新租户失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "更新租户失败"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
        return
    }

    InfoLogger("tenant", "更新租户成功", &LogMetadata{
        Action: "update",
        UserID: c.GetString("username"),
        ExtraInfo: map[string]interface{}{
            "tenant_id":   id,
            "tenant_code": tenant.Code,
        },
    })

    c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteTenant 删除租户
func DeleteTenant(c *gin.Context) {
    id := c.Param("id")

    // 获取租户信息用于日志记录
    var tenantCode string
    err := DB.QueryRow("SELECT code FROM tenants WHERE id = $1", id).Scan(&tenantCode)
    if err != nil && err != sql.ErrNoRows {
        ErrorLogger("tenant", "获取租户信息失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
        return
    }

    result, err := DB.Exec("DELETE FROM tenants WHERE id = $1", id)
    if err != nil {
        ErrorLogger("tenant", "删除租户失败", &LogMetadata{
            ExtraInfo: map[string]interface{}{
                "error": err.Error(),
            },
        })
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除租户失败"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
        return
    }

    InfoLogger("tenant", "删除租户成功", &LogMetadata{
        Action: "delete",
        UserID: c.GetString("username"),
        ExtraInfo: map[string]interface{}{
            "tenant_id":   id,
            "tenant_code": tenantCode,
        },
    })

    c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}