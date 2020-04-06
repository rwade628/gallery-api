package http_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	router "github.com/rwade628/gallery-api/http"
	"github.com/rwade628/gallery-api/server"
)

var (
	gallery router.Gallery
	r       *gin.Engine
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setup(dir string) *gin.Engine {
	dbpath := filepath.Join(dir, "gallery.db")

	db, _ := storm.Open(dbpath)
	defer db.Close()

	gallery = router.Gallery{
		ID:   1,
		Name: "test-gallery",
		Type: "photo",
		Tags: "tag1,tag2",
		Files: []router.File{
			{
				Src:    "/some/path",
				Width:  1,
				Height: 2,
			},
		},
	}

	err := db.Init(&router.Gallery{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Save(&gallery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Saved gallery to", dbpath)

	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()
	server.Setup(r, ".", dbpath)

	return r
}
