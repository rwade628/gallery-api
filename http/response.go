package http

import "github.com/gin-gonic/gin"

// Wraps the response 'data' object in a map[string].
type r map[string]interface{}

func writeError(c *gin.Context, status int, err error) {
	c.Status(status)
	c.Error(err).SetType(gin.ErrorTypePublic)
}
