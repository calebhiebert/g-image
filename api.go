package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jkomyno/nanoid"
)

// APIKeyArgs json arguments for creating an api key
type APIKeyArgs struct {
	Create    bool `json:"create"`
	Read      bool `json:"read"`
	Update    bool `json:"update"`
	Delete    bool `json:"delete"`
	Admin     bool `json:"admin"`
	SizeLimit int  `json:"sizeLimit"`
}

func createAPIKey(c *gin.Context) {
	key, _ := c.Get("apikey")
	if !key.(APIKey).Admin {
		c.JSON(403, gin.H{"error": "Missing admin permissions"})
		return
	}

	args := APIKeyArgs{}

	err := c.BindJSON(&args)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	id, _ := nanoid.Nanoid(48)

	sizeLimit := 52428800

	if args.SizeLimit > 0 {
		sizeLimit = args.SizeLimit
	}

	if err = db.Create(APIKey{
		Key:       id,
		Create:    args.Create,
		Read:      args.Read,
		Update:    args.Update,
		Delete:    args.Delete,
		Admin:     args.Admin,
		SizeLimit: sizeLimit,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{
		"key": id,
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

func getAndPrintAdminKey() error {
	key := APIKey{}

	if err := db.First(&key, APIKey{Create: true, Read: true, Update: true, Delete: true, Admin: true}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			id, _ := nanoid.Nanoid(48)

			key = APIKey{
				Key:       id,
				Create:    true,
				Read:      true,
				Update:    true,
				Delete:    true,
				Admin:     true,
				SizeLimit: 52428800,
			}

			if err = db.Create(&key).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	fmt.Printf("Admin API Key: %s\n", key.Key)
	return nil
}
