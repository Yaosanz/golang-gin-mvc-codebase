package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement;column:id;type:bigint" example:"1"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;column:created_at;type:timestamp" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;column:updated_at;type:timestamp" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;column:deleted_at;type:timestamp" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
}

type ActionLog struct {
	CreatedBy *int64 `gorm:"column:created_by;type:bigint" json:"created_by,omitempty"` // bigint, nullable
	UpdatedBy *int64 `gorm:"column:updated_by;type:bigint" json:"updated_by,omitempty"` // bigint, nullable
	DeletedBy *int64 `gorm:"column:deleted_by;type:bigint" json:"deleted_by,omitempty"` // bigint, nullable
}
