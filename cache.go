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
	fmt.Println("Checking Cache")
	fmt.Printf("Allowed cache size is %dmb\n", config.CacheSize)

	files, err := getCacheContents()
	if err != nil {
		fmt.Println(err)
	}

	var totalCacheSize int64
	var toDelete []string

	fmt.Printf("Cache contains %d entries\n", len(files))

	for _, file := range files {
		totalCacheSize += file.Size()

		if totalCacheSize/1000000 > config.CacheSize {
			toDelete = append(toDelete, file.Name())
		}
	}

	fmt.Printf("Total Cache size: %dmb\n", totalCacheSize/1000000)

	if len(toDelete) > 0 {
		for _, filePath := range toDelete {
			os.Remove(config.DataDir + filePath)
		}

		fmt.Printf("Removed %d entries from cache\n", len(toDelete))
	}
}
