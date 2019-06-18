package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rwade628/gallery-api/server"
)

var (
	r = gin.Default()
)

// Main function - starts up the server.
func main() {
	r.Run()
}

func init() {
	server.Setup(r)
}
