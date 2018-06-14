package main

import (
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	validator "gopkg.in/go-playground/validator.v8"
)

// Environment app config
type Environment struct {
	DataDir     string `env:"DATA_DIR" envDefault:"./data/"`
	BucketName  string `env:"BUCKET_NAME"`
	S3Endpoint  string `env:"S3_ENDPOINT"`
	S3AccessKey string `env:"S3_ACCESS_KEY"`
	S3Secret    string `env:"S3_SECRET"`
	S3SSL       bool   `env:"S3_SSL" envDefault:"false"`
	WebhookURL  string `env:"WEBHOOK_URL"`
	CacheSize   int64  `env:"CACHE_SIZE" envDefault:"50"` // Cache size in MB
}

var config Environment
var validate *validator.Validate

func main() {
	config = Environment{}
	validate = validator.New(&validator.Config{})

	err := godotenv.Load()
	if err != nil {
		println("Could not load .env file. Probably because there isn't one, and that is totally okay")
	}

	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	err = setupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	getAndPrintAdminKey()
	cacheCheck()

	r := gin.Default()
	r.Use(keyChecker())

	r.POST("/apikey", createAPIKey)

	r.GET("/:id", getFile)
	r.GET("/:id/info", getFileInfo)
	r.POST("/upload", putFile)
	r.DELETE("/:id", deleteFile)
	r.Run()
}
