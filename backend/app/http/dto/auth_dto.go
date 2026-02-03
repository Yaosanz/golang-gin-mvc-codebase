package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login success"`
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}
