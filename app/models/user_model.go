package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleCMS   = "cms"
	TokenTypeUser = "user"
	TokenTypeCMS  = "cms"
)

type User struct {
	ID        uuid.UUID     	`gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time			`gorm:"autoCreateTime"`
	UpdatedAt time.Time			`gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt 	`gorm:"index"`
	// Legacy role field for backward compatibility
	Role      string         `gorm:"type:varchar(50);default:'user'" json:"role"`
	Name      string		 `gorm:"type:varchar(100);not null" json:"name"`
	Username  string		 `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string		 `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone     *string		`gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
	Password  string		`gorm:"type:varchar(255);not null" json:"-"`
	IsActive  bool		 	`gorm:"default:true" json:"is_active"`
	// Legacy role ID for backward compatibility
	RoleID    uuid.UUID 		`gorm:"type:uuid" json:"-"`

	// Multi-role relationships
	UserRoles []UserRole     `gorm:"foreignKey:UserID" json:"-"`
	Roles     []Role         `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles,omitempty"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	User   User      `gorm:"foreignKey:UserID" json:"-"`
	Role   Role      `gorm:"foreignKey:RoleID" json:"-"`
}

// Hook BeforeCreate → generate UUID otomatis sebelum insert
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// Hook BeforeCreate for UserRole
func (ur *UserRole) BeforeCreate(tx *gorm.DB) (err error) {
	if ur.ID == uuid.Nil {
		ur.ID = uuid.New()
	}
	return
}

// GetPrimaryRole returns the first role for backward compatibility
func (u *User) GetPrimaryRole() string {
	if len(u.Roles) > 0 {
		return u.Roles[0].Name
	}
	return u.Role // fallback to legacy field
}

// GetPrimaryRoleID returns the first role ID for backward compatibility
func (u *User) GetPrimaryRoleID() uuid.UUID {
	if len(u.Roles) > 0 {
		return u.Roles[0].ID
	}
	return u.RoleID // fallback to legacy field
}

// GetAllRoleIDs returns all role IDs for the user
func (u *User) GetAllRoleIDs() []uuid.UUID {
	roleIDs := make([]uuid.UUID, len(u.Roles))
	for i, role := range u.Roles {
		roleIDs[i] = role.ID
	}
	return roleIDs
}

// GetAllRoleNames returns all role names for the user
func (u *User) GetAllRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}

// TableName override
func (u *User) TableName() string {
	return "users"
}

// TableName override for UserRole
func (ur *UserRole) TableName() string {
	return "user_roles"
}
