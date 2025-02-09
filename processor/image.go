package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	router "github.com/rwade628/gallery-api/http"
)

func addPhotoSet(path, rootPath, fullPath string, modTime time.Time, galleriesPtr *[]router.Gallery, existingGalleries *[]router.Gallery) error {
	galleries := *galleriesPtr

	// fmt.Println(path, rootPath, fullPath)

	for _, gallery := range *existingGalleries {
		srcPath := fullPath
		if strings.Contains(rootPath, "/public") {
			srcPath = strings.TrimSuffix(rootPath, "/public")
		}
		src, err := filepath.Rel(srcPath, fullPath)
		if err != nil {
			return err
		}
		if src == gallery.Files[0].Src {
			// fmt.Println("found an existing gallery", src)
			*galleriesPtr = append(*galleriesPtr, gallery)
			return nil
		}
	}

	pathSplit := strings.Split(path, "/")
	name := pathSplit[len(pathSplit)-1]

	// fileInfo, err := getFileInfo(path)
	// if err != nil {
	// 	return err
	// }
	//
	// setName := fileInfo.Name()
	// fmt.Println(setName)

	for _, gallery := range galleries {
		if gallery.Name == name {
			return nil
		}
	}

	noRoot := strings.SplitAfter(path, rootPath+"/")[1]
	tags := strings.Split(noRoot, "/")
	tags = tags[:len(tags)-1]

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// fmt.Println("Adding set with modtime:", fileInfo.ModTime())

	gallery := router.Gallery{
		Name:      name,
		Length:    len(files),
		CreatedAt: modTime,
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

func getFiles(files []fs.DirEntry, path, rootPath string) ([]router.File, error) {
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
