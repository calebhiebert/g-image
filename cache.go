package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func getCacheContents() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(config.DataDir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	return files, nil
}

func cacheCheck() {
	files, err := getCacheContents()
	if err != nil {
		fmt.Println(err)
	}

	var totalCacheSize int64
	var toDelete []string

	for _, file := range files {
		totalCacheSize += file.Size()

		if totalCacheSize/1000000 > config.CacheSize {
			toDelete = append(toDelete, file.Name())
		}
	}

	if len(toDelete) > 0 {
		for _, filePath := range toDelete {
			os.Remove(config.DataDir + filePath)
		}

		fmt.Printf("Removed %d entries from cache\n", len(toDelete))
	}
}
