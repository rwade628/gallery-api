package http

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Define the API.
func Initialize(r *gin.Engine, path string) {
	r.Use(static.Serve("/", static.LocalFile(path, false)))
	api := r.Group("/v1")
	{
		api.GET("/galleries", GetGalleries)
		api.POST("/galleries", AddGallery)
		api.PUT("/galleries", UpdateGallery)
		api.DELETE("/galleries", DeleteGallery)

		api.GET("/tags", GetTags)
	}
}
