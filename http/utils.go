package http

import (
	"fmt"
	"net/url"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func SelectFilters(queries url.Values, db *storm.DB) (storm.Query, error) {
	matchers := []q.Matcher{}
	for k, v := range queries {
		if len(v) > 1 && k != "tag" {
			return nil, fmt.Errorf("can only have 1 type of %s", k)
		}
		switch k {
		case "type":
			if v[0] != "all" {
				matchers = append(matchers, q.Eq("Type", v[0]))
			}
		case "tag":
			tags := []q.Matcher{}
			for _, val := range v {
				tags = append(tags, q.Re("Tags", fmt.Sprintf(".*%s.*", val)))
			}
			matchers = append(matchers, q.Or(tags...))
		case "orderBy":
			switch v[0] {
			case "longest", "shortest":
				matchers = append(matchers, q.Gt("Length", 0))
			}
		}
	}
	dbSelect := db.Select(matchers...)
	if len(queries["orderBy"]) == 1 {
		switch queries["orderBy"][0] {
		case "newest":
			dbSelect = dbSelect.OrderBy("CreatedAt").Reverse()
		case "oldest":
			dbSelect = dbSelect.OrderBy("CreatedAt")
		case "longest":
			dbSelect = dbSelect.OrderBy("Length").Reverse()
		case "shortest":
			dbSelect = dbSelect.OrderBy("Length")
		case "mostViews":
			dbSelect = dbSelect.OrderBy("Views").Reverse()
		case "leastViews":
			dbSelect = dbSelect.OrderBy("Views")
		}
	}
	return dbSelect, nil
}
