package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rwade628/gallery-api/http"
	"github.com/rwade628/gallery-api/middleware"
)

// Define all middlewares to use then set up the API.
func Setup(r *gin.Engine) {
	// Limit request size.
	r.Use(middleware.HandleErrors())

	http.Initialize(r)
}
