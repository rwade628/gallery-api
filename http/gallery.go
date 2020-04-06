package http

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

func GetGalleries(c *gin.Context) {
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

	c.JSON(http.StatusOK, galleries)
}

func AddGallery(c *gin.Context) {
	gallery := Gallery{}

	var err error

	if err = c.ShouldBindJSON(&gallery); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	dbpath := c.MustGet("DBPath").(string)

	//ignore error, its okay if the file doesn't exist yet we'll create it
	db, _ := storm.Open(dbpath)
	defer db.Close()

	err = db.Init(&Gallery{})
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	err = db.Save(&gallery)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	//gin doesn't like the custom unmarshaler
	//we'll do it ourselves and send the string back
	bytes, err := gallery.MarshalJSON()
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusCreated, string(bytes))
}

func UpdateGallery(c *gin.Context) {
	gallery := Gallery{}

	var err error

	if err = c.ShouldBindJSON(&gallery); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

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

	err = db.Update(&gallery)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	//gin doesn't like the custom unmarshaler
	//we'll do it ourselves and send the string back
	bytes, err := gallery.MarshalJSON()
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, string(bytes))
}
func DeleteGallery(c *gin.Context) {}
