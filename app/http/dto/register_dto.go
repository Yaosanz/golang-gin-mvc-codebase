package dto

type RegisterDTO struct {
	Name     string `json:"name" binding:"required,min=3"`
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Phone    *string `json:"phone"`
	Password string `json:"password" binding:"required,min=8"`
}