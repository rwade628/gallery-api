package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetDBPath(p string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DBPath", p)
		c.Next()
	}
}

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		// If an error occured during handling the request, write the error as a JSON response.
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			c.JSON(c.Writer.Status(), gin.H{"error": err.Error()})
		}
	}
}
