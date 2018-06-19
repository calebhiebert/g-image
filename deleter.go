package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func deleteFile(c *gin.Context) {
	apiKey, _ := c.Get("apikey")
	if apiKey == nil || !apiKey.(APIKey).Delete {
		c.JSON(401, gin.H{
			"error": "Missing delete permissions",
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

	os.Remove(config.DataDir + data.ID)

	if canUseS3() {
		client, err := getMinioClient()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
				"code":  "MinioClientError",
			})
		}

		err = client.RemoveObject(config.BucketName, id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
	}

	if isWebhookSet() {
		err = webhookDelete(id)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = deleteEntry(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error":    err.Error(),
			"code":     "DbError",
			"full_err": err,
		})
		return
	}

	c.JSON(200, data)
}
