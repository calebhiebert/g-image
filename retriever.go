package main

import (
	"errors"
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

	if config.CacheSize == 0 {
		err = copyResponseFromS3(c.Writer, &data)
		if err != nil {
			c.JSON(500, gin.H{
				"error":    "reading",
				"full_err": err,
			})
		}
	} else {
		fileExists := fileExists(config.DataDir + data.ID)

		if fileExists {
			err = copyResponseFromCache(c.Writer, &data)
			if err != nil {
				c.JSON(500, gin.H{
					"error":    "error while reading file from cache",
					"full_err": err,
				})
			}
		} else {
			defer func() {
				go cacheCheck()
			}()

			err = copyAndCacheResponseFromS3(c.Writer, &data)
			if err != nil {
				c.JSON(500, gin.H{
					"error":    "error while downloading and caching file",
					"full_err": err,
				})
				return
			}
		}
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

func getObjectFile(data *Entry) (io.ReadCloser, error) {
	stat, err := os.Stat(config.DataDir + data.ID)
	if err != nil {
		return nil, err
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

		return nil, errors.New("File size mismatch")
	}

	return file, nil
}

func writeHeaders(writer gin.ResponseWriter, data *Entry) {
	writer.Header().Set("Content-Type", data.Mime)
	writer.Header().Set("Content-Length", strconv.FormatInt(data.Size, 10))
}

func copyResponseFromS3(writer gin.ResponseWriter, data *Entry) error {
	objectReader, err := getObjectReader(data.ID)
	if err != nil {
		return err
	}

	writer.Header().Set("X-Loaded-From", "S3")
	writeHeaders(writer, data)

	_, err = io.Copy(writer, objectReader)
	if err != nil {
		println(err)
	}

	objectReader.Close()

	return nil
}

func copyResponseFromCache(writer gin.ResponseWriter, data *Entry) error {
	fileReader, err := getObjectFile(data)
	if err != nil {
		return err
	}

	writer.Header().Set("X-Loaded-From", "S3")
	writeHeaders(writer, data)

	io.Copy(writer, fileReader)
	fileReader.Close()

	return nil
}

func copyAndCacheResponseFromS3(writer gin.ResponseWriter, data *Entry) error {
	objectReader, err := getObjectReader(data.ID)
	if err != nil {
		return err
	}

	ensureDirectory(config.DataDir)

	defer objectReader.Close()
	fileWriter, err := os.OpenFile(config.DataDir+data.ID, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, objectReader)
	if err != nil {
		fileWriter.Close()
		return err
	}

	fileWriter.Close()

	fileReader, err := os.Open(config.DataDir + data.ID)
	if err != nil {
		return err
	}

	defer fileReader.Close()

	writer.Header().Set("X-Loaded-From", "S3-cache")
	writeHeaders(writer, data)

	_, err = io.Copy(writer, fileReader)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}
