package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getFile(c *gin.Context) {
	apiKey, _ := c.Get("apikey")

	if !apiKey.(APIKey).Read {
		c.JSON(401, gin.H{
			"error": "Missing read permissions",
		})
		return
	}

	var data Entry
	id := c.Param("id")

	data, err := webhookGetInfo(id)
	if err != nil {
		fmt.Println(err)
	}

	if data.ID == "" {
		data, err = readEntry(id)
		if err != nil {
			c.JSON(404, gin.H{
				"error":    "not found",
				"full_err": err,
			})
			return
		}
	}

	if data.ID == "" {
		c.JSON(404, gin.H{
			"error": "not found",
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

	c.Writer.Header().Set("Content-Type", data.Mime)
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(data.Size, 10))
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		println(err)
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
		}

		return nil, err
	}
	return file, nil
}
