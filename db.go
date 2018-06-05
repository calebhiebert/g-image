package main

import (
	"encoding/json"
	"io/ioutil"
)

// Entry a file
type Entry struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Mime     string `json:"mime"`
	Size     int64  `json:"size"`
	Sha256   string `json:"hash"`
}

func writeEntry(entry Entry) error {
	json, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(DataDir+entry.ID+".json", json, 0666)
	if err != nil {
		return err
	}

	return nil
}

func readEntry(id string) (Entry, error) {
	entry := Entry{}

	jsonEntry, err := ioutil.ReadFile(DataDir + id + ".json")
	if err != nil {
		return entry, err
	}

	json.Unmarshal(jsonEntry, &entry)
	return entry, nil
}
