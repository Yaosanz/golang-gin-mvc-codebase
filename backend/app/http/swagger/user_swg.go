package swagger

// swagger:model UserResponse
type UserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// swagger:model UserListResponse
type UserListResponse struct {
	Contents []UserResponse `json:"contents"`
}

// swagger:model UserPaginatedResponse
type UserPaginatedResponse struct {
	Contents   []UserResponse `json:"contents"`
	Pagination Pagination     `json:"pagination"`
}
