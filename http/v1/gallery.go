package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/gin-gonic/gin"
)

var dbFile = "gallery.db"

func GetGalleries(c *gin.Context) {

	galleries := []Gallery{}
	var err error

	db, err := storm.Open(dbFile)
	defer db.Close()
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

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
		matchers := []q.Matcher{}
		for k, v := range queries {
			fmt.Printf("key[%s] value[%s]\n", k, v)
			if len(v) > 1 {
				writeError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("Can only have 1 type of %s", k)))
				return
			}
			val := v[0]
			switch k {
			case "type":
				matchers = append(matchers, q.Eq("Type", val))
			case "orderBy":
				switch val {
				case "longest", "shortest":
					matchers = append(matchers, q.Gt("Length", 0))
				}
			}
		}
		dbSelect := db.Select(matchers...)
		if len(queries["orderBy"]) == 1 {
			fmt.Println("OrderBy", queries["orderBy"][0])
			switch queries["orderBy"][0] {
			case "newest":
				dbSelect = dbSelect.OrderBy("CreatedAt").Reverse()
			case "oldest":
				fmt.Println("sorting by oldest")
				dbSelect = dbSelect.OrderBy("CreatedAt")
			case "longest":
				dbSelect = dbSelect.OrderBy("Length").Reverse()
			case "shortest":
				dbSelect = dbSelect.OrderBy("Length")
			case "mostViews":
				dbSelect = dbSelect.OrderBy("Views").Reverse()
			case "lestsViews":
				dbSelect = dbSelect.OrderBy("Views")

			}
		}

		err = dbSelect.Find(&galleries)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": galleries})

	// writeError(c, http.StatusBadRequest, errors.New("Must bind via query parameters"))
}

func AddGallery(c *gin.Context) {
	gallery := Gallery{}
	var err error

	if err = c.ShouldBindJSON(&gallery); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	db, _ := storm.Open(dbFile)
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

	c.JSON(http.StatusOK, gin.H{"data": gallery})
}

func UpdateGallery(c *gin.Context) {
	gallery := Gallery{}
	var err error

	if err = c.ShouldBindJSON(&gallery); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	db, _ := storm.Open(dbFile)
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

	c.JSON(http.StatusOK, gin.H{"data": gallery})
}
func DeleteGallery(c *gin.Context) {}
