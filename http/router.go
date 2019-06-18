package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/rwade628/gallery-api/http/v1"
)

// Define the API.
func Initialize(r *gin.Engine) {
	api := r.Group("")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello world!",
			})

		})
	}

	api = r.Group("/v1")
	{
		// api.GET("/galleries", v1.GetAllGalleries)
		api.GET("/galleries", v1.GetGalleries)
		api.POST("/galleries", v1.AddGallery)
		api.PUT("/galleries", v1.UpdateGallery)
		api.DELETE("/galleries", v1.DeleteGallery)

	}

}
