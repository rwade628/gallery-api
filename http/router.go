package http

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Define the API.
func Initialize(r *gin.Engine, path string) {
	r.Use(static.Serve("/public", static.LocalFile(path, false)))
	v1 := r.Group("/v1")
	{
		v1.GET("/galleries", GetGalleries)
		v1.POST("/galleries", AddGallery)
		v1.PUT("/galleries", UpdateGallery)
		v1.DELETE("/galleries", DeleteGallery)

		v1.GET("/tags", GetTags)

		v1.GET("/reindex", ReIndex)
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/galleries", GetV2Galleries)
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
