package main

import (
	"io"

	"github.com/gin-gonic/gin"
)

const DataDir = "./data/"

func main() {
	r := gin.Default()

	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")

		data, err := readEntry(id)
		if err != nil {
			c.JSON(404, gin.H{
				"error":    "not found",
				"full_err": err,
			})
			return
		}

		file, err := loadFile(id)
		if err != nil {
			c.JSON(404, gin.H{
				"error":    "not found",
				"full_err": err,
			})
			return
		}

		c.Writer.Header().Set("Content-Type", data.Mime)
		io.Copy(c.Writer, file)
	})

	r.POST("/upload", handlePOST)

	r.Run()
}
