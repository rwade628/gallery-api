package http_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	router "github.com/rwade628/gallery-api/http"
)

func TestGetV1Galleries(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "testing")
	defer os.RemoveAll(dir)

	r := setup(dir)

	w := performRequest(r, "GET", "/v1/galleries", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected response status code %d to be %d", w.Code, http.StatusOK)
	}

	var response []router.Gallery
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatalf("expected json unmarshal error %e to be nil", err)
	}

	if len(response) != 1 && gallery.Name != response[0].Name {
		t.Fatalf("expected galleries length %d to be %d and first gallery's name %s to be %s",
			len(response), 1, gallery.Name, response[0].Name)
	}
}

func TestPostV1Galleries(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "testing")
	defer os.RemoveAll(dir)

	r := setup(dir)

	newGallery := router.Gallery{
		ID:        1,
		Name:      "new-gallery",
		Type:      "movie",
		Length:    0,
		CreatedAt: time.Now(),
		Tags:      "tag1,tag2",
		Files: []router.File{
			{
				Src:    "/some/path",
				Width:  1,
				Height: 2,
			},
		},
	}
	galleryBytes, err := newGallery.MarshalJSON() // tags need to get marshaled properly
	if err != nil {
		t.Fatalf("expected newGallery marshals error %e to be nil", err)
	}

	w := performRequest(r, "POST", "/v1/galleries", bytes.NewBuffer(galleryBytes))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected response status code %d to be %d", w.Code, http.StatusCreated)
	}

	var response router.Gallery
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatalf("expected body unmsrhal error %e to be nil", err)
	}

	if newGallery.Name != response.Name {
		t.Fatalf("expected response gallery name %s to be %s", response.Name, newGallery.Name)
	}
}

func TestPutV1Galleries(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "testing")
	defer os.RemoveAll(dir)

	r := setup(dir)

	newGallery := router.Gallery{
		ID:        1,
		Name:      "test-gallery",
		Type:      "photo",
		Length:    0,
		CreatedAt: time.Now(),
		Tags:      "tag1,tag2",
		Files: []router.File{
			{
				Src:    "/some/path",
				Width:  1,
				Height: 2,
			},
		},
	}
	galleryBytes, err := newGallery.MarshalJSON() // tags need to get marshaled properly
	if err != nil {
		t.Fatalf("expected newGallery marshals error %e to be nil", err)
	}

	w := performRequest(r, "PUT", "/v1/galleries", bytes.NewBuffer(galleryBytes))
	if w.Code != http.StatusOK {
		t.Fatalf("expected response status code %d to be %d", w.Code, http.StatusOK)
	}

	var response router.Gallery
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatalf("expected body unmsrhal error %e to be nil", err)
	}

	if newGallery.Name != response.Name {
		t.Fatalf("expected response gallery name %s to be %s", response.Name, newGallery.Name)
	}
}
