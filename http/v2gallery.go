package http

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

func GetV2Galleries(c *gin.Context) {
	galleries := []Gallery{}

	var err error

	dbpath := c.MustGet("DBPath").(string)

	db, err := storm.Open(dbpath)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	err = db.Init(&Gallery{})
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	queries := c.Request.URL.Query()
	fmt.Println(len(queries))

	if len(queries) == 0 {
		err = db.All(&galleries)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err)
			return
		}
	} else {
		dbSelect, err := SelectFilters(queries, db)
		if err != nil {
			writeError(c, http.StatusBadRequest, err)
			return
		}

		err = dbSelect.Find(&galleries)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err)
			return
		}
	}

	v2galleries := []Gallery{}
	for _, g := range galleries {
		g.Photos = g.Files
		g.Files = nil
		g.Src = g.Photos[0].Src
		g.Thumb = g.Photos[0].Thumb
		g.Width = g.Photos[0].Width
		g.Height = g.Photos[0].Height
		v2galleries = append(v2galleries, g)
	}

	c.JSON(http.StatusOK, v2galleries)
}
