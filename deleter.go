package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func deleteFile(c *gin.Context) {
	id := c.Param("id")

	os.Remove(config.DataDir + id)
	os.Remove(config.DataDir + id + ".json")

	client, err := getMinioClient()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = client.RemoveObject(config.BucketName, id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
}
