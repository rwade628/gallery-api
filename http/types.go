package http

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Gallery struct {
	ID         int         `storm:"id,increment"`
	Name       string      `json:"name" storm:"index"`
	Length     int         `json:"length" storm:"index"`
	CreatedAt  time.Time   `json:"createdAt" storm:"index"`
	Views      int         `json:"views" storm:"index"`
	Type       GalleryType `json:"type" storm:"index"`
	Tags       string      `json:"tags" storm:"index"`
	Preference string      `json:"preference" storm:"index"`
	Files      []File      `json:"files"`
}

type File struct {
	Src    string `json:"src"`
	Thumb  string `json:"thumb"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (g *Gallery) UnmarshalJSON(data []byte) error {
	// Define a secondary type so that we don't end up with a recursive call to json.Unmarshal
	type Alias Gallery

	aux := &struct {
		*Alias
		Tags []string `json:"tags"`
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	g.Tags = strings.Join(aux.Tags, ",")

	// Validate the valid enum values
	switch g.Type {
	case movie, photo, tag:
		return nil
	default:
		g.Type = ""
		return errors.New("invalid value for Key")
	}
}

func (g *Gallery) MarshalJSON() ([]byte, error) {
	type Alias Gallery

	return json.Marshal(&struct {
		*Alias
		Tags []string `json:"tags"`
	}{
		Tags:  strings.Split(g.Tags, ","),
		Alias: (*Alias)(g),
	})
}

type GalleryType string

const (
	movie GalleryType = "movie"
	photo GalleryType = "photo"
	tag   GalleryType = "tag"
)
