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

	fileReader, err := loadObjectFile(&data)
	if err != nil {
		c.JSON(404, gin.H{
			"error":    "not found",
			"full_err": err,
		})
		return
	}
	defer fileReader.Close()
	defer func() {
		go cacheCheck()
	}()

	c.Writer.Header().Set("Content-Type", data.Mime)
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(data.Size, 10))
	_, err = io.Copy(c.Writer, fileReader)

	if err != nil {
		println(err)
	}
}

func getFileInfo(c *gin.Context) {
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

	c.JSON(200, data)
}

func loadObjectFile(data *Entry) (io.ReadCloser, error) {
	if config.CacheSize == 0 {
		objectReader, err := getObjectReader(data.ID)
		if err != nil {
			return nil, err
		}

		return objectReader, nil
	}

	stat, err := os.Stat(config.DataDir + data.ID)
	if err != nil {
		objectReader, err := getObjectReader(data.ID)
		if err != nil {
			return nil, err
		}

		return objectReader, nil
	}

	file, err := os.Open(config.DataDir + data.ID)
	if err != nil {
		return nil, err
	}

	if stat.Size() != data.Size {
		file.Close()
		err := os.Remove(config.DataDir + stat.Name())
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}
