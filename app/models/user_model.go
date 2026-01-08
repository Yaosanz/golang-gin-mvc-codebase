package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        uuid.UUID     	`gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time			`gorm:"autoCreateTime"`
	UpdatedAt time.Time			`gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt 	`gorm:"index"`
	Role      string         `gorm:"type:varchar(50);default:'user'" json:"role"`
	Name      string		 `gorm:"type:varchar(100);not null" json:"name"`
	Username  string		 `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string		 `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone     *string		`gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
	Password  string		`gorm:"type:varchar(255);not null" json:"-"`
	IsActive  bool		 	`gorm:"default:true" json:"is_active"`
	RoleID uuid.UUID 		`gorm:"type:uuid"`
}

// Hook BeforeCreate → generate UUID otomatis sebelum insert
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// TableName override
func (u *User) TableName() string {
	return "users"
}
