package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func putFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
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

	fileInfo := Entry{
		ID:       id,
		Filename: file.Filename,
		Mime:     file.Header.Get("Content-Type"),
		Size:     file.Size,
		Sha256:   hash,
	}

	writeEntry(fileInfo)

	c.JSON(200, fileInfo)
}

func saveFile(file *multipart.FileHeader, id string) (string, error) {
	ensureDirectory(DataDir)

	f, err := os.OpenFile(DataDir+id, os.O_WRONLY|os.O_CREATE, 0666)
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

	if _, err := io.Copy(h, openedFile); err != nil {
		return "", err
	}

	if _, err := io.Copy(f, openedFile); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func loadFile(id string) (*os.File, error) {
	data, err := os.Open(DataDir + id)
	if err != nil {
		panic(err)
	}

	return data, nil
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}
}

func uploadFile(details *Entry, file *multipart.FileHeader) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	openedFile, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer openedFile.Close()

	// TODO Upload
}
