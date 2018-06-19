package main

import (
	validation "github.com/go-ozzo/ozzo-validation"
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

// Validate validate an entry struct
func (e Entry) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.ID, validation.Required),
		validation.Field(&e.Filename, validation.Required),
		validation.Field(&e.Mime, validation.Required),
		validation.Field(&e.Size, validation.Required))
}

// APIKey a single api key
type APIKey struct {
	Key       string `json:"key" gorm:"type:CHAR(48);UNIQUE_INDEX;PRIMARY_KEY"`
	Create    bool   `json:"create" gorm:"NOT_NULL;DEFAULT:true"`
	Read      bool   `json:"read" gorm:"NOT_NULL;DEFAULT:true"`
	Update    bool   `json:"update" gorm:"NOT_NULL;DEFAULT:true"`
	Delete    bool   `json:"delete" gorm:"NOT_NULL;DEFAULT:true"`
	Admin     bool   `json:"admin" gorm:"NOT_NULL;DEFAULT:false"`
	SizeLimit int    `json:"sizeLimit" gorm:"NOT_NULL;DEFAULT:52428800"`
}

var db *gorm.DB

func setupDB() error {
	var err error
	db, err = gorm.Open("sqlite3", config.DataDir+"data.db")
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
