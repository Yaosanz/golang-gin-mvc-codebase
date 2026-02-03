package models

import (
	"time"

	"github.com/google/uuid"
)

type ShortenLink struct {
	ID          uuid.UUID 		`gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ShortCode   string			`gorm:"uniqueIndex;not null" json:"short_code"`
	OriginalURL string    		`gorm:"not null" json:"original_url"`
	UserID      uuid.UUID		`gorm:"type:uuid;not null" json:"user_id"`
	CreatedBy   uuid.UUID		`gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time		`gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time		`gorm:"default:now()" json:"updated_at"`
}


func (ShortenLink) TableName() string {
	return "shortenlink"
}
