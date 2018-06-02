package main

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func handlePOST(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	id := xid.New().String()

	saveFile(file, id)

	writeEntry(Entry{
		ID:       id,
		Filename: file.Filename,
		Mime:     file.Header.Get("Content-Type"),
		Size:     file.Size,
	})

	c.JSON(200, gin.H{
		"id":       id,
		"filename": file.Filename,
		"mime":     file.Header.Get("Content-Type"),
		"size":     file.Size,
	})
}

func saveFile(file *multipart.FileHeader, id string) {
	ensureDirectory(DataDir)

	f, err := os.OpenFile(DataDir+id, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()

	openedFile, err := file.Open()
	if err != nil {
		panic(err)
		return
	}
	defer openedFile.Close()

	io.Copy(f, openedFile)
}

func loadFile(id string) (*os.File, error) {
	data, err := os.Open(DataDir + id)
	if err != nil {
		panic(err)
		return nil, err
	}

	return data, nil
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}
}
