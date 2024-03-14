package Shared

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
	Error      string `json:"error,omitempty"`
}

// SendApiResponse - a generic helper function to send API responses.
func SendApiResponse[T any](w http.ResponseWriter, statusCode int, message string, data T, errMessage string) {
	response := ApiResponse[T]{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Error:      errMessage,
	}
	w.WriteHeader(statusCode) // Set the status code header
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		// Log the error, as there's little else we can do if encoding fails
		fmt.Println(encodeErr)
	}
}

func SendGinGenericApiResponse[T any](c *gin.Context, statusCode int, message string, data T, errMessage string) {
	c.JSON(statusCode, ApiResponse[T]{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Error:      errMessage,
	})
}
