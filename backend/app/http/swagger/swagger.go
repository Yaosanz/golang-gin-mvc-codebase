package swagger

// swagger:model Response
type ResponseOk struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error,omitempty"`
}

// swagger:model Pagination
type Pagination struct {
	TotalRow    int64 `json:"total_row"`
	PerPage     int64 `json:"per_page"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	TotalPage   int64 `json:"total_page"`
}

type ResponseError struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"something went wrong!"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
