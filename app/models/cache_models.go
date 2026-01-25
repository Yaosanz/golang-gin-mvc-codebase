package models

import (
	"github.com/google/uuid"
)

// UserCacheData represents a lightweight version of User for caching
// This avoids circular dependencies by not importing the models package
type UserCacheData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	IsActive bool   `json:"is_active"`
	Role     string `json:"role"`
	RoleID   string `json:"role_id"`
}

// NewUserCacheData creates a new UserCacheData from basic user information
func NewUserCacheData(id uuid.UUID, name, username, email string, phone *string, isActive bool, role string, roleID uuid.UUID) *UserCacheData {
	phoneStr := ""
	if phone != nil {
		phoneStr = *phone
	}

	return &UserCacheData{
		ID:       id.String(),
		Name:     name,
		Username: username,
		Email:    email,
		Phone:    phoneStr,
		IsActive: isActive,
		Role:     role,
		RoleID:   roleID.String(),
	}
}

// ToUserData converts cache data to basic user information
func (ucd *UserCacheData) ToUserData() (uuid.UUID, string, string, string, *string, bool, string, uuid.UUID, error) {
	userID, err := uuid.Parse(ucd.ID)
	if err != nil {
		return uuid.Nil, "", "", "", nil, false, "", uuid.Nil, err
	}

	roleID, err := uuid.Parse(ucd.RoleID)
	if err != nil {
		return uuid.Nil, "", "", "", nil, false, "", uuid.Nil, err
	}

	var phone *string
	if ucd.Phone != "" {
		phone = &ucd.Phone
	}

	return userID, ucd.Name, ucd.Username, ucd.Email, phone, ucd.IsActive, ucd.Role, roleID, nil
}
