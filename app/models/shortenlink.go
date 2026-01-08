package models

import (
	"time"
	"github.com/google/uuid"
)

type ShortenLink struct {
	ID          uuid.UUID 		`gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShortCode       string		`gorm:"uniqueIndex;not null"`
	OriginalURL 	string    	`gorm:"not null"`
	UserID      uuid.UUID		`gorm:"type:uuid;not null"`
	CreatedBy   uuid.UUID		`gorm:"type:uuid;not null"`
	CreatedAt   time.Time		`gorm:"default:now()"`
	UpdatedAt   time.Time		`gorm:"default:now()"`
}


func (ShortenLink) TableName() string {
	return "shortenlink"
}
