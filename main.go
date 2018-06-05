package main

import (
	"github.com/gin-gonic/gin"
)

const DataDir = "./data/"

func main() {
	r := gin.Default()
	r.GET("/:id", getFile)
	r.POST("/upload", putFile)
	r.Run()
}
