package models

import (
	"github.com/google/uuid"
	"time"
)

type Permission struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time		`gorm:"default:now()"`
}

func (Permission) TableName() string {
	return "permissions"
}