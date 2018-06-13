package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

type APIKeyArgs struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
	Admin  bool `json:"admin"`
}

func createAPIKey(c *gin.Context) {
	args := APIKeyArgs{}

	err := c.BindJSON(&args)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	id := xid.New()

	if err = db.Create(APIKey{
		Key:    id.String(),
		Create: args.Create,
		Read:   args.Read,
		Update: args.Update,
		Delete: args.Delete,
		Admin:  args.Admin,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{
		"key": id.String(),
	})
}

func keyChecker() func(c *gin.Context) {
	return func(c *gin.Context) {
		apiKey := c.Query("key")

		var dbKey APIKey

		if err := db.First(&dbKey, APIKey{Key: apiKey}).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(401, gin.H{"error": "Invalid api key"})
				return
			}

			c.JSON(500, gin.H{"error": err})
			return
		}

		c.Set("apikey", dbKey)
		c.Next()
	}
}
