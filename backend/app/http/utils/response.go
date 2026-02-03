package utils

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	TotalData   int  `json:"total_data"`
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	EndOfPage   bool `json:"end_of_page"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// WithPagination Wrapper for pagination data
type WithPagination struct {
	Contents   interface{}
	Pagination Pagination
}

// SendError response
func SendError(ctx *gin.Context, statusCode int, message string, err interface{}) {
	ctx.JSON(statusCode, ApiResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
	ctx.Abort()
}

// NewPagination builds a pagination object and calculates total_page + last_page
func NewPagination(totalData, currentPage, perPage int) Pagination {
	if perPage <= 0 {
		perPage = 10 // default to 10 if not set
	}
	if currentPage <= 0 {
		currentPage = 1 // default to 1 if not set
	}

	lastPage := (totalData + perPage - 1) / perPage

	return Pagination{
		TotalData:   totalData,
		TotalPages:  lastPage,
		CurrentPage: currentPage,
		PerPage:     perPage,
		EndOfPage:   currentPage >= lastPage,
	}
}

// SendSuccess automatically formats response based on data type
// - struct/map -> single
// - slice/array -> list
// - WithPagination -> paginated list
func SendSuccess(ctx *gin.Context, data interface{}, message string) {
	switch v := data.(type) {

	// If data is wrapped with pagination
	case WithPagination:
		ctx.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Message: message,
			Data: gin.H{
				"contents":   v.Contents,
				"pagination": v.Pagination,
			},
		})
		return

	default:
		rv := reflect.ValueOf(data)
		kind := rv.Kind()

		// Slice or Array = list
		if kind == reflect.Slice || kind == reflect.Array {
			ctx.JSON(http.StatusOK, ApiResponse{
				Success: true,
				Message: message,
				Data: gin.H{
					"contents": data,
				},
			})
			return
		}

		// Struct, Map, or anything else = single
		ctx.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Message: message,
			Data:    data,
		})
	}
}

// SendOne send response format for single data
// This method is used to send a success response (200 OK by default) with a single data object.
// If needed, you can customize the HTTP status with the `statusCode` parameter.
func SendOne(ctx *gin.Context, data interface{}, message string, statusCode *int) {
	status := http.StatusOK
	if statusCode != nil {
		status = *statusCode
	}

	ctx.JSON(status, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendList send response format for list of data (without pagination)
// This method is used to send a success response (200 OK by default) with list of data object.
// If needed, you can customize the HTTP status with the `statusCode` parameter.
func SendList(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, ApiResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"contents": data,
		},
	})
}

// SendPaginated send response format for list of data with pagination
// This method is used to send a success response (200 OK by default) with list of data object with pagination.
// If needed, you can customize the HTTP status with the `statusCode` parameter.
func SendPaginated(ctx *gin.Context, contents interface{}, pagination Pagination, message string) {
	ctx.JSON(http.StatusOK, ApiResponse{
		Success: true,
		Message: message,
		Data: gin.H{
			"contents":   contents,
			"pagination": pagination,
		},
	})
}
