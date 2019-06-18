package v1

import (
	"encoding/json"
	"errors"
	"time"
)

type Gallery struct {
	ID        int         `storm:"id,increment"`
	Name      string      `json:"name" storm:"index"`
	Location  string      `json:"location" storm:"unique"`
	CreatedAt time.Time   `json:"createdAt" storm:"index"`
	Views     int         `json:"views" storm:"index"`
	Type      GalleryType `json:"type" storm:"index"`
	Tags      []string    `json:"tags" storm:"index"`
	Length    int         `json:"length" storm:"index"`
}

func (g *Gallery) UnmarshalJSON(data []byte) error {
	// Define a secondary type so that we don't end up with a recursive call to json.Unmarshal
	type Aux Gallery
	var a *Aux = (*Aux)(g)
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	// Validate the valid enum values
	switch g.Type {
	case movie, photo:
		return nil
	default:
		g.Type = ""
		return errors.New("invalid value for Key")
	}
}

type GalleryType string

const (
	movie GalleryType = "movie"
	photo GalleryType = "photo"
)

// type OrderBy string
//
// const (
// 	newest     OrderBy = "newest"
// 	oldest     OrderBy = "oldest"
// 	mostViews  OrderBy = "mostViews"
// 	leastViews OrderBy = "lestViews"
// 	largest    OrderBy = "largest"
// 	smallest   OrderBy = "smallest"
// )
//
// type Query struct {
// 	OrderBy string   `json:"orderBy"`
// 	Type    string   `json:"type"`
// 	Tags    []string `json:"tags"`
// }
