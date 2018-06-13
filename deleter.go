package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func deleteFile(c *gin.Context) {
	apiKey, _ := c.Get("apikey")
	if !apiKey.(APIKey).Delete {
		c.JSON(401, gin.H{"error": "Missing delete permissions"})
		return
	}

	id := c.Param("id")

	os.Remove(config.DataDir + id)

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

	deleteEntry(id)
}
