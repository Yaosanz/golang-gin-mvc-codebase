package models

type User struct {
	BaseModel
	ActionLog

	Name     string  `gorm:"column:name;type:citext;not null" json:"name"` // citext, not nullable
	Username string  `gorm:"column:username;type:citext;not null" json:"username"`
	Email    string  `gorm:"column:email;type:text" json:"email"` // text, nullable
	Phone    *string `gorm:"column:phone;type:text" json:"phone,omitempty"`
	Password string  `gorm:"column:password;type:text" json:"password"`               // text, nullable
	IsActive *bool   `gorm:"column:is_active;not null;default:true" json:"is_active"` // boolean, default true
}

// TableName sets the insert table name for this struct type
func (m *User) TableName() string {
	return "users"
}
