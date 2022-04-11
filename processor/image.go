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

	// fmt.Println("Adding set with modtime:", fileInfo.ModTime())

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
	defer file.Close()

	return file.Stat()
}

func getFiles(files []os.FileInfo, path, rootPath string) ([]router.File, error) {
	var filesToAdd []router.File

	for _, f := range files {
		if !strings.Contains(f.Name(), "th_") {
			filePath := filepath.Join(path, f.Name())
			reader, err := os.Open(filePath)
			if err == nil {

				im, _, err := image.DecodeConfig(reader)
				if err != nil {
					fmt.Println(path)
					return nil, err
				}

				srcPath := rootPath
				if strings.Contains(rootPath, "/public") {
					srcPath = strings.TrimSuffix(rootPath, "/public")
				}

				src, err := filepath.Rel(srcPath, filePath)
				if err != nil {
					return nil, err
				}

				thumbPath := filepath.Join(path, "th_"+f.Name())
				thumbSrc, err := filepath.Rel(srcPath, thumbPath)
				if err != nil {
					return nil, err
				}

				// fmt.Println("Adding gallery file with src:", src)

				image := router.File{
					Src:    src,
					Thumb:  thumbSrc,
					Width:  im.Width,
					Height: im.Height,
				}
				filesToAdd = append(filesToAdd, image)
			} else {
				return nil, err
			}
			reader.Close()
		}
	}

	return filesToAdd, nil
}
