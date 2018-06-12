package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Entry a file
type Entry struct {
	ID       string `json:"id" gorm:"type:CHAR(20);UNIQUE_INDEX;PRIMARY_KEY"`
	Filename string `json:"filename"`
	Mime     string `json:"mime"`
	Size     int64  `json:"size"`
	Sha256   string `json:"hash"`
}

// APIKey a single api key
type APIKey struct {
	Key    string `json:"key" gorm:"type:CHAR(20);UNIQUE_INDEX;PRIMARY_KEY"`
	Create bool   `json:"create" gorm:"NOT_NULL;DEFAULT:true"`
	Read   bool   `json:"read"`
	Update bool   `json:"update"`
	Delete bool   `json:"delete"`
	Admin  bool   `json:"admin"`
}

var db *gorm.DB

func setupDB() error {
	var err error
	db, err = gorm.Open("sqlite3", "data.db")
	if err != nil {
		return err
	}

	db.AutoMigrate(&Entry{})
	db.AutoMigrate(&APIKey{})
	return nil
}

func writeEntry(entry Entry) error {
	if err := db.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}

func readEntry(id string) (Entry, error) {
	entry := Entry{}

	if err := db.Find(&entry, Entry{ID: id}).Error; err != nil {
		return entry, err
	}

	return entry, nil
}

func deleteEntry(id string) error {
	if err := db.Delete(Entry{ID: id}).Error; err != nil {
		return err
	}

	return nil
}
