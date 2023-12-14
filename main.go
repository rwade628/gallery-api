package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/rwade628/gallery-api/server"
)

var (
	path   = os.Getenv("FILE_PATH")
	dbpath = os.Getenv("DB_PATH")
)

var r *gin.Engine

// Main function - starts up the server.
func main() {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()

	fmt.Println("Serving files at", path)
	fmt.Println("Creating db file at", dbpath)

	server.Setup(r, path, dbpath)

	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
