package dto

type CreateUserDTO struct {
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Username string  `json:"username" validate:"required,min=2,max=255" example:"John Doe"`
	Email    string  `json:"email" validate:"required,email" example:"john.doe@email.com"`
	Phone    *string `json:"phone" validate:"omitempty,e164" example:"+1234567890"`
	IsActive *bool   `json:"is_active" validate:"omitempty" example:"true"`
	Password string  `json:"password" validate:"required,min=8,max=255" example:"strongpassword"`
}

type UpdateUserDTO struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=100"`
	Username *string `json:"username" validate:"omitempty,min=2,max=255" example:"John Doe"`
	Email    *string `json:"email" validate:"omitempty,email" example:""`
	Phone    *string `json:"phone" validate:"omitempty,e164" example:"+1234567890"`
	IsActive *bool   `json:"is_active" validate:"omitempty" example:"true"`
	Password *string `json:"password" validate:"omitempty,min=8,max=255" example:"strongpassword"`
}

type ResponseUserDTO struct {
	ID       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"john.doe"`
	Email    string `json:"email" example:"john.doe@gmail.com"`
}
