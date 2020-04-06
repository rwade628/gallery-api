package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	router "github.com/rwade628/gallery-api/http"
)

func addPhotoSet(path, rootPath string, galleriesPtr *[]router.Gallery) error {

	galleries := *galleriesPtr

	fileInfo, err := getFileInfo(path)
	if err != nil {
		return err
	}

	setName := fileInfo.Name()

	for _, gallery := range galleries {
		if gallery.Name == setName {
			return nil
		}
	}

	noRoot := strings.SplitAfter(path, rootPath+"/")[1]
	tags := strings.Split(noRoot, "/")
	tags = tags[:len(tags)-1]

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	gallery := router.Gallery{
		Name:      setName,
		Length:    len(files),
		CreatedAt: fileInfo.ModTime(),
		Type:      "photo",
		Tags:      strings.Join(tags, ","),
	}

	filesToAdd, err := getFiles(files, path, rootPath)
	if err != nil {
		return err
	}

	gallery.Files = filesToAdd

	*galleriesPtr = append(*galleriesPtr, gallery)

	return nil
}

func getFileInfo(path string) (os.FileInfo, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return nil, err
	}

	return file.Stat()
}

func getFiles(files []os.FileInfo, path, rootPath string) ([]router.File, error) {
	var filesToAdd []router.File

	for _, f := range files {
		filePath := filepath.Join(path, f.Name())
		if reader, err := os.Open(filePath); err == nil {
			defer reader.Close()

			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				return nil, err
			}

			src, err := filepath.Rel(rootPath, filePath)
			if err != nil {
				return nil, err
			}

			fmt.Println("Adding gallery file with src:", src)

			image := router.File{
				Src:    src,
				Width:  im.Width,
				Height: im.Height,
			}
			filesToAdd = append(filesToAdd, image)
		} else {
			return nil, err
		}
	}

	return filesToAdd, nil
}
