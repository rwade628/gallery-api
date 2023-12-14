package http_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/asdine/storm"

	router "github.com/rwade628/gallery-api/http"
)

var galleries []router.Gallery

func filterSetup(t *testing.T) *storm.DB {
	dir, _ := ioutil.TempDir(os.TempDir(), "testing")
	defer os.RemoveAll(dir)

	dbpath := filepath.Join(dir, "gallery.db")

	db, _ := storm.Open(dbpath)

	gallery1 := router.Gallery{
		Name:      "test1",
		Type:      "photo",
		Tags:      "tag1,tag2",
		CreatedAt: time.Now(),
		Length:    1,
		Views:     0,
	}
	galleries = append(galleries, gallery1)

	err := db.Save(&gallery1)
	if err != nil {
		t.Fatalf("expected db.Save error %e to be nil", err)
	}

	gallery2 := router.Gallery{
		Name:      "test2",
		Type:      "movie",
		Tags:      "tag2,tag3",
		CreatedAt: time.Now(),
		Length:    2,
		Views:     1,
	}
	galleries = append(galleries, gallery2)

	err = db.Save(&gallery2)
	if err != nil {
		t.Fatalf("expected db.Save error %e to be nil", err)
	}

	gallery3 := router.Gallery{
		Name:      "test3",
		Type:      "movie",
		Tags:      "tag4",
		CreatedAt: time.Now(),
		Length:    0,
		Views:     1,
	}
	galleries = append(galleries, gallery3)

	err = db.Save(&gallery3)
	if err != nil {
		t.Fatalf("expected db.Save error %e to be nil", err)
	}

	return db
}

func TestSelectFiltersType(t *testing.T) {
	db := filterSetup(t)
	defer db.Close()

	queryOptions := []url.Values{}
	queryResults := []int{}

	all := url.Values{}
	all.Add("type", "all")
	queryOptions = append(queryOptions, all)
	queryResults = append(queryResults, 3)

	photo := url.Values{}
	photo.Add("type", "photo")
	queryOptions = append(queryOptions, photo)
	queryResults = append(queryResults, 1)

	movie := url.Values{}
	movie.Add("type", "movie")
	queryOptions = append(queryOptions, movie)
	queryResults = append(queryResults, 2)

	for i, query := range queryOptions {
		dbSelect, err := router.SelectFilters(query, db)
		if err != nil {
			t.Fatalf("expected SelectFilters error %e to be nil for test %d", err, i)
		}

		var results []router.Gallery
		err = dbSelect.Find(&results)
		if err != nil {
			t.Fatalf("expected dbSelect.Find error %e to be nil for test %d", err, i)
		}

		if len(results) != queryResults[i] {
			t.Fatalf("expected results length %d to be %d for test %d", len(results), queryResults[i], i)
		}
	}
}

func TestSelectFiltersMultipleKeys(t *testing.T) {
	db := filterSetup(t)
	defer db.Close()

	queries := url.Values{}
	queries.Add("type", "one")
	queries.Add("type", "two")

	_, err := router.SelectFilters(queries, db)

	if err == nil {
		t.Fatalf("expected SelectFilters error %e to not be nil", err)
	}
}

func TestSelectFiltersTags(t *testing.T) {
	db := filterSetup(t)
	defer db.Close()

	queryOptions := []url.Values{}
	queryResults := []int{}

	tag1 := url.Values{}
	tag1.Add("tag", "tag1")
	queryOptions = append(queryOptions, tag1)
	queryResults = append(queryResults, 1)

	tag2 := url.Values{}
	tag2.Add("tag", "tag2")
	queryOptions = append(queryOptions, tag2)
	queryResults = append(queryResults, 2)

	tag3 := url.Values{}
	tag3.Add("tag", "tag3")
	queryOptions = append(queryOptions, tag3)
	queryResults = append(queryResults, 1)

	multiTag := url.Values{}
	multiTag.Add("tag", "tag1")
	multiTag.Add("tag", "tag3")
	queryOptions = append(queryOptions, multiTag)
	queryResults = append(queryResults, 2)

	for i, query := range queryOptions {
		dbSelect, err := router.SelectFilters(query, db)
		if err != nil {
			t.Fatalf("expected SelectFilters error %e to be nil for test %d", err, i)
		}

		var results []router.Gallery
		err = dbSelect.Find(&results)
		if err != nil {
			t.Fatalf("expected dbSelect.Find error %e to be nil for test %d", err, i)
		}

		if len(results) != queryResults[i] {
			t.Fatalf("expected results length %d to be %d for test %d", len(results), queryResults[i], i)
		}
	}
}

func TestSelectFiltersOrder(t *testing.T) {
	db := filterSetup(t)
	defer db.Close()

	queryOptions := []url.Values{}
	queryResults := []router.Gallery{}

	query := url.Values{}
	query.Add("orderBy", "newest")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[2])

	query = url.Values{}
	query.Add("orderBy", "oldest")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[0])

	query = url.Values{}
	query.Add("orderBy", "longest")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[1])
	//should not include length=0 galleries

	query = url.Values{}
	query.Add("orderBy", "shortest")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[0])

	query = url.Values{}
	query.Add("orderBy", "mostViews")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[2])

	query = url.Values{}
	query.Add("orderBy", "leastViews")
	queryOptions = append(queryOptions, query)
	queryResults = append(queryResults, galleries[0])

	for i, query := range queryOptions {
		dbSelect, err := router.SelectFilters(query, db)
		if err != nil {
			t.Fatalf("expected SelectFilters error %e to be nil for test %d", err, i)
		}

		var results []router.Gallery
		err = dbSelect.Find(&results)
		if err != nil {
			t.Fatalf("expected dbSelect.Find error %e to be nil for test %d", err, i)
		}

		if results[0].Name != queryResults[i].Name {
			t.Fatalf("expected first result %s to be %s for test %d", results[0].Name, queryResults[i].Name, i)
		}
	}
}
