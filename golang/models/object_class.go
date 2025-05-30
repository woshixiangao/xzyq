package models

import "time"

type ObjectClass struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	OrgID       uint      `json:"org_id" gorm:"not null"`
	ParentID    *uint     `json:"parent_id"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Organization  Organization  `json:"organization" gorm:"foreignKey:OrgID"`
	Parent        *ObjectClass  `json:"parent" gorm:"foreignKey:ParentID"`
	Children      []ObjectClass `json:"children" gorm:"foreignKey:ParentID"`
	CreatedByUser User          `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

// TableName 指定表名
func (ObjectClass) TableName() string {
	return "object_class"
}
