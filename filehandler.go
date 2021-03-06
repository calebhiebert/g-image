package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	minio "github.com/minio/minio-go"
	"github.com/rs/xid"
)

func putFile(c *gin.Context) {
	apiKey, _ := c.Get("apikey")
	if apiKey == nil || !apiKey.(APIKey).Create {
		c.JSON(401, gin.H{
			"error": "Missing create permissions",
			"code":  "PermError",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{
			"error":    err.Error(),
			"code":     "FormParseError",
			"full_err": err,
		})
		return
	}

	if file.Size > int64(apiKey.(APIKey).SizeLimit) {
		c.AbortWithStatusJSON(400, gin.H{
			"error": fmt.Sprintf("Size %d above limit %d", file.Size, apiKey.(APIKey).SizeLimit),
			"code":  "FileTooLarge",
		})
		return
	}

	id := xid.New().String()

	hash, err := saveFile(file, id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	defer func() {
		go cacheCheck()
	}()

	fileInfo := Entry{
		ID:       id,
		Filename: file.Filename,
		Mime:     file.Header.Get("Content-Type"),
		Size:     file.Size,
		Sha256:   hash,
	}

	if canUseS3() {
		go doUpload(&fileInfo, file)
	}

	err = writeEntry(fileInfo)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	if isWebhookSet() {
		err = webhookPutInfo(&fileInfo)
		if err != nil {
			fmt.Println(err)
		}
	}

	c.JSON(200, fileInfo)
}

func saveFile(file *multipart.FileHeader, id string) (string, error) {
	ensureDirectory(config.DataDir)

	f, err := os.OpenFile(config.DataDir+id, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	h := sha256.New()

	if _, err := io.Copy(f, io.TeeReader(openedFile, h)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func loadFile(id string) (*os.File, error) {
	data, err := os.Open(config.DataDir + id)
	if err != nil {
		panic(err)
	}

	return data, nil
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

func doUpload(details *Entry, file *multipart.FileHeader) {
	fmt.Printf("Starting file upload for %s\n", details.Filename)

	err := uploadFile(details, file)
	if err != nil {
		fmt.Println("Upload failed with error ", err)
	}

	fmt.Printf("Finished file upload for %s\n", details.Filename)
	go cacheCheck()
}

func uploadFile(details *Entry, file *multipart.FileHeader) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	openedFile, err := file.Open()
	if err != nil {
		return err
	}

	err = ensureBucket(config.BucketName)
	if err != nil {
		return err
	}

	_, err = client.PutObject(config.BucketName, details.ID, openedFile, details.Size, minio.PutObjectOptions{
		ContentType: details.Mime,
	})
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(id string) error {
	objectReader, err := getObjectReader(id)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(config.DataDir+id, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, objectReader)
	if err != nil {
		f.Close()
		os.Remove(config.DataDir + id)
		return err
	}

	return nil
}

func getObjectReader(id string) (*minio.Object, error) {
	client, err := getMinioClient()
	if err != nil {
		return nil, err
	}

	obj, err := client.GetObject(config.BucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func ensureBucket(bucketName string) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	exists, err := client.BucketExists(bucketName)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	return client.MakeBucket(bucketName, "nyc3")
}

func getMinioClient() (*minio.Client, error) {
	endpoint := config.S3Endpoint
	accessKeyID := config.S3AccessKey
	secretAccessKey := config.S3Secret
	useSSL := config.S3SSL

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return client, nil
}
