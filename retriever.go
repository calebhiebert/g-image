package main

import (
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getFile(c *gin.Context) {
	id := c.Param("id")

	data, err := readEntry(id)
	if err != nil {
		c.JSON(404, gin.H{
			"error":    "not found",
			"full_err": err,
		})
		return
	}

	file, err := loadObjectFile(data.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":    "server error",
			"full_err": err,
		})
		return
	}
	defer file.Close()

	println(file.Name())

	c.Writer.Header().Set("Content-Type", data.Mime)
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(data.Size, 10))
	bytes, err := io.Copy(c.Writer, file)

	if err != nil {
		println(err)
	} else {
		println("wrote bytes " + strconv.FormatInt(bytes, 10))
	}

}

func loadObjectFile(ID string) (*os.File, error) {
	file, err := os.Open(config.DataDir + ID)
	if err != nil {
		if config.BucketName != "" {
			err = downloadFile(ID)
			if err != nil {
				return nil, err
			}
			return loadObjectFile(ID)
		} else {
			return nil, err
		}
	}
	return file, nil
}
