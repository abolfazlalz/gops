package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Database() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	db, err := gorm.Open(sqlite.Open("db.sqlite"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
