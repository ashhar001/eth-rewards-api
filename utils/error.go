package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleInternalServerError is a utility function that sends a JSON response with a 500 Internal Server Error status code.
// It takes a gin.Context object and a message string as parameters.
// The function constructs a JSON response with an "error" key and the provided message as its value.
func HandleInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}
