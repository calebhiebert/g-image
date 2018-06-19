package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/minio/minio-go"

	"github.com/gin-gonic/gin"
)

type HashMismatchError struct {
	code     string
	err      string
	original string
	current  string
}

func (e HashMismatchError) Error() string {
	return e.err
}

func getFile(c *gin.Context) {
	apiKey, _ := c.Get("apikey")

	if apiKey == nil || !apiKey.(APIKey).Read {
		c.JSON(401, gin.H{
			"error": "Missing read permissions",
			"code":  "PermError",
		})
		return
	}

	var data Entry
	var err error
	id := c.Param("id")

	if isWebhookSet() {
		data, err = webhookGetInfo(id)
		if err != nil {
			fmt.Println(err)
		}
	}

	if data.ID == "" {
		data, err = readEntry(id)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(404, gin.H{
					"error": "File not found in database",
					"code":  "NotFound",
				})
			} else {
				c.JSON(500, gin.H{
					"error":    "Error while doing local db lookup",
					"code":     "DbError",
					"full_err": err,
				})
			}
			return
		}
	}

	fileExists := fileExists(config.DataDir + data.ID)

	if fileExists {
		err := copyResponseFromCache(c.Writer, &data)
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

		err := copyAndCacheResponseFromS3(c.Writer, &data)
		if err != nil {
			handleRetrievalErr(err, c)
		}
	}
}

func getFileInfo(c *gin.Context) {
	apiKey, _ := c.Get("apikey")

	if apiKey == nil || !apiKey.(APIKey).Read {
		c.JSON(401, gin.H{
			"error": "Missing read permissions",
			"code":  "PermError",
		})
		return
	}

	var data Entry
	var err error
	id := c.Param("id")

	if isWebhookSet() {
		data, err = webhookGetInfo(id)
		if err != nil {
			fmt.Println(err)
		}
	}

	if data.ID == "" {
		data, err = readEntry(id)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(404, gin.H{
					"error": "File not found in database",
					"code":  "NotFound",
				})
			} else {
				c.JSON(500, gin.H{
					"error":    "Error while doing local db lookup",
					"code":     "DbError",
					"full_err": err,
				})
			}
			return
		}
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

	h := sha256.New()

	_, err = io.Copy(fileWriter, io.TeeReader(objectReader, h))
	if err != nil {
		fileWriter.Close()
		os.Remove(config.DataDir + data.ID)
		return err
	}
	fileWriter.Close()

	if hex.EncodeToString(h.Sum(nil)) != data.Sha256 {
		return &HashMismatchError{
			code:     "HashMismatch",
			err:      "The file downloaded from S3 was different than the original",
			original: data.Sha256,
			current:  hex.EncodeToString(h.Sum(nil)),
		}
	}

	fileReader, err := os.Open(config.DataDir + data.ID)
	if err != nil {
		return err
	}

	defer fileReader.Close()
	defer func() {
		cacheCheck()
	}()

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

func handleRetrievalErr(err error, c *gin.Context) {
	switch err.(type) {
	case minio.ErrorResponse:
		switch err.(minio.ErrorResponse).Code {
		case "NoSuchBucket":
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("A bucket with the name %s does not exist", config.BucketName),
				"code":  "NoSuchBucket",
			})
			break
		default:
			c.JSON(500, gin.H{
				"error": "S3 Error",
				"code":  err.(minio.ErrorResponse).Code,
			})
			break
		}
		break
	case *url.Error:
		fmt.Println("Url error type", reflect.TypeOf(err.(*url.Error).Err))
		c.JSON(500, gin.H{
			"error": err.(*url.Error).Error(),
			"code":  "CONNECTION_ERROR",
		})
		break
	case *HashMismatchError:
		c.JSON(500, gin.H{
			"code":     err.(*HashMismatchError).code,
			"original": err.(*HashMismatchError).original,
			"current":  err.(*HashMismatchError).current,
			"error":    err.Error(),
		})
		break
	default:
		fmt.Println("Error while getting file from bucket", reflect.TypeOf(err), err)
		c.JSON(500, gin.H{
			"error":    "error while downloading and caching file",
			"full_err": err,
			"code":     "GENERIC_ERR",
		})
		break
	}
}
