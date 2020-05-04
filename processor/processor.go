package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pillash/mp4util"

	router "github.com/rwade628/gallery-api/http"
)

var path string

func init() {
	usage := "Path that processor will search for files"
	flag.StringVar(&path, "path", ".", usage)
	flag.StringVar(&path, "p", ".", usage+" (shorthand)")
}

func main() {
	flag.Parse()
	fmt.Println("Reading files at", path)

	err := listFiles(path)
	if err != nil {
		panic(err)
	}
}

func listFiles(rootPath string) error {
	var galleries []router.Gallery

	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Failed walking path")
				return err
			}
			ext := filepath.Ext(path)
			if ext == ".jpg" || ext == ".jpeg" {
				moviePath := strings.ReplaceAll(path, "jpg", "mp4")

				_, temperr := os.Stat(moviePath)
				if os.IsNotExist(temperr) {
					// jpg is not thumbnail of mp4 so process the gallery
					setPath := strings.Split(path, "/"+info.Name())[0]
					err = addPhotoSet(setPath, rootPath, &galleries)
				}
			} else if ext == ".mp4" {
				err = addMovie(path, rootPath, info, &galleries)
			}
			return err
		})

	if err != nil {
		return err
	}

	res, err := http.Get("http://localhost:8080/v1/galleries")
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var existingGalleries []router.Gallery

	if res.StatusCode == 200 {
		err = json.NewDecoder(res.Body).Decode(&existingGalleries)
		if err != nil {
			return err
		}
	} else {
		// no galleries existing
		existingGalleries = []router.Gallery{}
	}

	return insertGalleries(galleries, existingGalleries)
}

func insertGalleries(galleries, existingGalleries []router.Gallery) error {
	fmt.Println("inserting galleries", len(galleries))
	for _, gallery := range galleries {
		found := false

		for _, existingGallery := range existingGalleries {
			if gallery.Files[0].Src == existingGallery.Files[0].Src {
				fmt.Println("Found an existing gallery")

				found = true

				break
			}
		}

		if !found {
			galleryBytes, err := gallery.MarshalJSON() // tags need to get marshaled properly
			if err != nil {
				return err
			}

			res, err := http.Post("http://localhost:8080/v1/galleries", "application/json", bytes.NewBuffer(galleryBytes))
			if err != nil {
				fmt.Println("Error creating gallery", err.Error())
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			if res.StatusCode != http.StatusCreated {
				fmt.Println("Gallery not created", string(body))
			}
		}
	}

	return nil
}

func addMovie(path, rootPath string, info os.FileInfo, galleries *[]router.Gallery) error {
	thumbPath := strings.ReplaceAll(path, "mp4", "jpg")

	var files []router.File

	if reader, err := os.Open(thumbPath); err == nil {
		defer reader.Close()

		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			fmt.Println("unable to decode image")
			return err
		}

		thumbSrc, err := filepath.Rel(rootPath, thumbPath)
		if err != nil {
			fmt.Println("unable to get relative thumb path")
			return err
		}

		thumb := router.File{
			Src:    thumbSrc,
			Width:  im.Width,
			Height: im.Height,
		}
		files = append(files, thumb)

		src, err := filepath.Rel(rootPath, path)
		if err != nil {
			fmt.Println("unable to get relative movie path")
			return err
		}

		movie := router.File{
			Src:    src,
			Width:  im.Width,
			Height: im.Height,
		}
		files = append(files, movie)
	}

	duration, err := mp4util.Duration(path)
	if err != nil {
		fmt.Println("unable to get movie duration")
		duration = 0
	}

	noRoot := strings.SplitAfter(path, rootPath+"/")[1]
	tags := strings.Split(noRoot, "/")
	tags = tags[:len(tags)-1]

	gallery := router.Gallery{
		Name:      strings.Split(info.Name(), ".")[0],
		Length:    duration,
		CreatedAt: info.ModTime(),
		Type:      "movie",
		Tags:      strings.Join(tags, ","),
		Files:     files,
	}

	*galleries = append(*galleries, gallery)

	return nil
}
