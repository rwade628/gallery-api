package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rwade628/gallery-api/server"
)

var r *gin.Engine

// Main function - starts up the server.
func main() {
	flag.Parse()
	fmt.Println("Serving files at", path)
	fmt.Println("Creating db file at", dbpath)

	server.Setup(r, path, dbpath)

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

var path, dbpath string

func init() {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()

	usage := "Path that processor will search for files"
	flag.StringVar(&path, "path", ".", usage)
	flag.StringVar(&path, "p", ".", usage+" (shorthand)")

	dbusage := "Path that db file will be created at"
	flag.StringVar(&dbpath, "dbpath", "gallery.db", dbusage)
	flag.StringVar(&dbpath, "dbp", "gallery.db", dbusage+" (shorthand)")
}
