package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	router "github.com/rwade628/gallery-api/http"
)

func TestGetV1Tags(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "testing")
	defer os.RemoveAll(dir)

	r := setup(dir)

	w := performRequest(r, "GET", "/v1/tags", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected response status code %d to be %d", w.Code, http.StatusOK)
	}

	var response []router.Gallery
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatalf("expected body unmsrhal error %e to be nil", err)
	}

	if len(response) != 2 {
		t.Fatalf("expected response length %d to be %d", len(response), 2)
	}
}
