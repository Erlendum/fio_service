package v1

import (
	"github.com/gin-gonic/gin"
)

type Resposne struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, Resposne{Message: message})
}
