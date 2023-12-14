package http

import (
	"net/http"

	"github.com/asdine/storm"
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

		api.GET("/reindex", ReIndex)
	}
}

func ReIndex(c *gin.Context) {
	var err error

	dbpath := c.MustGet("DBPath").(string)

	//ignore error, its okay if the file doesn't exist yet we'll create it
	db, _ := storm.Open(dbpath)
	defer db.Close()

	err = db.Init(&Gallery{})
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	err = db.ReIndex(&Gallery{})
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Reindex complete")
}
