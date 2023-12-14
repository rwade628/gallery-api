package http

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	galleries := []Gallery{}
	var err error

	dbpath := c.MustGet("DBPath").(string)

	db, err := storm.Open(dbpath)
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

	tags := make(map[string][]Gallery)

	err = db.Select().Each(new(Gallery), func(record interface{}) error {
		u := record.(*Gallery)

		for _, tag := range strings.Split(u.Tags, ",") {
			tags[tag] = append(tags[tag], *u)
		}
		return nil
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	for key, tag := range tags {
		i := rand.Intn(len(tag))
		gallery := Gallery{
			Name: key,
			Files: []File{
				File{
					Src:    tag[i].Files[0].Src,
					Width:  tag[i].Files[0].Width,
					Height: tag[i].Files[0].Height,
				},
			},
			Type: "tag",
		}

		galleries = append(galleries, gallery)
	}

	c.JSON(http.StatusOK, galleries)
}
